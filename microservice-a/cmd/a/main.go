package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/glitchdawg/synthetic_sensors/proto/ingestpb"
)

var freq int64 = 1000

func main() {
	e := echo.New()

	//UPDATE FREQUENCY
	e.PUT("/config/frequency", func(c echo.Context) error {
		type Req struct {
			Frequency int64 `json:"frequency_ms"`
		}
		req := new(Req)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		atomic.StoreInt64(&freq, req.Frequency)
		return c.JSON(http.StatusOK, req)
	})

	// gRPC goroutine
	go runGenerator()

	log.Println("Microservice A running on :8081")
	e.Logger.Fatal(e.Start(":8081"))
}

func runGenerator() {
	conn, err := grpc.Dial("microservice-b:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect B: %v", err)
	}
	defer conn.Close()

	client := pb.NewIngestServiceClient(conn)
	stream, err := client.Write(context.Background())
	if err != nil {
		log.Fatalf("failed to open stream: %v", err)
	}

	sensorType := os.Getenv("SENSOR_TYPE")
	if sensorType == "" {
		sensorType = "temperature"
	}

	for {
		now := time.Now().UTC().Format(time.RFC3339Nano)
		msg := &pb.Reading{
			Value:      rand.Float64() * 100,
			SensorType: sensorType,
			Id1:        "A",
			Id2:        int32(rand.Intn(10)),
			Timestamp:  now,
		}
		if err := stream.Send(msg); err != nil {
			log.Printf("send error: %v", err)
		}
		time.Sleep(time.Duration(atomic.LoadInt64(&freq)) * time.Millisecond)
	}
}
