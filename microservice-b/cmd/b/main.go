//	@title			Synthetic Sensors API
//	@version		1.0
//	@description	API for managing sensor data readings with microservice architecture
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.email	joydeeppaul9000@gmail.com
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"

	pb "github.com/glitchdawg/synthetic_sensors/proto/ingestpb"
	"github.com/glitchdawg/synthetic_sensors/microservice-b/internal/handler"
	customMiddleware "github.com/glitchdawg/synthetic_sensors/microservice-b/internal/middleware"
	"github.com/glitchdawg/synthetic_sensors/microservice-b/internal/repository"
	"github.com/glitchdawg/synthetic_sensors/microservice-b/internal/service"
)

func main() {
	// Configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:root@postgres:5432/sensors?sslmode=disable"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	// Wait for database to be ready
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		log.Printf("waiting for database... (%d/30)", i+1)
		time.Sleep(1 * time.Second)
	}

	// Run database migrations
	if err := runMigrations(db); err != nil {
		log.Printf("migration failed: %v", err)
	}

	// Initialize layers
	repo := repository.NewPostgresRepository(db)
	sensorService := service.NewSensorService(repo)
	sensorHandler := handler.NewSensorHandler(sensorService)
	authHandler := handler.NewAuthHandler()
	grpcHandler := handler.NewGRPCHandler(sensorService)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatal("failed to listen:", err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterIngestServiceServer(grpcServer, grpcHandler)
		log.Printf("Microservice B gRPC server listening on :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("failed to serve gRPC:", err)
		}
	}()

	// Setup Echo HTTP server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(100)))

	// Public routes
	e.POST("/api/auth/login", authHandler.Login)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "healthy"})
	})
	
	// Swagger documentation
	e.Static("/docs", "./microservice-b/static")
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(302, "/docs/index.html")
	})

	// Protected routes
	api := e.Group("/api")
	api.Use(customMiddleware.JWTMiddleware)

	// Sensor readings endpoints
	api.GET("/readings", sensorHandler.GetReadings)
	api.GET("/readings/:id", sensorHandler.GetReadingByID)
	api.POST("/readings", sensorHandler.CreateReading, customMiddleware.RequireRole("admin"))
	api.PUT("/readings/:id", sensorHandler.UpdateReading, customMiddleware.RequireRole("admin"))
	api.DELETE("/readings", sensorHandler.DeleteReadings, customMiddleware.RequireRole("admin"))

	log.Printf("Microservice B REST API listening on :%s", httpPort)
	e.Logger.Fatal(e.Start(":" + httpPort))
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db",
		"postgres", driver)
	if err != nil {
		return err
	}

	return m.Up()
}
