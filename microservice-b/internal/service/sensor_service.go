package service

import (
	"context"
	"fmt"
	"time"

	"github.com/glitchdawg/synthetic_sensors/microservice-b/internal/domain"
	sharedDomain "github.com/glitchdawg/synthetic_sensors/shared/domain"
)

type SensorService struct {
	repo domain.SensorReadingRepository
}

func NewSensorService(repo domain.SensorReadingRepository) *SensorService {
	return &SensorService{repo: repo}
}

func (s *SensorService) CreateReading(ctx context.Context, reading *sharedDomain.SensorReading) error {
	if reading.Timestamp.IsZero() {
		reading.Timestamp = time.Now().UTC()
	}
	return s.repo.Create(ctx, reading)
}

func (s *SensorService) GetReadings(ctx context.Context, filter *sharedDomain.SensorReadingFilter) (*sharedDomain.PaginatedResponse, error) {
	return s.repo.GetByFilter(ctx, filter)
}

func (s *SensorService) UpdateReading(ctx context.Context, id int, reading *sharedDomain.SensorReading) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("sensor reading with id %d not found", id)
	}
	reading.ID = id
	return s.repo.Update(ctx, id, reading)
}

func (s *SensorService) DeleteReadings(ctx context.Context, filter *sharedDomain.SensorReadingFilter) (int64, error) {
	return s.repo.Delete(ctx, filter)
}

func (s *SensorService) GetReadingByID(ctx context.Context, id int) (*sharedDomain.SensorReading, error) {
	return s.repo.GetByID(ctx, id)
}