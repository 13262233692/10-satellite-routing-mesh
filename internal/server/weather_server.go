package server

import (
	"context"
	"time"

	routingpb "github.com/aerospace/leo-routing-mesh/api/proto"
	"github.com/aerospace/leo-routing-mesh/internal/weather"
	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

type SpaceWeatherServer struct {
	routingpb.UnimplementedSpaceWeatherServiceServer
	predictor *weather.LinkQualityPredictor
}

func NewSpaceWeatherServer(predictor *weather.LinkQualityPredictor) *SpaceWeatherServer {
	return &SpaceWeatherServer{
		predictor: predictor,
	}
}

func (s *SpaceWeatherServer) IngestSpaceWeather(ctx context.Context, req *routingpb.IngestSpaceWeatherRequest) (*routingpb.IngestSpaceWeatherResponse, error) {
	data := req.GetData()
	if data == nil {
		return &routingpb.IngestSpaceWeatherResponse{Success: false}, nil
	}

	swData := &model.SpaceWeatherData{
		LinkID:         data.GetLinkId(),
		SatA:           model.SatelliteID(data.GetSatA()),
		SatB:           model.SatelliteID(data.GetSatB()),
		Timestamp:      data.GetTimestamp(),
		SNRdB:          data.GetSnrDb(),
		Attenuation:    data.GetAttenuation(),
		WeatherType:    model.SpaceWeatherType(data.GetWeatherType()),
		Severity:       data.GetSeverity(),
		CloudThickness: data.GetCloudThickness(),
		ParticleFlux:   data.GetParticleFlux(),
	}

	if swData.Timestamp == 0 {
		swData.Timestamp = time.Now().UnixNano()
	}

	s.predictor.IngestData(swData)

	return &routingpb.IngestSpaceWeatherResponse{Success: true}, nil
}

func (s *SpaceWeatherServer) BatchIngestSpaceWeather(ctx context.Context, req *routingpb.BatchIngestSpaceWeatherRequest) (*routingpb.BatchIngestSpaceWeatherResponse, error) {
	dataList := req.GetData()
	if len(dataList) == 0 {
		return &routingpb.BatchIngestSpaceWeatherResponse{Success: true, IngestedCount: 0}, nil
	}

	modelData := make([]model.SpaceWeatherData, len(dataList))
	now := time.Now().UnixNano()

	for i, data := range dataList {
		modelData[i] = model.SpaceWeatherData{
			LinkID:         data.GetLinkId(),
			SatA:           model.SatelliteID(data.GetSatA()),
			SatB:           model.SatelliteID(data.GetSatB()),
			Timestamp:      data.GetTimestamp(),
			SNRdB:          data.GetSnrDb(),
			Attenuation:    data.GetAttenuation(),
			WeatherType:    model.SpaceWeatherType(data.GetWeatherType()),
			Severity:       data.GetSeverity(),
			CloudThickness: data.GetCloudThickness(),
			ParticleFlux:   data.GetParticleFlux(),
		}
		if modelData[i].Timestamp == 0 {
			modelData[i].Timestamp = now
		}
	}

	s.predictor.BatchIngest(modelData)

	return &routingpb.BatchIngestSpaceWeatherResponse{
		Success:       true,
		IngestedCount: int32(len(dataList)),
	}, nil
}

func (s *SpaceWeatherServer) GetLinkQualityPrediction(ctx context.Context, req *routingpb.GetLinkQualityPredictionRequest) (*routingpb.GetLinkQualityPredictionResponse, error) {
	satA := model.SatelliteID(req.GetSatA())
	satB := model.SatelliteID(req.GetSatB())

	now := time.Now().UnixNano()
	pred := s.predictor.PredictLinkQuality(satA, satB, now)

	pbPred := &routingpb.LinkQualityPrediction{
		LinkId:           pred.LinkID,
		SatA:             uint32(pred.SatA),
		SatB:             uint32(pred.SatB),
		CurrentSnr:       pred.CurrentSNR,
		PredictedSnr:     pred.PredictedSNR,
		Trend:            pred.Trend,
		Confidence:       pred.Confidence,
		PredictionTime:   pred.PredictionTime,
		WeightMultiplier: pred.WeightMultiplier,
		WarningLevel:     pred.WarningLevel,
	}

	return &routingpb.GetLinkQualityPredictionResponse{
		Found:      pred.Confidence > 0,
		Prediction: pbPred,
	}, nil
}

func (s *SpaceWeatherServer) GetAffectedLinks(ctx context.Context, req *routingpb.GetAffectedLinksRequest) (*routingpb.GetAffectedLinksResponse, error) {
	affected := s.predictor.GetAllAffectedLinks()
	now := time.Now().UnixNano()

	result := make([]*routingpb.LinkQualityPrediction, 0, len(affected))
	for _, key := range affected {
		pred := s.predictor.PredictLinkQuality(key.A, key.B, now)
		result = append(result, &routingpb.LinkQualityPrediction{
			LinkId:           pred.LinkID,
			SatA:             uint32(pred.SatA),
			SatB:             uint32(pred.SatB),
			CurrentSnr:       pred.CurrentSNR,
			PredictedSnr:     pred.PredictedSNR,
			Trend:            pred.Trend,
			Confidence:       pred.Confidence,
			PredictionTime:   pred.PredictionTime,
			WeightMultiplier: pred.WeightMultiplier,
			WarningLevel:     pred.WarningLevel,
		})
	}

	return &routingpb.GetAffectedLinksResponse{
		AffectedLinks: result,
	}, nil
}
