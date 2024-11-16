# Ticket Management API

## Overview
A RESTful API service built with Go for managing ticket allocations and purchases. The application provides endpoints for creating tickets, retrieving ticket information, and processing ticket purchases. It uses PostgreSQL for data persistence and Docker for containerization.

## API Endpoints

### Tickets
- `GET /tickets/{id}` - Retrieve ticket details by ID
- `POST /tickets/{id}/purchases` - Purchase tickets
- `POST /ticketsuser` - Create a new ticket

For detailed API documentation, visit `/swagger/index.html` after starting the application.

## Project Structure
```
project/
├── cmd/                      # Application entry points
├── controller/               # API Controllers
├── docs/                     # API documentation
├── domain/                   # Domain models and interfaces
├── infrastructure/           # Infrastructure Implementations (Database, Cache, etc.)
├── interface/                # Interface Implementations (HTTP, GRPC, etc.)
├── mock/                     # Test mocks
├── pkg/                      # Shared packages
├── service/                  # Business logic layer
├── valueobject/              # Value objects
```

## How to Run

### Prerequisites
- Go 1.22 or higher
- Docker and Docker Compose
- PostgreSQL (if running locally)

### Environment Variables
Create a `.env` file using the `.env.example` file as a template. (Copy the `.env` file under `cmd/` folder if you want to run the application locally)

### Testing
Run `go generate ./...` to generate the mocks before running any tests.

#### Unit tests:
```bash
go test ./...
```

#### Integration tests
```bash
INTEGRATION=true go test ./... #unix
$env:INTEGRATION="true"; go test ./... #powershell
```
> :warning: **Integration tests use Docker, you need to start docker before running any integration test otherwise test will fail**


### Docker Setup (Recommended)
1. Build and start the services:
```bash
docker-compose up --build
```

2. The API will be available at `http://localhost:{PORT}`
3. Swagger documentation will be available at `http://localhost:{PORT}/swagger/index.html`

### Local Development Setup
1. Install dependencies:
```bash
go mod download
```

2. Generate Swagger documentation (If any changes are made to the API):
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g ./controller/ticket/ticket.go --parseDependency true
```

3. Start the application:
```bash
go run cmd/main.go
```

## Example API Calls

### Create a New Ticket
```bash
curl -X POST 'http://localhost:8080/ticketsuser' \
-H 'Content-Type: application/json' \
-d '{
    "name": "Concert Ticket",
    "description": "VIP Concert Access",
    "allocation": 100
}'
```

### Get Ticket by ID
```bash
curl -X GET 'http://localhost:8080/tickets/1' \
-H 'Content-Type: application/json'
```

### Purchase Tickets
```bash
curl -X POST 'http://localhost:8080/tickets/1/purchases' \
-H 'Content-Type: application/json' \
-d '{
    "quantity": 2,
    "user_id": "f64e1422-f67e-4629-af14-111f85ac6655"
}'
```


## Example Responses

### Successful Ticket Creation
```json
{
    "id": 1,
    "name": "Concert Ticket",
    "description": "VIP Concert Access",
    "allocation": 100
}
```

### Successful Ticket Retrieval
```json
{
    "id": 1,
    "name": "Concert Ticket",
    "description": "VIP Concert Access",
    "allocation": 100
}
```

### Error Response
```json
{
    "status": 400,
    "message": "Invalid request",
    "errors": [
        {
            "failed_field": "quantity",
            "tag": "min",
            "message": "quantity must be greater than 0"
        }
    ]
}
```

## Error Handling
The API returns standardized error responses with appropriate HTTP status codes:
- 400: Bad Request
- 404: Not Found
- 422: Unprocessable Entity
- 500: Internal Server Error
