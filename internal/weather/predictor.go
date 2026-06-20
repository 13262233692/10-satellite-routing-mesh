package weather

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

const (
	DefaultHistoryWindowSize   = 60
	DefaultPredictionHorizonSec = 300
	DefaultBaseSNR             = 30.0
	DefaultSNRThreshold        = 10.0
	DefaultMaxWeightMult       = 10.0
	DefaultMinWeightMult       = 1.0
)

type LinkQualityPredictor struct {
	mu       sync.RWMutex
	windows  map[uint64]*model.MovingWindowStats
	attenMap map[uint64]float64

	historySize    int
	predictionSec  int
	baseSNR        float64
	snrThreshold   float64
	maxWeightMult  float64
	minWeightMult  float64

	lastUpdate atomic.Int64
	predictionCount atomic.Int64
}

type LinkKey struct {
	A model.SatelliteID
	B model.SatelliteID
}

func NewLinkQualityPredictor() *LinkQualityPredictor {
	return &LinkQualityPredictor{
		windows:       make(map[uint64]*model.MovingWindowStats),
		attenMap:      make(map[uint64]float64),
		historySize:   DefaultHistoryWindowSize,
		predictionSec: DefaultPredictionHorizonSec,
		baseSNR:       DefaultBaseSNR,
		snrThreshold:  DefaultSNRThreshold,
		maxWeightMult: DefaultMaxWeightMult,
		minWeightMult: DefaultMinWeightMult,
	}
}

func EncodeLinkID(a, b model.SatelliteID) uint64 {
	if a > b {
		a, b = b, a
	}
	return uint64(a)<<32 | uint64(b)
}

func (p *LinkQualityPredictor) IngestData(data *model.SpaceWeatherData) {
	linkID := EncodeLinkID(data.SatA, data.SatB)

	p.mu.Lock()
	window, ok := p.windows[linkID]
	if !ok {
		window = model.NewMovingWindowStats(p.historySize)
		p.windows[linkID] = window
	}

	snr := data.SNRdB
	if snr == 0 && data.Attenuation > 0 {
		snr = p.baseSNR - data.Attenuation
	}

	window.Add(data.Timestamp, snr)

	if data.Attenuation > 0 {
		p.attenMap[linkID] = data.Attenuation
	}

	p.mu.Unlock()
	p.lastUpdate.Store(time.Now().UnixNano())
}

func (p *LinkQualityPredictor) BatchIngest(dataList []model.SpaceWeatherData) {
	now := time.Now().UnixNano()

	p.mu.Lock()
	for _, data := range dataList {
		linkID := EncodeLinkID(data.SatA, data.SatB)

		window, ok := p.windows[linkID]
		if !ok {
			window = model.NewMovingWindowStats(p.historySize)
			p.windows[linkID] = window
		}

		snr := data.SNRdB
		if snr == 0 && data.Attenuation > 0 {
			snr = p.baseSNR - data.Attenuation
		}

		window.Add(data.Timestamp, snr)

		if data.Attenuation > 0 {
			p.attenMap[linkID] = data.Attenuation
		}
	}
	p.mu.Unlock()
	p.lastUpdate.Store(now)
	p.predictionCount.Add(int64(len(dataList)))
}

func (p *LinkQualityPredictor) PredictLinkQuality(satA, satB model.SatelliteID, nowUnixNano int64) model.LinkQualityPrediction {
	linkID := EncodeLinkID(satA, satB)

	p.mu.RLock()
	window, ok := p.windows[linkID]
	curAtten := p.attenMap[linkID]
	p.mu.RUnlock()

	pred := model.LinkQualityPrediction{
		LinkID:           linkID,
		SatA:             satA,
		SatB:             satB,
		WeightMultiplier: p.minWeightMult,
		WarningLevel:     0,
		Confidence:       0,
		PredictionTime:   nowUnixNano,
	}

	if !ok || window.Count() < 3 {
		pred.CurrentSNR = p.baseSNR - curAtten
		pred.PredictedSNR = pred.CurrentSNR
		return pred
	}

	curSNR := window.Mean()
	pred.CurrentSNR = curSNR
	pred.Confidence = float64(minInt(window.Count(), p.historySize)) / float64(p.historySize)

	futureTime := nowUnixNano + int64(p.predictionSec)*int64(time.Second)
	predSNR, ok := window.PredictAt(futureTime)
	if !ok {
		pred.PredictedSNR = curSNR
		return pred
	}

	pred.PredictedSNR = predSNR

	slope, _, _ := window.LinearTrend()
	pred.Trend = slope

	weightMult := p.computeWeightMultiplier(curSNR, predSNR, pred.Confidence)
	pred.WeightMultiplier = weightMult
	pred.WarningLevel = p.computeWarningLevel(curSNR, predSNR)

	return pred
}

