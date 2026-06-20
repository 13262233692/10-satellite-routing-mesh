package model

const (
	DefaultPredictionWindowSec = 300
	DefaultHistoryWindowSize   = 60
)

type SpaceWeatherType int32

const (
	SpaceWeather_None           SpaceWeatherType = 0
	SpaceWeather_SolarFlare     SpaceWeatherType = 1
	SpaceWeather_CosmicRay      SpaceWeatherType = 2
	SpaceWeather_Atmospheric    SpaceWeatherType = 3
	SpaceWeather_MagneticStorm  SpaceWeatherType = 4
)

type SpaceWeatherData struct {
	LinkID      uint64
	SatA        SatelliteID
	SatB        SatelliteID
	Timestamp   int64
	SNRdB       float64
	Attenuation float64
	WeatherType SpaceWeatherType
	Severity    float64
	CloudThickness float64
	ParticleFlux float64
}

type LinkQualityPrediction struct {
	LinkID         uint64
	SatA           SatelliteID
	SatB           SatelliteID
	CurrentSNR     float64
	PredictedSNR   float64
	Trend          float64
	Confidence     float64
	PredictionTime int64
	WeightMultiplier float64
	WarningLevel   int32
}

type SlidingWindowSample struct {
	Timestamp int64
	Value     float64
}

type MovingWindowStats struct {
	samples []SlidingWindowSample
	head    int
	tail    int
	count   int
	size    int

	sum    float64
	sumSq  float64
	sumXY  float64
	sumX   float64
	sumXSq float64

	windowStart int64
	windowEnd   int64
}

func NewMovingWindowStats(windowSize int) *MovingWindowStats {
	if windowSize < 2 {
		windowSize = 2
	}
	return &MovingWindowStats{
		samples: make([]SlidingWindowSample, windowSize),
		size:    windowSize,
	}
}

func (m *MovingWindowStats) Add(timestamp int64, value float64) {
	if m.count >= m.size {
		old := m.samples[m.tail]
		m.sum -= old.Value
		m.sumSq -= old.Value * old.Value
		m.sumXY -= float64(old.Timestamp) * old.Value
		m.sumX -= float64(old.Timestamp)
		m.sumXSq -= float64(old.Timestamp) * float64(old.Timestamp)

		m.tail = (m.tail + 1) % m.size
		m.count--
	}

	idx := (m.tail + m.count) % m.size
	m.samples[idx] = SlidingWindowSample{
		Timestamp: timestamp,
		Value:     value,
	}

	m.sum += value
	m.sumSq += value * value
	m.sumXY += float64(timestamp) * value
	m.sumX += float64(timestamp)
	m.sumXSq += float64(timestamp) * float64(timestamp)
	m.count++

	if m.windowStart == 0 || timestamp < m.windowStart {
		m.windowStart = timestamp
	}
	if timestamp > m.windowEnd {
		m.windowEnd = timestamp
	}
}

func (m *MovingWindowStats) Mean() float64 {
	if m.count == 0 {
		return 0
	}
	return m.sum / float64(m.count)
}

func (m *MovingWindowStats) StdDev() float64 {
	if m.count < 2 {
		return 0
	}
	mean := m.Mean()
	variance := (m.sumSq/float64(m.count) - mean*mean)
	if variance < 0 {
		return 0
	}
	return mathSqrt(variance)
}

func (m *MovingWindowStats) LinearTrend() (slope float64, intercept float64, ok bool) {
	if m.count < 2 {
		return 0, 0, false
	}

	baseT := float64(m.windowStart)

	var sumX, sumY, sumXY, sumXSq float64
	n := float64(m.count)

	for i := 0; i < m.count; i++ {
		idx := (m.tail + i) % m.size
		s := m.samples[idx]
		x := float64(s.Timestamp) - baseT
		y := s.Value
		sumX += x
		sumY += y
		sumXY += x * y
		sumXSq += x * x
	}

	denom := n*sumXSq - sumX*sumX
	if denom == 0 {
		return 0, 0, false
	}

	slope = (n*sumXY - sumX*sumY) / denom
	interceptRel := (sumY - slope*sumX) / n

	intercept = interceptRel - slope*baseT

	return slope, intercept, true
}

func (m *MovingWindowStats) PredictAt(timestamp int64) (float64, bool) {
	slope, intercept, ok := m.LinearTrend()
	if !ok {
		return 0, false
	}
	return slope*float64(timestamp) + intercept, true
}

func (m *MovingWindowStats) Count() int {
	return m.count
}

func (m *MovingWindowStats) Reset() {
	m.head = 0
	m.tail = 0
	m.count = 0
	m.sum = 0
	m.sumSq = 0
	m.sumXY = 0
	m.sumX = 0
	m.sumXSq = 0
	m.windowStart = 0
	m.windowEnd = 0
}

func mathSqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}
