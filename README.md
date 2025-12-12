# Monitoring Energy Service

Go microservice base template with hexagonal architecture.

## Tech Stack

- **Go 1.24** - Programming language
- **Gin** - HTTP REST framework
- **GORM** - PostgreSQL ORM
- **Kafka** - Message broker (confluent-kafka-go)
- **PostgreSQL 15** - Database with extensions:
  - PostGIS 3.6
  - TimescaleDB 2.24
- **Goose** - Database migrations
- **Atlas** - Migration generation from GORM

## Project Structure

```
├── cmd/
│   └── atlasloader/          # Entity loader for Atlas
├── db/
│   ├── Dockerfile            # PostgreSQL image with extensions
│   └── init-db.sql           # Initialization script
├── internal/
│   ├── api/                  # Application services
│   │   └── kafka_service.go
│   ├── domain/
│   │   ├── entities/         # Domain entities
│   │   └── ports/
│   │       ├── input/        # Service interfaces
│   │       └── output/       # Repository/adapter interfaces
│   └── infrastructure/
│       ├── adapters/
│       │   ├── http/webhook/ # Webhook adapter
│       │   ├── kafka/        # Kafka adapter
│       │   ├── repositories/ # GORM repositories
│       │   └── rest/         # Gin router and handlers
│       ├── conf/
│       │   ├── database/     # PostgreSQL configuration
│       │   └── kafkaconf/    # Kafka factory
│       └── container/        # Dependency injection
├── migrations/               # SQL migrations (Goose)
├── .air.toml                 # Hot reload (Air)
├── modd.conf                 # File watcher (Modd)
├── atlas.hcl                 # Atlas configuration
├── docker-compose.yml        # Development services
└── main.go                   # Entry point
```

## Local Development

### Requirements

- Go 1.24+
- Docker and Docker Compose
- Make

### Install dev tools

```bash
make install-dev-tools
```

This installs:
- Air (hot reload)
- Modd (file watcher)
- Goose (migrations)
- Atlas (migration generation)
- Swag (Swagger documentation)

### Docker Services

```bash
# Start services
make docker-up

# Stop services
make docker-down
```

| Service | Port | URL/Connection |
|---------|------|----------------|
| PostgreSQL | 5432 | `postgres://postgres:postgres@localhost:5432/monitoring_energy` |
| pgAdmin | 5050 | http://localhost:5050 |
| Kafka | 9092 | `localhost:9092` |
| Kafka UI | 8080 | http://localhost:8080 |
| Zookeeper | 2181 | `localhost:2181` |

### Credentials

| Service | User | Password |
|---------|------|----------|
| PostgreSQL | postgres | postgres |
| pgAdmin | admin@admin.com | admin |

### PostgreSQL Extensions

The database includes the following extensions:

| Extension | Version | Description |
|-----------|---------|-------------|
| pgcrypto | 1.3 | Cryptographic functions |
| PostGIS | 3.6.1 | Geospatial data |
| TimescaleDB | 2.24.0 | Time series |

### Run the application

```bash
# Configure environment variables
cp .env.example .env

# Development with hot reload (Air)
make dev

# Development with Modd
make dev-modd

# Run directly
make run
```

### Build

```bash
make build
```

## Migrations

### Create new migration

```bash
make migrate-create name=migration_name
```

### Apply migrations

```bash
make goose-up
```

Migrations run automatically on application startup.

### Rollback

```bash
# Last migration
make goose-down

# To specific version
make goose-down-to version=20241211000000
```

### Migration status

```bash
make goose-status
```

## REST API

### Example endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/examples | List all |
| GET | /api/v1/examples/:id | Get by ID |
| POST | /api/v1/examples | Create new |
| PUT | /api/v1/examples/:id | Update |
| DELETE | /api/v1/examples/:id | Delete |

### Health checks

| Endpoint | Description |
|----------|-------------|
| GET /healthz | Liveness probe |
| GET /readyz | Readiness probe |

### Swagger Documentation

Swagger UI is automatically available in development mode:
- URL: http://localhost:9000/swagger/index.html

Swagger docs are regenerated automatically when running `make dev`, `make dev-modd`, or `make run`.

## Environment Variables

```bash
# Server
PORT=9000
ENVIRONMENT=dev

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=monitoring_energy
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_SCHEMA=public

# Kafka
LIST_KAFKA_BROKERS=localhost:9092
CONSUMER_GROUP=monitoring-energy-group
CONSUMER_TOPIC=events.default
PRODUCER_TOPIC=events.output

# Webhook
WEBHOOK_ENABLED=false
WEBHOOK_URL=

# HTTP
HTTP_CLIENT_TIMEOUT=30

# CORS
ALLOWED_CORS_SUFFIXES=.spotcloud.io
```

## Hexagonal Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        INFRASTRUCTURE                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   REST API  │  │    Kafka    │  │      Webhook        │  │
│  │   (Gin)     │  │   Adapter   │  │      Adapter        │  │
│  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘  │
│         │                │                     │             │
│  ┌──────▼────────────────▼─────────────────────▼──────────┐ │
│  │                     PORTS (Input)                       │ │
│  │               Service Interfaces                        │ │
│  └──────────────────────┬──────────────────────────────────┘ │
│                         │                                    │
│  ┌──────────────────────▼──────────────────────────────────┐ │
│  │                      DOMAIN                              │ │
│  │              Entities + Business Logic                   │ │
│  └──────────────────────┬──────────────────────────────────┘ │
│                         │                                    │
│  ┌──────────────────────▼──────────────────────────────────┐ │
│  │                    PORTS (Output)                        │ │
│  │              Repository Interfaces                       │ │
│  └──────┬───────────────┬───────────────────┬──────────────┘ │
│         │               │                   │                │
│  ┌──────▼──────┐ ┌──────▼──────┐ ┌──────────▼─────────────┐  │
│  │ PostgreSQL  │ │    Kafka    │ │      External APIs     │  │
│  │   (GORM)    │ │  Producer   │ │       (Webhook)        │  │
│  └─────────────┘ └─────────────┘ └────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Make Commands

| Command | Description |
|---------|-------------|
| `make dev` | Development with Air (hot reload) |
| `make dev-modd` | Development with Modd |
| `make run` | Run application |
| `make build` | Build binary |
| `make test` | Run tests |
| `make docker-up` | Start Docker services |
| `make docker-down` | Stop Docker services |
| `make migrate-create name=X` | Create migration |
| `make goose-up` | Apply migrations |
| `make goose-down` | Rollback last migration |
| `make goose-status` | Migration status |
| `make install-dev-tools` | Install dev tools |
| `make swagger` | Generate Swagger docs manually |
