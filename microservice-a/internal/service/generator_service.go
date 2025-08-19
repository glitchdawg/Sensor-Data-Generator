package service

import (
	"context"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	pb "github.com/glitchdawg/synthetic_sensors/proto/ingestpb"
	"github.com/glitchdawg/synthetic_sensors/microservice-a/internal/domain"
	"google.golang.org/grpc"
)

type GeneratorService struct {
	client     pb.IngestServiceClient
	config     *domain.GeneratorConfig
	frequency  *int64
}

func NewGeneratorService(conn *grpc.ClientConn, config *domain.GeneratorConfig) *GeneratorService {
	freq := config.FrequencyMs
	return &GeneratorService{
		client:    pb.NewIngestServiceClient(conn),
		config:    config,
		frequency: &freq,
	}
}

func (s *GeneratorService) UpdateFrequency(freq int64) {
	atomic.StoreInt64(s.frequency, freq)
}

func (s *GeneratorService) GetFrequency() int64 {
	return atomic.LoadInt64(s.frequency)
}

func (s *GeneratorService) StartGenerator(ctx context.Context) error {
	stream, err := s.client.Write(ctx)
	if err != nil {
		return err
	}

	idPool := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	for {
		select {
		case <-ctx.Done():
			return stream.CloseSend()
		default:
			now := time.Now().UTC().Format(time.RFC3339Nano)
			msg := &pb.Reading{
				Value:      rand.Float64() * 100,
				SensorType: s.config.SensorType,
				Id1:        idPool[rand.Intn(len(idPool))],
				Id2:        int32(rand.Intn(100)),
				Timestamp:  now,
			}
			
			if err := stream.Send(msg); err != nil {
				log.Printf("send error: %v", err)
			}
			
			time.Sleep(time.Duration(s.GetFrequency()) * time.Millisecond)
		}
	}
}