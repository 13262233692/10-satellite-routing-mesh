package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

const (
	defaultDialTimeout      = 5 * time.Second
	defaultRequestTimeout   = 3 * time.Second
	topologyPrefix          = "/leo-routing/topology/"
	satellitePrefix         = "/leo-routing/satellites/"
	ephemerisPrefix         = "/leo-routing/ephemeris/"
	configKey               = "/leo-routing/config"
	leaderElectionKey       = "/leo-routing/leader"
)

type ConfigManager struct {
	client   *clientv3.Client
	topoMgr  TopologyUpdater
	config   Config
	watchers []clientv3.WatchChan
	mu       sync.RWMutex

	isLeader bool
	leaderMu sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
}

type Config struct {
	MaxLinkDistance float64 `json:"max_link_distance"`
	MaxLasersPerSat int     `json:"max_lasers_per_sat"`
	RebuildInterval int64   `json:"rebuild_interval_ms"`
}

type TopologyUpdater interface {
	UpdateEphemeris(ephem *model.Ephemeris)
	BatchUpdateEphemerides(ephems []model.Ephemeris)
	RebuildTopology() uint64
	AddSatellite(sat model.Satellite)
}

func NewConfigManager(endpoints []string, topoMgr TopologyUpdater) (*ConfigManager, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: defaultDialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cm := &ConfigManager{
		client:  cli,
		topoMgr: topoMgr,
		ctx:     ctx,
		cancel:  cancel,
		config: Config{
			MaxLinkDistance: model.MaxLinkDistanceMeters,
			MaxLasersPerSat: model.MaxLaserLinksPerSatellite,
			RebuildInterval: 1000,
		},
	}

	return cm, nil
}

func (cm *ConfigManager) Start() error {
	if err := cm.loadConfig(); err != nil {
		log.Printf("Warning: failed to load initial config: %v", err)
	}

	if err := cm.loadSatellites(); err != nil {
		log.Printf("Warning: failed to load satellites: %v", err)
	}

	if err := cm.loadEphemerides(); err != nil {
		log.Printf("Warning: failed to load ephemerides: %v", err)
	}

	go cm.watchConfig()
	go cm.watchSatellites()
	go cm.watchEphemerides()

	return nil
}

func (cm *ConfigManager) Stop() {
	cm.cancel()
	cm.client.Close()
}

func (cm *ConfigManager) loadConfig() error {
	ctx, cancel := context.WithTimeout(cm.ctx, defaultRequestTimeout)
	defer cancel()

	resp, err := cm.client.Get(ctx, configKey)
	if err != nil {
		return err
	}

	if len(resp.Kvs) == 0 {
		return nil
	}

	var cfg Config
	if err := json.Unmarshal(resp.Kvs[0].Value, &cfg); err != nil {
		return fmt.Errorf("invalid config JSON: %w", err)
	}

	cm.mu.Lock()
	cm.config = cfg
	cm.mu.Unlock()

	return nil
}

func (cm *ConfigManager) loadSatellites() error {
	ctx, cancel := context.WithTimeout(cm.ctx, defaultRequestTimeout)
	defer cancel()

	resp, err := cm.client.Get(ctx, satellitePrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		var sat model.Satellite
		if err := json.Unmarshal(kv.Value, &sat); err != nil {
			log.Printf("Warning: invalid satellite data for key %s: %v", kv.Key, err)
			continue
		}
		cm.topoMgr.AddSatellite(sat)
	}

	return nil
}

