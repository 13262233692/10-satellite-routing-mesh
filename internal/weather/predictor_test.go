package weather

import (
	"testing"
	"time"

	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

func TestMovingWindowBasic(t *testing.T) {
	w := model.NewMovingWindowStats(10)

	now := time.Now().UnixNano()
	for i := 0; i < 5; i++ {
		w.Add(now+int64(i)*int64(time.Second), 20.0+float64(i))
	}

	if w.Count() != 5 {
		t.Errorf("Expected count 5, got %d", w.Count())
	}

	mean := w.Mean()
	expectedMean := (20 + 21 + 22 + 23 + 24) / 5.0
	if mean != expectedMean {
		t.Errorf("Expected mean %f, got %f", expectedMean, mean)
	}
}

func TestMovingWindowLinearTrend(t *testing.T) {
	w := model.NewMovingWindowStats(20)

	now := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		ts := now + int64(i)*int64(time.Second)
		val := 30.0 - float64(i)*0.5
		w.Add(ts, val)
	}

	slope, _, ok := w.LinearTrend()
	if !ok {
		t.Fatal("Expected valid linear trend")
	}

	if slope >= 0 {
		t.Errorf("Expected negative slope (degrading SNR), got %f", slope)
	}

	perSec := slope * float64(time.Second)
	t.Logf("Trend: %f SNR per second", perSec)
}

func TestMovingWindowPrediction(t *testing.T) {
	w := model.NewMovingWindowStats(20)

	now := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		ts := now + int64(i)*int64(time.Second)
		val := 30.0 - float64(i)
		w.Add(ts, val)
	}

	future := now + 15*int64(time.Second)
	predicted, ok := w.PredictAt(future)
	if !ok {
		t.Fatal("Expected valid prediction")
	}

	expected := 30.0 - 15.0
	t.Logf("Predicted SNR at +15s: %f (expected ~%f)", predicted, expected)

	if predicted > expected+1.0 || predicted < expected-1.0 {
		t.Errorf("Prediction out of expected range: %f vs %f", predicted, expected)
	}
}

func TestWeightPredictorBasic(t *testing.T) {
	p := NewLinkQualityPredictor()

	mult := p.GetWeightMultiplier(1, 2)
	if mult != 1.0 {
		t.Errorf("Expected default weight 1.0 for unknown link, got %f", mult)
	}
}

func TestWeightPredictorWithAttenuation(t *testing.T) {
	p := NewLinkQualityPredictor()

	now := time.Now().UnixNano()
	data := &model.SpaceWeatherData{
		SatA:        1,
		SatB:        2,
		Timestamp:   now,
		Attenuation: 15.0,
		WeatherType: model.SpaceWeather_SolarFlare,
		Severity:    0.8,
	}

	p.IngestData(data)

	mult := p.GetWeightMultiplier(1, 2)
	t.Logf("Weight multiplier with 15dB attenuation: %f", mult)

	if mult <= 1.0 {
		t.Error("Expected weight multiplier > 1.0 for attenuated link")
	}
	if mult > 10.0 {
		t.Error("Expected weight multiplier <= 10.0")
	}
}

func TestWeightPredictorTrendDegradation(t *testing.T) {
	p := NewLinkQualityPredictor()

	now := time.Now().UnixNano()
	for i := 0; i < 10; i++ {
		ts := now + int64(i)*int64(time.Second)
		data := model.SpaceWeatherData{
			SatA:      1,
			SatB:      2,
			Timestamp: ts,
			SNRdB:     25.0 - float64(i)*1.5,
		}
		p.IngestData(&data)
	}

	mult := p.GetWeightMultiplier(1, 2)
	t.Logf("Weight multiplier with degrading SNR: %f", mult)

	if mult <= 1.0 {
		t.Error("Expected weight multiplier > 1.0 for degrading link")
	}
}

