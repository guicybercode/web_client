# Go WebSocket Server with Elixir Client

## Overview
This project delivers a concurrent WebSocket server written in Go and a resilient WebSocket client built with Elixir. The server manages connection lifecycles through a hub that broadcasts every inbound message to all connected peers. The Elixir client maintains a persistent connection, publishes periodic heartbeat messages, and logs everything it receives from the hub.

Messages are enriched with metadata (client ID, timestamp) to prevent echo and provide better observability.

## Requirements
- Go 1.21+
- Elixir 1.14+
- Docker & Docker Compose (optional)

## Running Locally

### Go Server
```
cd server
go run . --port 8080
```
The server listens on `ws://localhost:8080/ws` by default. Run `go test ./...` to ensure the server builds cleanly.

**Configuration:**
- `--port`: Port to listen on (default 8080)
- `--host`: Host to bind to (default empty)
- `PORT` environment variable
- `HOST` environment variable

### Elixir Client
```
cd client
mix deps.get
mix run --no-halt
```
Set `WS_CLIENT_URL` to point at a different WebSocket endpoint if needed. The client also respects the `interval` value in `config/config.exs` to control its heartbeat frequency. Execute `mix test` to validate the client without opening a live WebSocket connection.

**Configuration:**
- `WS_CLIENT_URL`: WebSocket server URL (default `ws://localhost:8080/ws`)
- `interval`: Heartbeat interval in ms (default 2000)
- `reconnect_delay`: Delay between reconnection attempts in ms (default 1000)

### Running with Docker Compose
```
docker-compose up --build
```
This will start both the Go server and Elixir client in containers, with automatic service discovery.

## Features

### Message Flow
1. Clients connect to `ws://localhost:8080/ws`.
2. Any text message from one client is broadcast to every **other** connected client (no echo).
3. Messages include metadata: `{"client_id": "...", "timestamp": 1234567890, "content": "message text"}`
4. The Elixir client emits timestamped heartbeats every two seconds and prints received messages with sender ID.

### Logging
Both server and client emit structured JSON logs for monitoring:
- Client connections/disconnections with IDs
- Message broadcasts with counts
- Errors and reconnection attempts

### Resilience
- Elixir client auto-reconnects on disconnection
- Multiple client instances can run simultaneously
- Asynchronous startup prevents boot failures

## Development

### Testing
```
# Go server
cd server && go test -v

# Elixir client
cd client && mix test
```

### Building Docker Images
```
# Server
docker build -f Dockerfile.server -t ws-server .

# Client
docker build -f Dockerfile.client -t ws-client .
```

요한복음 3:16 - 하나님이 세상을 이처럼 사랑하사 독생자를 주셨으니 이는 그를 믿는 자마다 멸망하지 않고 영생을 얻게 하려 하심이라.
