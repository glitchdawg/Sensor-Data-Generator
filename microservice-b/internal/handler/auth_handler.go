package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/glitchdawg/synthetic_sensors/microservice-b/internal/middleware"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

// POST /api/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Simple hardcoded authentication for demo
	// In production, verify against database
	var userID, role string
	if req.Username == "admin" && req.Password == "admin123" {
		userID = "1"
		role = "admin"
	} else if req.Username == "user" && req.Password == "user123" {
		userID = "2"
		role = "user"
	} else {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	token, err := middleware.GenerateToken(userID, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		Type:  "Bearer",
	})
}