func (p *LinkQualityPredictor) computeWeightMultiplier(currentSNR, predictedSNR, confidence float64) float64 {
	degradation := currentSNR - predictedSNR
	if degradation < 0 {
		degradation = 0
	}

	ratio := degradation / (p.baseSNR - p.snrThreshold)
	if ratio < 0 {
		ratio = 0
	}

	if predictedSNR <= p.snrThreshold {
		return p.maxWeightMult
	}

	adjRatio := ratio * confidence
	mult := p.minWeightMult + adjRatio*(p.maxWeightMult-p.minWeightMult)

	if mult > p.maxWeightMult {
		mult = p.maxWeightMult
	}
	if mult < p.minWeightMult {
		mult = p.minWeightMult
	}

	return mult
}

func (p *LinkQualityPredictor) computeWarningLevel(currentSNR, predictedSNR float64) int32 {
	degradation := currentSNR - predictedSNR
	worst := math.Min(currentSNR, predictedSNR)

	if worst <= p.snrThreshold {
		return 3
	}
	if degradation > 10 {
		return 2
	}
	if degradation > 5 {
		return 1
	}
	return 0
}

func (p *LinkQualityPredictor) GetWeightMultiplier(satA, satB model.SatelliteID) float64 {
	linkID := EncodeLinkID(satA, satB)

	p.mu.RLock()
	window, ok := p.windows[linkID]
	curAtten := p.attenMap[linkID]
	p.mu.RUnlock()

	if !ok || window.Count() < 3 {
		if curAtten > 0 {
			snr := p.baseSNR - curAtten
			if snr <= p.snrThreshold {
				return p.maxWeightMult
			}
			ratio := curAtten / (p.baseSNR - p.snrThreshold)
			mult := p.minWeightMult + ratio*(p.maxWeightMult-p.minWeightMult)
			if mult > p.maxWeightMult {
				mult = p.maxWeightMult
			}
			return mult
		}
		return p.minWeightMult
	}

	curSNR := window.Mean()
	futureTime := time.Now().UnixNano() + int64(p.predictionSec)*int64(time.Second)
	predSNR, ok := window.PredictAt(futureTime)
	if !ok {
		predSNR = curSNR
	}

	return p.computeWeightMultiplier(curSNR, predSNR, 0.8)
}

func (p *LinkQualityPredictor) GetAllAffectedLinks() []LinkKey {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]LinkKey, 0, len(p.windows))
	for linkID, window := range p.windows {
		if window.Count() < 3 {
			continue
		}
		curSNR := window.Mean()
		futureTime := time.Now().UnixNano() + int64(p.predictionSec)*int64(time.Second)
		predSNR, _ := window.PredictAt(futureTime)

		if curSNR-predSNR > 3 || curSNR < p.snrThreshold+5 {
			a := model.SatelliteID(linkID >> 32)
			b := model.SatelliteID(linkID & 0xFFFFFFFF)
			result = append(result, LinkKey{A: a, B: b})
		}
	}

	return result
}

func (p *LinkQualityPredictor) ActiveLinkCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.windows)
}

func (p *LinkQualityPredictor) LastUpdateTime() int64 {
	return p.lastUpdate.Load()
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (p *LinkQualityPredictor) GetStats() (activeLinks int, totalSamples int64) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	total := int64(0)
	for _, w := range p.windows {
		total += int64(w.Count())
	}
	return len(p.windows), total
}
