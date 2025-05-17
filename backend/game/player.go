package game

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type PlayerMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type Player struct {
	Conn    *websocket.Conn
	Name    string
	IsAlive bool
	Id      uint32
}

func NewPlayer(conn *websocket.Conn, name string) *Player {
	id := uuid.New().ID()

	return &Player{
		Name:    name,
		Conn:    conn,
		IsAlive: true,
		Id:      id,
	}
}
