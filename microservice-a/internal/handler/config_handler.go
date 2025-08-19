package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/glitchdawg/synthetic_sensors/microservice-a/internal/service"
)

type ConfigHandler struct {
	service   *service.GeneratorService
	validator *validator.Validate
}

func NewConfigHandler(service *service.GeneratorService) *ConfigHandler {
	return &ConfigHandler{
		service:   service,
		validator: validator.New(),
	}
}

type UpdateFrequencyRequest struct {
	FrequencyMs int64 `json:"frequency_ms" validate:"required,min=100"`
}

// PUT /config/frequency
func (h *ConfigHandler) UpdateFrequency(c echo.Context) error {
	req := new(UpdateFrequencyRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	h.service.UpdateFrequency(req.FrequencyMs)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"frequency_ms": req.FrequencyMs,
		"message":      "frequency updated successfully",
	})
}

// GET /config/frequency
func (h *ConfigHandler) GetFrequency(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"frequency_ms": h.service.GetFrequency(),
	})
}