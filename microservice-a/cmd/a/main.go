package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/glitchdawg/synthetic_sensors/microservice-a/internal/domain"
	"github.com/glitchdawg/synthetic_sensors/microservice-a/internal/handler"
	"github.com/glitchdawg/synthetic_sensors/microservice-a/internal/service"
)

func main() {
	// Configuration
	sensorType := os.Getenv("SENSOR_TYPE")
	if sensorType == "" {
		sensorType = "temperature"
	}

	grpcAddr := os.Getenv("GRPC_ADDRESS")
	if grpcAddr == "" {
		grpcAddr = "microservice-b:9090"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	config := &domain.GeneratorConfig{
		FrequencyMs: 1000,
		SensorType:  sensorType,
	}

	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to microservice-b: %v", err)
	}
	defer conn.Close()

	// Initialize service and handler
	genService := service.NewGeneratorService(conn, config)
	configHandler := handler.NewConfigHandler(genService)

	// Start generator in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			if err := genService.StartGenerator(ctx); err != nil {
				log.Printf("generator error: %v, retrying in 5 seconds", err)
				time.Sleep(5 * time.Second)
			}
		}
	}()

	// Setup Echo server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.PUT("/config/frequency", configHandler.UpdateFrequency)
	e.GET("/config/frequency", configHandler.GetFrequency)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "healthy",
			"sensor_type": sensorType,
		})
	})

	log.Printf("Microservice A (sensor: %s) running on :%s", sensorType, port)
	e.Logger.Fatal(e.Start(":" + port))
}
