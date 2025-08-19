package domain

type GeneratorConfig struct {
	FrequencyMs int64  `json:"frequency_ms" validate:"required,min=100"`
	SensorType  string `json:"sensor_type"`
}