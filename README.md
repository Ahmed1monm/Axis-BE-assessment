# Axis Backend Assessment

A Go application using Echo framework and MongoDB.

## Prerequisites

- Go 1.16 or higher
- MongoDB

## Installation

1. Clone the repository

2. Install dependencies:

```bash
go mod tidy
```

## Environment Variables

- `MONGO_URI`: MongoDB connection string (default: "mongodb://localhost:27017")
- `PORT`: Application port (default: "8080")

## Running the Application

```bash
go run main.go
```

The server will start on [http://localhost:8080](http://localhost:8080)
