package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	pb "github.com/glitchdawg/synthetic_sensors/proto/ingestpb"
)

var db *sql.DB

type server struct {
	pb.UnimplementedIngestServiceServer
}

func (s *server) Write(stream pb.IngestService_WriteServer) error {
	count := uint64(0)
	for {
		reading, err := stream.Recv()
		if err != nil {
			return stream.SendAndClose(&pb.WriteAck{Count: count})
		}
		_, err = db.Exec(`INSERT INTO sensor_readings (id1, id2, sensor_type, value, ts) VALUES (?, ?, ?, ?, ?)`,
			reading.Id1, reading.Id2, reading.SensorType, reading.Value, reading.Timestamp)
		if err != nil {
			log.Println("insert error:", err)
		}
		count++
	}
}

func main() {
	// Connect to PostgreSQL database
	var err error
	db, err = sql.Open("postgres", "user=postgres password=root dbname=sensors sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatal(err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterIngestServiceServer(grpcServer, &server{})
		log.Println("Microservice B gRPC on :9090")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	// REST API
	e := echo.New()

	e.GET("/readings", func(c echo.Context) error {
		id1 := c.QueryParam("id1")
		id2 := c.QueryParam("id2")
		from := c.QueryParam("from")
		to := c.QueryParam("to")
		limit := c.QueryParam("limit")
		if limit == "" {
			limit = "100"
		}

		q := `SELECT id, id1, id2, sensor_type, value, ts 
      FROM sensor_readings WHERE 1=1`
		args := []interface{}{}
		i := 1

		if id1 != "" {
			q += fmt.Sprintf(" AND id1=$%d", i)
			args = append(args, id1)
			i++
		}
		if id2 != "" {
			q += fmt.Sprintf(" AND id2=$%d", i)
			args = append(args, id2)
			i++
		}
		if from != "" {
			q += fmt.Sprintf(" AND ts >= $%d", i)
			args = append(args, from)
			i++
		}
		if to != "" {
			q += fmt.Sprintf(" AND ts <= $%d", i)
			args = append(args, to)
			i++
		}

		q += fmt.Sprintf(" ORDER BY ts DESC LIMIT %s", limit)

		rows, err := db.Query(q, args...)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer rows.Close()

		type Reading struct {
			Id1        string    `json:"id1"`
			Id2        int       `json:"id2"`
			SensorType string    `json:"sensor_type"`
			Value      float64   `json:"value"`
			Ts         time.Time `json:"timestamp"`
		}
		readings := []Reading{}
		for rows.Next() {
			var r Reading
			if err := rows.Scan(&r.Id1, &r.Id2, &r.SensorType, &r.Value, &r.Ts); err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			readings = append(readings, r)
		}
		return c.JSON(http.StatusOK, readings)
	})

	log.Println("Microservice B REST on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
