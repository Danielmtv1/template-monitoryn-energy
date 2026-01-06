FROM golang:1.24 AS builder

WORKDIR /app

# Install build dependencies required by confluent-kafka-go (librdkafka)
RUN apt-get update && apt-get install -y --no-install-recommends \
    pkg-config \
    librdkafka-dev \
    build-essential \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build (enable CGO for librdkafka)
COPY . .
RUN CGO_ENABLED=1 \
    go build -v -ldflags "-s -w" -o /monitoring-energy-service main.go

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
 && rm -rf /var/lib/apt/lists/*
COPY --from=builder /monitoring-energy-service /monitoring-energy-service
WORKDIR /
EXPOSE 9000
ENTRYPOINT ["/monitoring-energy-service"]
