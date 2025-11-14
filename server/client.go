package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type Client struct {
	id   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	id := generateID()
	return &Client{
		id:   id,
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			logger.Info("client read error", "client_id", c.id, "error", err)
			break
		}
		msg := NewMessage(c.id, string(rawMessage))
		logger.Info("message received", "client_id", c.id, "length", len(rawMessage))
		c.hub.broadcast <- msg.ToJSON()
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case rawMessage, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			msg, err := MessageFromJSON(rawMessage)
			if err != nil {
				logger.Info("invalid message format", "error", err)
				continue
			}
			if msg.ClientID == c.id {
				continue
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Info("write error", "client_id", c.id, "error", err)
				return
			}
			if _, err := w.Write(rawMessage); err != nil {
				logger.Info("write error", "client_id", c.id, "error", err)
				return
			}
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err := w.Write([]byte{'\n'}); err != nil {
					logger.Info("write error", "client_id", c.id, "error", err)
					return
				}
				nextRaw := <-c.send
				nextMsg, err := MessageFromJSON(nextRaw)
				if err != nil {
					logger.Info("invalid message format", "error", err)
					continue
				}
				if nextMsg.ClientID == c.id {
					continue
				}
				if _, err := w.Write(nextRaw); err != nil {
					logger.Info("write error", "client_id", c.id, "error", err)
					return
				}
			}
			if err := w.Close(); err != nil {
				logger.Info("write error", "client_id", c.id, "error", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Info("ping error", "client_id", c.id, "error", err)
				return
			}
		}
	}
}

func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}
