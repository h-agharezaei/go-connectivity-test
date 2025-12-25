# Go Connectivity Test

A simple Go application to test connectivity to various databases and message brokers, including PostgreSQL, Redis, MongoDB, NATS, NATS JetStream, and Kafka.

## Features

- **PostgreSQL**: Tests connection and ping.
- **Redis**: Tests connection and ping.
- **MongoDB**: Tests connection and ping.
- **NATS**: Tests connection and publishes a message.
- **NATS JetStream**: Tests JetStream functionality, including stream creation, publishing, and consuming.
- **Kafka**: Tests connection to the broker.

## Prerequisites

- Go 1.25.2 or later
- Running instances of the following services (or adjust `.env` accordingly):
  - PostgreSQL
  - Redis
  - MongoDB
  - NATS Server
  - Kafka

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/h-agharezaei/go-connectivity-test.git
   cd go-connectivity-test
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Copy the example environment file and configure it:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` with your actual service credentials and ports.

## Usage

Build and run the application:

```bash
make build
make run
```

Or manually:

```bash
go build -o connectivity-test ./cmd
./connectivity-test
```

The application will log the connectivity status for each service.

## Configuration

Environment variables in `.env`:

- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_HOST`, `POSTGRES_PORT`
- `REDIS_HOST`, `REDIS_PORT`
- `MONGO_URI`
- `NATS_URL`, `NATS_JS_URL`
- `KAFKA_HOST`, `KAFKA_PORT`

## Makefile Targets

- `make build`: Build the application
- `make run`: Run the application
- `make clean`: Remove build artifacts
- `make test`: Run tests (if any)

## Contributing

Feel free to submit issues or pull requests.

## License

This project is licensed under the MIT License.