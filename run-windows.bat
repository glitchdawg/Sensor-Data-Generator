@echo off
echo Starting Synthetic Sensors Application...
echo.

REM Check if Docker is running
docker version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Docker is not running. Please start Docker Desktop.
    pause
    exit /b 1
)

REM Build and start services
echo Building Docker images...
docker-compose build

echo.
echo Starting services...
docker-compose up -d

echo.
echo Waiting for services to be ready...
timeout /t 10 /nobreak >nul

echo.
echo Services status:
docker-compose ps

echo.
echo Application URLs:
echo - Microservice B REST API: http://localhost:8080
echo - Microservice B gRPC: localhost:9090
echo - Microservice A (temperature): http://localhost:8081
echo - Microservice A (humidity): http://localhost:8082
echo - Microservice A (pressure): http://localhost:8083
echo - PostgreSQL: localhost:5432
echo.
echo Default credentials:
echo - Admin: username=admin, password=admin123
echo - User: username=user, password=user123
echo.
echo To stop the application, run: docker-compose down
echo To view logs, run: docker-compose logs -f
pause