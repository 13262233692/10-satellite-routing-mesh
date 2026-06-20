package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	routingpb "github.com/aerospace/leo-routing-mesh/api/proto"
	etcdmgr "github.com/aerospace/leo-routing-mesh/internal/etcd"
	"github.com/aerospace/leo-routing-mesh/internal/routing"
	"github.com/aerospace/leo-routing-mesh/internal/server"
	"github.com/aerospace/leo-routing-mesh/internal/topology"
	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

var (
	grpcPort        = flag.Int("port", 50051, "gRPC server port")
	maxSatellites   = flag.Int("max-satellites", model.MaxSatellites, "Maximum number of satellites")
	etcdEndpoints   = flag.String("etcd-endpoints", "localhost:2379", "Comma-separated etcd endpoints")
	useEtcd         = flag.Bool("use-etcd", false, "Enable etcd integration")
	rebuildInterval = flag.Int("rebuild-interval-ms", 1000, "Topology rebuild interval in milliseconds")
	cacheCapacity   = flag.Int("cache-capacity", 10000, "Route cache capacity")
	enableCache     = flag.Bool("enable-cache", true, "Enable route result caching")
)

func main() {
	flag.Parse()

	log.Printf("Starting LEO Routing Mesh Server...")
	log.Printf("Max satellites: %d", *maxSatellites)
	log.Printf("gRPC port: %d", *grpcPort)
	log.Printf("Rebuild interval: %dms", *rebuildInterval)
	log.Printf("Route cache enabled: %v, capacity: %d", *enableCache, *cacheCapacity)

	topoMgr := topology.NewTopologyManager(*maxSatellites)
	router := routing.NewRouter()

	var routeCache *routing.RouteCache
	if *enableCache {
		routeCache = routing.NewRouteCache(*cacheCapacity)
	}

	routingSrv := server.NewRoutingServer(topoMgr, router, routeCache)
	ephemerisSrv := server.NewEphemerisServer(topoMgr)
	topologySrv := server.NewTopologyServer(topoMgr)

	var etcdMgr *etcdmgr.ConfigManager
	if *useEtcd {
		endpoints := strings.Split(*etcdEndpoints, ",")
		log.Printf("Connecting to etcd endpoints: %v", endpoints)

		var err error
		etcdMgr, err = etcdmgr.NewConfigManager(endpoints, topoMgr)
		if err != nil {
			log.Printf("Warning: failed to create etcd manager: %v", err)
		} else {
			if err := etcdMgr.Start(); err != nil {
				log.Printf("Warning: failed to start etcd manager: %v", err)
			} else {
				log.Println("Etcd manager started successfully")
			}
		}
	}

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(10000),
	)

	routingpb.RegisterRoutingServiceServer(grpcServer, routingSrv)
	routingpb.RegisterEphemerisServiceServer(grpcServer, ephemerisSrv)
	routingpb.RegisterTopologyServiceServer(grpcServer, topologySrv)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go startRebuildTicker(topoMgr, time.Duration(*rebuildInterval)*time.Millisecond)

	go func() {
		log.Printf("gRPC server listening on :%d", *grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %v, shutting down...", sig)

	grpcServer.GracefulStop()

	if etcdMgr != nil {
		etcdMgr.Stop()
	}

	log.Println("Server stopped gracefully")
}

func startRebuildTicker(topoMgr *topology.TopologyManager, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	count := 0
	for range ticker.C {
		start := time.Now()
		topoMgr.RebuildTopology()
		elapsed := time.Since(start)
		count++

		if count%10 == 0 {
			snapshot := topoMgr.GetSnapshot()
			log.Printf("Topology rebuilt #%d: epoch=%d, nodes=%d, time=%v", count, snapshot.Epoch, snapshot.Matrix.Nodes(), elapsed)
		}
	}
}
