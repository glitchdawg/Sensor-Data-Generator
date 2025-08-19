package domain

import (
	"context"
	"github.com/glitchdawg/synthetic_sensors/shared/domain"
)

type SensorReadingRepository interface {
	Create(ctx context.Context, reading *domain.SensorReading) error
	GetByFilter(ctx context.Context, filter *domain.SensorReadingFilter) (*domain.PaginatedResponse, error)
	Update(ctx context.Context, id int, reading *domain.SensorReading) error
	Delete(ctx context.Context, filter *domain.SensorReadingFilter) (int64, error)
	GetByID(ctx context.Context, id int) (*domain.SensorReading, error)
}