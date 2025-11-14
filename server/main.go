package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var logger *slog.Logger

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	logger = slog.New(handler)
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("websocket upgrade failed", "error", err)
		return
	}
	client := NewClient(hub, conn)
	hub.register <- client
	logger.Info("client connected", "client_id", client.id, "remote_addr", r.RemoteAddr)
	go client.writePump()
	go client.readPump()
}

func main() {
	var port string
	var host string

	flag.StringVar(&port, "port", "", "Port to listen on (default 8080)")
	flag.StringVar(&host, "host", "", "Host to bind to (default empty)")
	flag.Parse()

	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}
	if _, err := strconv.Atoi(port); err != nil {
		logger.Error("invalid port", "port", port, "error", err)
		os.Exit(1)
	}

	if host == "" {
		host = os.Getenv("HOST")
	}

	hub := NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	addr := host + ":" + port
	logger.Info("server listening", "addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
