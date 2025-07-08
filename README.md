# Axis Backend Assessment

A Go application using Echo framework and MongoDB for financial transaction management with authentication.

## Prerequisites

- Go 1.16 or higher
- Docker and Docker Compose (for containerized setup)
- MongoDB (if running locally)

## Installation

1. Clone the repository

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Set up environment variables by copying the example file:

   ```bash
   cp .env.example .env
   ```

   Edit the `.env` file to configure your environment settings.

## Environment Variables

- `PORT`: Application port (default: "8080")
- `ENV`: Environment mode (development/production)
- `MONGO_URI`: MongoDB connection string (default: "mongodb://localhost:27017")
- `DB_NAME`: MongoDB database name
- `JWT_SECRET`: Secret key for JWT token generation
- `JWT_EXPIRATION`: JWT token expiration time

## Running with Docker Compose

The easiest way to run the application with all its dependencies is using Docker Compose:

```bash
docker-compose up -d
```

This will start the application and MongoDB in containers.

## Running the Application Locally

```bash
go run cmd/server/main.go
```

The server will start on [http://localhost:8080](http://localhost:8080)

## Running Tests

### Running All Tests

```bash
# Clear test cache first (recommended)
go clean -testcache

# Run all tests
go test ./tests/...
```

### Running Specific Test Suites

```bash
# Run auth handler tests
go test -v ./tests/handlers/auth_handler_test.go

# Run auth service tests
go test -v ./tests/services/auth_service_test.go

# Run transaction handler tests
go test -v ./tests/handlers/transaction_handler_test.go

# Run transaction service tests
go test -v ./tests/services/transaction_service_test.go
```

## API Documentation

Swagger documentation is available at `/swagger/index.html` when the server is running.

## Project Structure

The project follows clean architecture principles:

- `cmd/server`: Main application entry point
- `internal/`
  - `api/`: HTTP handlers, routes, middleware
  - `config/`: Application configuration
  - `dtos/`: Data transfer objects
  - `models/`: Domain models
  - `repository/`: Data access layer
  - `services/`: Business logic
- `pkg/`: Shared packages (database, jwt, logger)
- `tests/`: Test files and mocks
- `docs/`: Swagger documentation
