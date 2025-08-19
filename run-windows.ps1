Write-Host "Starting Synthetic Sensors Application..." -ForegroundColor Green
Write-Host ""

# Check if Docker is running
try {
    docker version | Out-Null
} catch {
    Write-Host "ERROR: Docker is not running. Please start Docker Desktop." -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# Build and start services
Write-Host "Building Docker images..." -ForegroundColor Yellow
docker-compose build

Write-Host ""
Write-Host "Starting services..." -ForegroundColor Yellow
docker-compose up -d

Write-Host ""
Write-Host "Waiting for services to be ready..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

Write-Host ""
Write-Host "Services status:" -ForegroundColor Cyan
docker-compose ps

Write-Host ""
Write-Host "Application URLs:" -ForegroundColor Green
Write-Host "- Microservice B REST API: http://localhost:8080"
Write-Host "- Microservice B gRPC: localhost:9090"
Write-Host "- Microservice A (temperature): http://localhost:8081"
Write-Host "- Microservice A (humidity): http://localhost:8082"
Write-Host "- Microservice A (pressure): http://localhost:8083"
Write-Host "- PostgreSQL: localhost:5432"

Write-Host ""
Write-Host "Default credentials:" -ForegroundColor Yellow
Write-Host "- Admin: username=admin, password=admin123"
Write-Host "- User: username=user, password=user123"

Write-Host ""
Write-Host "To stop the application, run: docker-compose down" -ForegroundColor Cyan
Write-Host "To view logs, run: docker-compose logs -f" -ForegroundColor Cyan

Read-Host "Press Enter to continue"