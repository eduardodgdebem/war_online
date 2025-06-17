package game

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	Conn    *websocket.Conn `json:"-"`
	Name    string          `json:"name"`
	IsAlive bool            `json:"isAlive"`
	IsReady bool            `json:"isReady"`
	Id      uint32          `json:"id"`
}

func NewPlayer(conn *websocket.Conn, name string) *Player {
	return &Player{
		Name:    name,
		Conn:    conn,
		IsAlive: true,
		Id:      uuid.New().ID(),
		IsReady: false,
	}
}

type PlayerMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
