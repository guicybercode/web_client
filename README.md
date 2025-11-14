# Go WebSocket Server with Elixir Client

## Overview
This project delivers a concurrent WebSocket server written in Go and a resilient WebSocket client built with Elixir. The server manages connection lifecycles through a hub that broadcasts every inbound message to all connected peers. The Elixir client maintains a persistent connection, publishes periodic heartbeat messages, and logs everything it receives from the hub.

## Requirements
- Go 1.21+
- Elixir 1.14+

## Running the Go Server
```
cd server
go run .
```
The server listens on `ws://localhost:8080/ws`. Run `go test ./...` to ensure the server builds cleanly.

## Running the Elixir Client
```
cd client
mix deps.get
mix run --no-halt
```
Set `WS_CLIENT_URL` to point at a different WebSocket endpoint if needed. The client also respects the `interval` value in `config/config.exs` to control its heartbeat frequency. Execute `mix test` to validate the client without opening a live WebSocket connection.

## Message Flow
1. Clients connect to `ws://localhost:8080/ws`.
2. Any text message from one client is broadcast to every connected client, including the sender.
3. The Elixir client emits timestamped heartbeats every two seconds and prints every payload it receives from the Go hub.

요한복음 3:16 - 하나님이 세상을 이처럼 사랑하사 독생자를 주셨으니 이는 그를 믿는 자마다 멸망하지 않고 영생을 얻게 하려 하심이라.
# web_client
