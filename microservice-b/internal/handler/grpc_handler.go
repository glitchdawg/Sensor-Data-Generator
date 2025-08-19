package handler

import (
	"io"
	"log"
	"time"

	pb "github.com/glitchdawg/synthetic_sensors/proto/ingestpb"
	"github.com/glitchdawg/synthetic_sensors/microservice-b/internal/service"
	"github.com/glitchdawg/synthetic_sensors/shared/domain"
)

type GRPCHandler struct {
	pb.UnimplementedIngestServiceServer
	service *service.SensorService
}

func NewGRPCHandler(service *service.SensorService) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

func (h *GRPCHandler) Write(stream pb.IngestService_WriteServer) error {
	count := uint64(0)
	
	for {
		reading, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.WriteAck{Count: count})
		}
		if err != nil {
			log.Printf("stream receive error: %v", err)
			return err
		}

		// Parse timestamp
		ts, err := time.Parse(time.RFC3339Nano, reading.Timestamp)
		if err != nil {
			log.Printf("timestamp parse error: %v", err)
			ts = time.Now().UTC()
		}

		// Convert to domain model
		sensorReading := &domain.SensorReading{
			ID1:        reading.Id1,
			ID2:        int(reading.Id2),
			SensorType: reading.SensorType,
			Value:      reading.Value,
			Timestamp:  ts,
		}

		// Save to database
		if err := h.service.CreateReading(stream.Context(), sensorReading); err != nil {
			log.Printf("failed to save reading: %v", err)
			continue
		}

		count++
	}
}