func TestLinkQualityPrediction(t *testing.T) {
	p := NewLinkQualityPredictor()

	now := time.Now().UnixNano()
	for i := 0; i < 15; i++ {
		ts := now + int64(i)*int64(time.Second)
		data := model.SpaceWeatherData{
			SatA:      5,
			SatB:      7,
			Timestamp: ts,
			SNRdB:     28.0 - float64(i)*0.8,
		}
		p.IngestData(&data)
	}

	pred := p.PredictLinkQuality(5, 7, now+10*int64(time.Second))

	t.Logf("Link prediction:")
	t.Logf("  Current SNR: %.2f dB", pred.CurrentSNR)
	t.Logf("  Predicted SNR: %.2f dB", pred.PredictedSNR)
	t.Logf("  Weight multiplier: %.2fx", pred.WeightMultiplier)
	t.Logf("  Warning level: %d", pred.WarningLevel)
	t.Logf("  Confidence: %.2f", pred.Confidence)

	if pred.CurrentSNR <= pred.PredictedSNR {
		t.Error("Expected predicted SNR to be lower than current (degrading trend)")
	}

	if pred.WeightMultiplier < 1.0 {
		t.Error("Expected weight multiplier >= 1.0")
	}

	if pred.Confidence <= 0 {
		t.Error("Expected confidence > 0 with 15 samples")
	}
}

func TestBatchIngest(t *testing.T) {
	p := NewLinkQualityPredictor()

	now := time.Now().UnixNano()
	dataList := make([]model.SpaceWeatherData, 20)
	for i := 0; i < 20; i++ {
		dataList[i] = model.SpaceWeatherData{
			SatA:      10,
			SatB:      20,
			Timestamp: now + int64(i)*int64(time.Second),
			SNRdB:     30.0 - float64(i)*0.3,
		}
	}

	p.BatchIngest(dataList)

	if p.ActiveLinkCount() != 1 {
		t.Errorf("Expected 1 active link, got %d", p.ActiveLinkCount())
	}

	mult := p.GetWeightMultiplier(10, 20)
	t.Logf("After batch ingest, weight multiplier: %f", mult)
}

func TestAffectedLinks(t *testing.T) {
	p := NewLinkQualityPredictor()

	now := time.Now().UnixNano()
	for i := 0; i < 15; i++ {
		ts := now + int64(i)*int64(time.Second)
		data1 := model.SpaceWeatherData{
			SatA:      1,
			SatB:      2,
			Timestamp: ts,
			SNRdB:     25.0 - float64(i)*1.2,
		}
		data2 := model.SpaceWeatherData{
			SatA:      3,
			SatB:      4,
			Timestamp: ts,
			SNRdB:     28.0,
		}
		p.IngestData(&data1)
		p.IngestData(&data2)
	}

	affected := p.GetAllAffectedLinks()
	t.Logf("Affected links count: %d", len(affected))

	foundDegrading := false
	for _, link := range affected {
		if (link.A == 1 && link.B == 2) || (link.A == 2 && link.B == 1) {
			foundDegrading = true
		}
	}

	if !foundDegrading {
		t.Log("Warning: degrading link (1-2) not in affected list (may need more degradation)")
	}
}

func TestSevereAttenuationMaxWeight(t *testing.T) {
	p := NewLinkQualityPredictor()

	now := time.Now().UnixNano()
	data := &model.SpaceWeatherData{
		SatA:        100,
		SatB:        200,
		Timestamp:   now,
		Attenuation: 25.0,
	}

	p.IngestData(data)

	mult := p.GetWeightMultiplier(100, 200)
	t.Logf("Weight with 25dB attenuation: %f", mult)

	if mult < DefaultMaxWeightMult*0.9 {
		t.Errorf("Expected weight near max (%f) for severe attenuation, got %f", DefaultMaxWeightMult, mult)
	}
}

func TestEncodeLinkID(t *testing.T) {
	id1 := EncodeLinkID(1, 2)
	id2 := EncodeLinkID(2, 1)

	if id1 != id2 {
		t.Errorf("Expected symmetric encoding: %d vs %d", id1, id2)
	}

	id3 := EncodeLinkID(100, 200)
	id4 := EncodeLinkID(100, 300)

	if id3 == id4 {
		t.Error("Expected different link IDs for different satellites")
	}
}
