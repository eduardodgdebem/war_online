package websocket

import (
	"crypto/rand"
	"log"
	"net/http"
	"time"

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

	name := r.URL.Query().Get("name")
	if name == "" {
		// name = "Anonymous"
		name = rand.Text()
	}

	game := game.GetGameInstance()
	if err := game.AddPlayer(conn, name); err != nil {
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, err.Error()),
			time.Now().Add(5*time.Second),
		)
		conn.Close()
		return
	}

	go game_instance.HandleConnection(conn)
}
