package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/glitchdawg/synthetic_sensors/microservice-b/internal/service"
	"github.com/glitchdawg/synthetic_sensors/shared/domain"
)

type SensorHandler struct {
	service   *service.SensorService
	validator *validator.Validate
}

func NewSensorHandler(service *service.SensorService) *SensorHandler {
	return &SensorHandler{
		service:   service,
		validator: validator.New(),
	}
}

// GET /api/readings
func (h *SensorHandler) GetReadings(c echo.Context) error {
	filter := &domain.SensorReadingFilter{}

	if id1 := c.QueryParam("id1"); id1 != "" {
		filter.ID1 = &id1
	}
	if id2Str := c.QueryParam("id2"); id2Str != "" {
		id2, err := strconv.Atoi(id2Str)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id2 format"})
		}
		filter.ID2 = &id2
	}
	if fromStr := c.QueryParam("from"); fromStr != "" {
		from, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid from date format"})
		}
		filter.From = &from
	}
	if toStr := c.QueryParam("to"); toStr != "" {
		to, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid to date format"})
		}
		filter.To = &to
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	filter.Page = page
	filter.PageSize = pageSize

	result, err := h.service.GetReadings(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// GET /api/readings/:id
func (h *SensorHandler) GetReadingByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}

	reading, err := h.service.GetReadingByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if reading == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "reading not found"})
	}

	return c.JSON(http.StatusOK, reading)
}

// POST /api/readings
func (h *SensorHandler) CreateReading(c echo.Context) error {
	reading := &domain.SensorReading{}
	if err := c.Bind(reading); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validator.Struct(reading); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.service.CreateReading(c.Request().Context(), reading); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, reading)
}

// PUT /api/readings/:id
func (h *SensorHandler) UpdateReading(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id format"})
	}

	reading := &domain.SensorReading{}
	if err := c.Bind(reading); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.validator.Struct(reading); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.service.UpdateReading(c.Request().Context(), id, reading); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "reading updated successfully"})
}

// DELETE /api/readings
func (h *SensorHandler) DeleteReadings(c echo.Context) error {
	filter := &domain.SensorReadingFilter{}

	if id1 := c.QueryParam("id1"); id1 != "" {
		filter.ID1 = &id1
	}
	if id2Str := c.QueryParam("id2"); id2Str != "" {
		id2, err := strconv.Atoi(id2Str)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id2 format"})
		}
		filter.ID2 = &id2
	}
	if fromStr := c.QueryParam("from"); fromStr != "" {
		from, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid from date format"})
		}
		filter.From = &from
	}
	if toStr := c.QueryParam("to"); toStr != "" {
		to, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid to date format"})
		}
		filter.To = &to
	}

	count, err := h.service.DeleteReadings(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "readings deleted successfully",
		"count":   count,
	})
}