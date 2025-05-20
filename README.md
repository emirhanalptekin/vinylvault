# VinylVault

VinylVault is a simple REST API for managing your vinyl record collection. It provides basic CRUD operations for albums, artists, and genres.

## Features

- Create, read, update, and delete albums
- Simple organizational structure for artists and genres
- Docker containerization for easy deployment
- PostgreSQL database for data storage
- Comprehensive testing suite (unit and integration tests)

## Technology Stack

- **Backend**: Go with Gin framework
- **Database**: PostgreSQL
- **Containerization**: Docker and Docker Compose
- **Testing**: Go's built-in testing package

## Architecture

VinylVault follows a clean architecture approach with separation of concerns:

- **Handlers**: Process HTTP requests and responses
- **Models**: Define data structures
- **Database**: Manages data persistence
- **Config**: Handles application configuration

## Getting Started

### Prerequisites

- Go 1.19+
- Docker and Docker Compose
- PostgreSQL (if running locally without Docker)

### Running Locally with Docker

1. Clone the repository
2. Navigate to the project directory
3. Run the application with Docker Compose:

```bash
docker-compose up -d
```

The API will be available at http://localhost:8080

### Running Locally without Docker

1. Start a PostgreSQL instance
2. Update the configuration in `internal/config/config.yml`
3. Run the application:

```bash
go run cmd/main.go
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | /albums  | Get all albums |
| GET    | /albums/:id | Get album by ID |
| POST   | /albums  | Create a new album |
| PUT    | /albums/:id | Update an album |
| DELETE | /albums/:id | Delete an album |
| GET    | /artists | Get all artists |
| GET    | /genres  | Get all genres |

## Testing

### Unit Tests

```bash
go test -v ./internal/...
```

### Integration Tests

```bash
go test -v ./tests/...
```

## Project Structure

The project follows a standard Go project layout:

- `cmd/`: Application entry points
- `internal/`: Private application code
- `pkg/`: Public libraries
- `tests/`: Integration tests

## Documentation

API documentation is available at the `/swagger/index.html` endpoint when the server is running.