func (cm *ConfigManager) loadEphemerides() error {
	ctx, cancel := context.WithTimeout(cm.ctx, defaultRequestTimeout)
	defer cancel()

	resp, err := cm.client.Get(ctx, ephemerisPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	ephems := make([]model.Ephemeris, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		var ephem model.Ephemeris
		if err := json.Unmarshal(kv.Value, &ephem); err != nil {
			log.Printf("Warning: invalid ephemeris data for key %s: %v", kv.Key, err)
			continue
		}
		ephems = append(ephems, ephem)
	}

	if len(ephems) > 0 {
		cm.topoMgr.BatchUpdateEphemerides(ephems)
		cm.topoMgr.RebuildTopology()
	}

	return nil
}

func (cm *ConfigManager) watchConfig() {
	ch := cm.client.Watch(cm.ctx, configKey)
	for resp := range ch {
		for _, ev := range resp.Events {
			if ev.Type == clientv3.EventTypePut {
				var cfg Config
				if err := json.Unmarshal(ev.Kv.Value, &cfg); err != nil {
					log.Printf("Warning: invalid config update: %v", err)
					continue
				}
				cm.mu.Lock()
				cm.config = cfg
				cm.mu.Unlock()
				log.Println("Config updated from etcd")
			}
		}
	}
}

func (cm *ConfigManager) watchSatellites() {
	ch := cm.client.Watch(cm.ctx, satellitePrefix, clientv3.WithPrefix())
	for resp := range ch {
		for _, ev := range resp.Events {
			if ev.Type == clientv3.EventTypePut {
				var sat model.Satellite
				if err := json.Unmarshal(ev.Kv.Value, &sat); err != nil {
					log.Printf("Warning: invalid satellite update: %v", err)
					continue
				}
				cm.topoMgr.AddSatellite(sat)
			}
		}
	}
}

func (cm *ConfigManager) watchEphemerides() {
	ch := cm.client.Watch(cm.ctx, ephemerisPrefix, clientv3.WithPrefix())
	for resp := range ch {
		ephems := make([]model.Ephemeris, 0, len(resp.Events))
		for _, ev := range resp.Events {
			if ev.Type == clientv3.EventTypePut {
				var ephem model.Ephemeris
				if err := json.Unmarshal(ev.Kv.Value, &ephem); err != nil {
					log.Printf("Warning: invalid ephemeris update: %v", err)
					continue
				}
				ephems = append(ephems, ephem)
			}
		}
		if len(ephems) > 0 {
			cm.topoMgr.BatchUpdateEphemerides(ephems)
		}
	}
}

func (cm *ConfigManager) GetConfig() Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

func (cm *ConfigManager) CampaignLeader(ttl int) error {
	sess, err := concurrency.NewSession(cm.client, concurrency.WithTTL(ttl))
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	election := concurrency.NewElection(sess, leaderElectionKey)

	go func() {
		for {
			if err := election.Campaign(cm.ctx, "routing-node"); err != nil {
				if cm.ctx.Err() != nil {
					return
				}
				log.Printf("Leader campaign error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			cm.leaderMu.Lock()
			cm.isLeader = true
			cm.leaderMu.Unlock()
			log.Println("Became leader")

			<-sess.Done()

			cm.leaderMu.Lock()
			cm.isLeader = false
			cm.leaderMu.Unlock()
			log.Println("Lost leadership")

			if cm.ctx.Err() != nil {
				return
			}

			sess, _ = concurrency.NewSession(cm.client, concurrency.WithTTL(ttl))
			election = concurrency.NewElection(sess, leaderElectionKey)
		}
	}()

	return nil
}

func (cm *ConfigManager) IsLeader() bool {
	cm.leaderMu.RLock()
	defer cm.leaderMu.RUnlock()
	return cm.isLeader
}

func (cm *ConfigManager) PublishEphemeris(ctx context.Context, ephem *model.Ephemeris) error {
	data, err := json.Marshal(ephem)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", ephemerisPrefix, ephem.ID)
	_, err = cm.client.Put(ctx, key, string(data))
	return err
}

func (cm *ConfigManager) PublishSatellite(ctx context.Context, sat *model.Satellite) error {
	data, err := json.Marshal(sat)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", satellitePrefix, sat.ID)
	_, err = cm.client.Put(ctx, key, string(data))
	return err
}
