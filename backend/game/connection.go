package game

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxPlayers     = 4
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

type connectionManager struct {
	game         *Game
	shutdownChan chan struct{}
}

func newConnectionManager(game *Game) *connectionManager {
	return &connectionManager{
		game:         game,
		shutdownChan: make(chan struct{}),
	}
}

func (cm *connectionManager) HandleConnection(conn *websocket.Conn) {
	if cm == nil || cm.game == nil {
		log.Println("Connection manager or game is nil")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in HandleConnection: %v", r)
		}
		conn.Close()
		cm.game.RemovePlayer(conn)
	}()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-cm.shutdownChan:
			return
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		default:
			var msg PlayerMessage
			if err := conn.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
					websocket.CloseNoStatusReceived) {
					log.Printf("Read error: %v", err)
				}
				return
			}
			cm.game.handleMessage(conn, msg)
		}
	}
}

func (cm *connectionManager) Shutdown() {
	close(cm.shutdownChan)
	cm.game.playersMutex.Lock()
	defer cm.game.playersMutex.Unlock()

	for player := range cm.game.players {
		if player.Conn != nil {
			// Send proper close message
			player.Conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				time.Now().Add(5*time.Second),
			)
			player.Conn.Close()
		}
	}
}
