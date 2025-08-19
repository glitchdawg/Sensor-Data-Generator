package domain

import (
	"time"
)

type SensorReading struct {
	ID         int       `json:"id" db:"id" example:"1"`                                                        // Unique identifier
	ID1        string    `json:"id1" db:"id1" validate:"required,alpha,uppercase" example:"A"`                 // First identifier (A-Z)
	ID2        int       `json:"id2" db:"id2" validate:"required,min=0,max=999" example:"42"`                  // Second identifier (0-999)
	SensorType string    `json:"sensor_type" db:"sensor_type" validate:"required" example:"temperature"`       // Type of sensor
	Value      float64   `json:"value" db:"value" validate:"required" example:"23.5"`                          // Sensor reading value
	Timestamp  time.Time `json:"timestamp" db:"ts" example:"2024-01-15T10:30:00Z"`                             // When reading was taken
}

type SensorReadingFilter struct {
	ID1       *string
	ID2       *int
	From      *time.Time
	To        *time.Time
	Page      int
	PageSize  int
}

type PaginatedSensorReadings struct {
	Data       []SensorReading `json:"data"`
	Page       int             `json:"page" example:"1"`
	PageSize   int             `json:"page_size" example:"10"`
	TotalItems int64           `json:"total_items" example:"100"`
	TotalPages int             `json:"total_pages" example:"10"`
}
type PaginatedResponse = PaginatedSensorReadings