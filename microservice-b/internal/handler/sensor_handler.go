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

//	@Summary		Get sensor readings
//	@Description	Get sensor readings with optional filtering by ID1, ID2, time range, and pagination
//	@Tags			Sensor Readings
//	@Accept			json
//	@Produce		json
//	@Param			id1			query		string	false	"Filter by ID1 (A-Z)"
//	@Param			id2			query		int		false	"Filter by ID2 (0-999)"
//	@Param			from		query		string	false	"Start timestamp (RFC3339 format)"
//	@Param			to			query		string	false	"End timestamp (RFC3339 format)"
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			page_size	query		int		false	"Items per page (default: 10, max: 100)"
//	@Success		200			{object}	domain.PaginatedSensorReadings	"Successfully retrieved readings"
//	@Failure		400			{object}	map[string]string				"Invalid request parameters"
//	@Failure		401			{object}	map[string]string				"Unauthorized"
//	@Failure		500			{object}	map[string]string				"Internal server error"
//	@Security		Bearer
//	@Router			/api/readings [get]
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

//	@Summary		Get sensor reading by ID
//	@Description	Get a specific sensor reading by its ID
//	@Tags			Sensor Readings
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Reading ID"
//	@Success		200	{object}	domain.SensorReading	"Successfully retrieved reading"
//	@Failure		400	{object}	map[string]string		"Invalid ID format"
//	@Failure		401	{object}	map[string]string		"Unauthorized"
//	@Failure		404	{object}	map[string]string		"Reading not found"
//	@Failure		500	{object}	map[string]string		"Internal server error"
//	@Security		Bearer
//	@Router			/api/readings/{id} [get]
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

//	@Summary		Create sensor reading
//	@Description	Create a new sensor reading (requires admin privileges)
//	@Tags			Sensor Readings
//	@Accept			json
//	@Produce		json
//	@Param			reading	body		domain.SensorReading	true	"Sensor reading data"
//	@Success		201		{object}	domain.SensorReading	"Reading created successfully"
//	@Failure		400		{object}	map[string]string		"Invalid request body or validation error"
//	@Failure		401		{object}	map[string]string		"Unauthorized"
//	@Failure		403		{object}	map[string]string		"Forbidden - Admin access required"
//	@Failure		500		{object}	map[string]string		"Internal server error"
//	@Security		Bearer
//	@Router			/api/readings [post]
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

//	@Summary		Update sensor reading
//	@Description	Update an existing sensor reading by ID (requires admin privileges)
//	@Tags			Sensor Readings
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Reading ID"
//	@Param			reading	body		domain.SensorReading	true	"Updated sensor reading data"
//	@Success		200		{object}	map[string]string		"Reading updated successfully"
//	@Failure		400		{object}	map[string]string		"Invalid ID format or request body"
//	@Failure		401		{object}	map[string]string		"Unauthorized"
//	@Failure		403		{object}	map[string]string		"Forbidden - Admin access required"
//	@Failure		404		{object}	map[string]string		"Reading not found"
//	@Failure		500		{object}	map[string]string		"Internal server error"
//	@Security		Bearer
//	@Router			/api/readings/{id} [put]
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

//	@Summary		Delete sensor readings
//	@Description	Delete sensor readings based on filter criteria (requires admin privileges)
//	@Tags			Sensor Readings
//	@Accept			json
//	@Produce		json
//	@Param			id1		query		string	false	"Filter by ID1 (A-Z)"
//	@Param			id2		query		int		false	"Filter by ID2 (0-999)"
//	@Param			from	query		string	false	"Start timestamp (RFC3339 format)"
//	@Param			to		query		string	false	"End timestamp (RFC3339 format)"
//	@Success		200		{object}	map[string]interface{}	"Readings deleted successfully with count"
//	@Failure		400		{object}	map[string]string		"Invalid request parameters"
//	@Failure		401		{object}	map[string]string		"Unauthorized"
//	@Failure		403		{object}	map[string]string		"Forbidden - Admin access required"
//	@Failure		500		{object}	map[string]string		"Internal server error"
//	@Security		Bearer
//	@Router			/api/readings [delete]
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