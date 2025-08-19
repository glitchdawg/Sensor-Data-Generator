# Synthetic Sensors - Microservice Architecture

A scalable sensor data processing system built with Go, implementing microservice architecture with gRPC communication and REST APIs.

## 🏗️ Architecture Overview

The system consists of multiple microservices:

- **Microservice A (Data Generators)**: Multiple instances generating sensor data
- **Microservice B (Data Processor)**: Receives, processes, and stores sensor data
- **PostgreSQL Database**: Persistent storage for sensor readings

### Architecture Diagram

```
┌─────────────────────┐     ┌─────────────────────┐     ┌─────────────────────┐
│  Microservice A     │     │  Microservice A     │     │  Microservice A     │
│  (Temperature)      │     │  (Humidity)         │     │  (Pressure)         │
│  Port: 8081         │     │  Port: 8082         │     │  Port: 8083         │
└──────────┬──────────┘     └──────────┬──────────┘     └──────────┬──────────┘
           │                           │                           │
           │         gRPC Stream       │         gRPC Stream       │
           └───────────────────────────┼───────────────────────────┘
                                       │
                                       ▼
                          ┌────────────────────────┐
                          │   Microservice B       │
                          │   ┌────────────────┐   │
                          │   │  gRPC Server   │   │
                          │   │  Port: 9090    │   │
                          │   └────────┬───────┘   │
                          │            │           │
                          │   ┌────────▼───────┐   │
                          │   │ Service Layer  │   │
                          │   └────────┬───────┘   │
                          │            │           │
                          │   ┌────────▼───────┐   │
                          │   │  REST API      │   │
                          │   │  Port: 8080    │   │
                          │   └────────────────┘   │
                          └────────────┬───────────┘
                                       │
                                       ▼
                          ┌────────────────────────┐
                          │    PostgreSQL DB       │
                          │    Port: 5432          │
                          └────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Docker Desktop (Windows/Mac/Linux)
- Git

### Windows Users

1. Clone the repository:
```bash
git clone https://github.com/yourusername/synthetic_sensors.git
cd synthetic_sensors
```

2. Run the application:
```bash
# Using Command Prompt
run-windows.bat

# Using PowerShell
.\run-windows.ps1
```

### Linux/Mac Users

```bash
docker-compose up -d
```

## 📋 Features

### Microservice A (Data Generator)
- ✅ Generates sensor data with configurable frequency
- ✅ Supports multiple sensor types (temperature, humidity, pressure)
- ✅ REST API endpoint to change data generation frequency
- ✅ Sends data via gRPC stream to Microservice B

### Microservice B (Data Processor)
- ✅ Receives sensor data via gRPC
- ✅ Stores data in PostgreSQL database
- ✅ REST API with full CRUD operations
- ✅ JWT-based authentication and authorization
- ✅ Pagination support for data retrieval
- ✅ Advanced filtering (by ID combination, time range)
- ✅ Rate limiting protection

## 🔐 Authentication

The system uses JWT tokens for authentication. Default credentials:

- **Admin**: username=`admin`, password=`admin123`
- **User**: username=`user`, password=`user123`

### Login

```bash
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

## 📡 API Endpoints

### Authentication
- `POST /api/auth/login` - Login and get JWT token

### Sensor Readings (Protected)
- `GET /api/readings` - Get sensor readings with pagination and filters
- `GET /api/readings/:id` - Get specific reading by ID
- `POST /api/readings` - Create new reading (Admin only)
- `PUT /api/readings/:id` - Update reading (Admin only)
- `DELETE /api/readings` - Delete readings by filter (Admin only)

### Configuration (Microservice A)
- `PUT /config/frequency` - Update data generation frequency
- `GET /config/frequency` - Get current frequency
- `GET /health` - Health check

### Query Parameters for GET /api/readings
- `id1` - Filter by ID1 (e.g., "A", "B", "C")
- `id2` - Filter by ID2 (integer)
- `from` - Start timestamp (RFC3339 format)
- `to` - End timestamp (RFC3339 format)
- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 10, max: 100)

## 🏛️ Clean Architecture

The project follows Clean Architecture principles:

```
microservice-a/
├── cmd/a/              # Application entry point
├── internal/
│   ├── domain/         # Business entities
│   ├── handler/        # HTTP handlers
│   └── service/        # Business logic

microservice-b/
├── cmd/b/              # Application entry point
├── internal/
│   ├── domain/         # Business entities & interfaces
│   ├── handler/        # HTTP & gRPC handlers
│   ├── middleware/     # Authentication middleware
│   ├── repository/     # Data persistence layer
│   └── service/        # Business logic

shared/
└── domain/             # Shared domain models
```

## 🗃️ Database Schema

```sql
CREATE TABLE sensor_readings (
    id SERIAL PRIMARY KEY,
    id1 VARCHAR(10) NOT NULL,
    id2 INT NOT NULL,
    sensor_type VARCHAR(50) NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    ts TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_sensor_id1_id2_ts ON sensor_readings (id1, id2, ts);
```

## 🔧 Configuration

Environment variables can be configured in `docker-compose.yml`:

### Microservice A
- `SENSOR_TYPE` - Type of sensor (temperature, humidity, pressure)
- `GRPC_ADDRESS` - Address of Microservice B gRPC server
- `PORT` - HTTP server port

### Microservice B
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT tokens
- `GRPC_PORT` - gRPC server port
- `HTTP_PORT` - HTTP REST API port

## 📊 Monitoring

View logs:
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f microservice-b
```

Check service status:
```bash
docker-compose ps
```

## 🧪 Testing

Use the provided Postman collection (`postman_collection.json`) to test all API endpoints.

## 🚦 Scaling

To add more sensor instances, add new services in `docker-compose.yml`:

```yaml
microservice-a-custom:
  build:
    context: .
    dockerfile: Dockerfile.microservice-a
  environment:
    SENSOR_TYPE: custom_sensor
    GRPC_ADDRESS: microservice-b:9090
    PORT: 8084
  ports:
    - "8084:8084"
```

## 🛑 Stopping the Application

```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## 📝 License

This project is created for educational purposes as part of a job interview assignment.