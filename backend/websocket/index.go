package websocket

import (
	"crypto/rand"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"war-backend/game"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	game_instance = game.GetGameInstance()
)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	game_instance.AddPlayer(conn, rand.Text())

	// BroadcastPlayerCount()

	// if game_instance.CanStartGame() {
	// 	game_instance.StartGame()
	// }

	go game_instance.HandleConnection(conn)
}

func BroadcastPlayerCount() {
	game_instance.BroadcastMessage("count", game_instance.PlayerCount(), nil)
}
