package game

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	GetGameState = "getGameState"
	Welcome      = "welcome"
	SetIsReady   = "setIsReady"
	Attack       = "attack"
	Reforce      = "reforce"
)

type messageHandler struct {
	game *Game
}

func newMessageHandler(game *Game) *messageHandler {
	return &messageHandler{game: game}
}

func (mh *messageHandler) SendMessage(msgType string, payload any, conn *websocket.Conn) {
	if mh == nil || conn == nil {
		log.Println("Message handler or connection is nil")
		return
	}

	msg := map[string]any{
		"type":    msgType,
		"payload": payload,
	}

	conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Send error to %s: %v", conn.RemoteAddr(), err)
		go mh.game.RemovePlayer(conn)
	}
	conn.SetWriteDeadline(time.Time{}) // Reset deadline
}

func (mh *messageHandler) BroadcastMessage(msgType string, payload any, excludeSender *websocket.Conn) {
	mh.game.playersMutex.RLock()
	defer mh.game.playersMutex.RUnlock()

	msg := map[string]any{
		"type":    msgType,
		"payload": payload,
	}

	for player := range mh.game.players {
		if player.Conn != excludeSender && player.Conn != nil {
			if err := player.Conn.WriteJSON(msg); err != nil {
				log.Printf("Broadcast error to player '%s': %v", player.Name, err)
				go mh.game.RemovePlayer(player.Conn)
			}
		}
	}
}

func (mh *messageHandler) handleMessage(conn *websocket.Conn, msg PlayerMessage) {
	if mh == nil || mh.game == nil {
		log.Println("Message handler or game is nil")
		return
	}

	switch msg.Type {
	case GetGameState:
		mh.handleGetGameState(conn)
	case SetIsReady:
		mh.handleSetIsReady(conn, msg)
	case Reforce:
		mh.handleReforce(conn, msg)
	case Attack:
		mh.handleAttack(conn)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

type reinforceMessage struct {
	cellX    int
	cellY    int
	playerId uint32
}

func (mh *messageHandler) handleReforce(conn *websocket.Conn, msg PlayerMessage) {
	rMsg, ok := msg.Payload.(reinforceMessage)
	if !ok {
		return
	}
	mh.game.ReinforceCell(rMsg.playerId, rMsg.cellX, rMsg.cellY)
	mh.handleGetGameState(conn)
}

func (mh *messageHandler) handleGetGameState(conn *websocket.Conn) {
	mh.game.SendMessage(GetGameState, mh.game.formatGameState(), conn)
}

func (mh *messageHandler) handleSetIsReady(conn *websocket.Conn, msg PlayerMessage) {
	player := mh.game.getPlayerByConn(conn)
	if player != nil {
		player.IsReady = true
		mh.game.BroadcastPlayerUpdate()
	}
}

func (mh *messageHandler) handleAttack(conn *websocket.Conn) {
	mh.game.BroadcastMessage(Attack, "attacking", conn)
}
