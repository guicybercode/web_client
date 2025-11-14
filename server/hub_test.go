package main

import (
	"testing"
	"time"
)

func TestHub_NewHub(t *testing.T) {
	hub := NewHub()
	if hub.clients == nil {
		t.Error("clients map not initialized")
	}
	if hub.register == nil {
		t.Error("register channel not initialized")
	}
	if hub.unregister == nil {
		t.Error("unregister channel not initialized")
	}
	if hub.broadcast == nil {
		t.Error("broadcast channel not initialized")
	}
}

func TestHub_Run_RegisterUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{hub: hub, send: make(chan []byte, 256)}

	hub.register <- client
	time.Sleep(10 * time.Millisecond)

	if len(hub.clients) != 1 {
		t.Errorf("expected 1 client, got %d", len(hub.clients))
	}

	hub.unregister <- client
	time.Sleep(10 * time.Millisecond)

	if len(hub.clients) != 0 {
		t.Errorf("expected 0 clients after unregister, got %d", len(hub.clients))
	}
}

func TestHub_Run_Broadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client1 := &Client{hub: hub, send: make(chan []byte, 256)}
	client2 := &Client{hub: hub, send: make(chan []byte, 256)}

	hub.register <- client1
	hub.register <- client2
	time.Sleep(10 * time.Millisecond)

	message := []byte("test message")
	hub.broadcast <- message
	time.Sleep(10 * time.Millisecond)

	select {
	case received := <-client1.send:
		if string(received) != string(message) {
			t.Errorf("client1 received wrong message")
		}
	default:
		t.Error("client1 did not receive message")
	}

	select {
	case received := <-client2.send:
		if string(received) != string(message) {
			t.Errorf("client2 received wrong message")
		}
	default:
		t.Error("client2 did not receive message")
	}
}
