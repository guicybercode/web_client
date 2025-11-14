package main

import (
	"encoding/json"
	"time"
)

type Message struct {
	ClientID  string `json:"client_id"`
	Timestamp int64  `json:"timestamp"`
	Content   string `json:"content"`
}

func NewMessage(clientID, content string) *Message {
	return &Message{
		ClientID:  clientID,
		Timestamp: time.Now().UnixMilli(),
		Content:   content,
	}
}

func (m *Message) ToJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func MessageFromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
