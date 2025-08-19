package domain

import (
	"time"
)

type SensorReading struct {
	ID         int       `json:"id" db:"id"`
	ID1        string    `json:"id1" db:"id1" validate:"required,alpha,uppercase"`
	ID2        int       `json:"id2" db:"id2" validate:"required,min=0,max=999"`
	SensorType string    `json:"sensor_type" db:"sensor_type" validate:"required"`
	Value      float64   `json:"value" db:"value" validate:"required"`
	Timestamp  time.Time `json:"timestamp" db:"ts"`
}

type SensorReadingFilter struct {
	ID1       *string
	ID2       *int
	From      *time.Time
	To        *time.Time
	Page      int
	PageSize  int
}

type PaginatedResponse struct {
	Data       []SensorReading `json:"data"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalItems int64           `json:"total_items"`
	TotalPages int             `json:"total_pages"`
}