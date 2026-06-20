package model

type RouteResult struct {
	Path      []SatelliteID
	TotalDist float64
	Latency   float64
	Hops      int
	Found     bool
}

type RouteRequest struct {
	From      SatelliteID
	To        SatelliteID
	Timestamp int64
	Algorithm string
}
