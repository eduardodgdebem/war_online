package game

import (
	"errors"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Game struct {
	*connectionManager
	*messageHandler
	*gameState
}

var (
	instance *Game
	once     sync.Once
)

func GetGameInstance() *Game {
	once.Do(func() {
		instance = &Game{}

		instance.gameState = newGameState()
		instance.connectionManager = newConnectionManager(instance)
		instance.messageHandler = newMessageHandler(instance)
	})
	return instance
}

func (g *Game) AddPlayer(conn *websocket.Conn, name string) error {
	g.playersMutex.Lock()
	defer g.playersMutex.Unlock()

	remoteAddr := conn.RemoteAddr().String()
	for p := range g.players {
		if p.Conn != nil && p.Conn.RemoteAddr().String() == remoteAddr {
			p.Conn.Close()
			delete(g.players, p)
			break
		}
	}

	if len(g.players) >= maxPlayers {
		return ErrGameFull
	}

	player := NewPlayer(conn, name)
	g.players[player] = true

	log.Printf("New player '%s' connected. Total players: %d/%d",
		name, len(g.players), maxPlayers)

	g.SendWelcomeMessage(player)
	return nil
}

func (g *Game) RemovePlayer(conn *websocket.Conn) {
	if g == nil {
		log.Println("Game instance is nil")
		return
	}

	g.playersMutex.Lock()
	defer g.playersMutex.Unlock()

	player := g.getPlayerByConn(conn)
	if player == nil {
		return
	}

	if player.Conn != nil {
		player.Conn.Close()
	}

	delete(g.players, player)
	log.Printf("Player disconnected. Remaining players: %d/%d", len(g.players), maxPlayers)
	g.BroadcastPlayerUpdate()
}

func (g *Game) getPlayerByConn(conn *websocket.Conn) *Player {
	for player := range g.players {
		if player.Conn == conn {
			return player
		}
	}
	return nil
}

func (g *Game) SendWelcomeMessage(player *Player) {
	msg := map[string]any{
		"type": Welcome,
		"payload": map[string]any{
			"playersCount": len(g.players),
			"maxPlayers":   maxPlayers,
			"player":       player,
		},
	}
	player.Conn.WriteJSON(msg)
}

func (g *Game) BroadcastPlayerUpdate() {
	g.BroadcastMessage("playerUpdate", g.formatGameState(), nil)
}

func (g *Game) nextPhase() {
	if g.gameState.round.phase == ReinforcePhase {
		g.gameState.round.phase = AttackPhase
	} else {
		g.gameState.round.phase = ReinforcePhase
	}
}

func (g *Game) ReinforceCell(playerId uint32, cellX int, cellY int) {
	cell := g.gameMap.Nodes[cellY][cellX]

	if playerId != g.round.currentPlayerId || playerId != cell.OwnerId {
		return
	}

	cell.Strengh++
}

var (
	ErrGameFull = errors.New("game is full")
)
