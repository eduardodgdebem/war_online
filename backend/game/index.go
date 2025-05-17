package game

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	maxPlayers = 4
)

var (
	instance *Game
	once     sync.Once
)

type Game struct {
	players      map[*Player]bool
	playersMutex sync.RWMutex
	gameStarted  bool
	shutdownChan chan struct{}
	gameMap      Map
}

func GetGameInstance() *Game {
	once.Do(func() {
		instance = &Game{
			players:      make(map[*Player]bool),
			shutdownChan: make(chan struct{}),
			gameMap:      *NewMap(),
		}
	})
	return instance
}

func (g *Game) getPlayerByConn(conn *websocket.Conn) *Player {
	g.playersMutex.RLock()
	defer g.playersMutex.RUnlock()

	for player := range g.players {
		if player.Conn == conn {
			return player
		}
	}
	return nil
}

func (g *Game) Shutdown() {
	close(g.shutdownChan)
	g.playersMutex.Lock()
	defer g.playersMutex.Unlock()

	for player := range g.players {
		player.Conn.Close()
		delete(g.players, player)
	}
}

func (g *Game) AddPlayer(conn *websocket.Conn, name string) error {
	g.playersMutex.Lock()
	defer g.playersMutex.Unlock()

	player := NewPlayer(conn, name)

	g.players[player] = true
	log.Printf("New player '%s' connected. Total players: %d/%d", name, len(g.players), maxPlayers)

	welcomeMsg := map[string]any{
		"type":         "welcome",
		"payload":      "Welcome to the War Game!",
		"playersCount": len(g.players),
		"max":          maxPlayers,
		"name":         name,
		"player": map[string]any{
			"name": player.Name,
			"id":   player.Id,
		},
	}

	if err := conn.WriteJSON(welcomeMsg); err != nil {
		delete(g.players, player)
		return err
	}

	return nil
}

func (g *Game) RemovePlayer(conn *websocket.Conn) {
	g.playersMutex.Lock()
	defer g.playersMutex.Unlock()

	player := g.getPlayerByConn(conn)
	if player != nil {
		conn.Close()
		log.Printf("Player '%s' disconnected. Remaining players: %d/%d", player.Name, len(g.players)-1, maxPlayers)
		delete(g.players, player)
	}
}

func (g *Game) BroadcastMessage(msgType string, payload any, excludeSender *websocket.Conn) {
	g.playersMutex.RLock()
	defer g.playersMutex.RUnlock()

	msg := map[string]any{
		"type":    msgType,
		"payload": payload,
	}

	for player := range g.players {
		if player.Conn != excludeSender {
			if err := player.Conn.WriteJSON(msg); err != nil {
				log.Printf("Broadcast error to player '%s': %v", player.Name, err)
				go g.RemovePlayer(player.Conn)
			}
		}
	}
}

func (g *Game) HandleConnection(conn *websocket.Conn) {
	defer g.RemovePlayer(conn)

	for {
		select {
		case <-g.shutdownChan:
			return
		default:
			var msg PlayerMessage
			if err := conn.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Read error: %v", err)
				}
				return
			}

			g.handleMessage(conn, msg)
		}
	}
}

func (g *Game) handleMessage(conn *websocket.Conn, msg PlayerMessage) {
	switch msg.Type {
	case "join":
		g.handleJoin(conn, msg.Payload)
	case "move":
		g.BroadcastMessage("move", msg.Payload, conn)
	case "attack":
		g.BroadcastMessage("attack", msg.Payload, conn)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func (g *Game) handleJoin(conn *websocket.Conn, name string) {
	response := map[string]any{
		"type":    "joinResponse",
		"payload": "Joined as " + name,
		"name":    name,
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Join response error: %v", err)
	}
}

func (g *Game) CanStartGame() bool {
	g.playersMutex.RLock()
	defer g.playersMutex.RUnlock()
	return len(g.players) >= 2 && !g.gameStarted
}

func (g *Game) StartGame() {
	g.gameStarted = true
	log.Println("Starting the war game!")

	playerNames := make([]string, 0, len(g.players))
	g.playersMutex.RLock()
	for player := range g.players {
		playerNames = append(playerNames, player.Name)
	}
	g.playersMutex.RUnlock()

	startMsg := map[string]any{
		"type":    "gameStart",
		"payload": "The war has begun!",
		"players": playerNames,
	}

	g.playersMutex.RLock()
	defer g.playersMutex.RUnlock()
	for player := range g.players {
		if err := player.Conn.WriteJSON(startMsg); err != nil {
			log.Printf("Game start error for player '%s': %v", player.Name, err)
			go g.RemovePlayer(player.Conn)
		}
	}
}

func (g *Game) PlayerCount() int {
	g.playersMutex.RLock()
	defer g.playersMutex.RUnlock()
	return len(g.players)
}

func (g *Game) IsGameStarted() bool {
	g.playersMutex.RLock()
	defer g.playersMutex.RUnlock()
	return g.gameStarted
}

func (g *Game) GetPlayerNames() []string {
	g.playersMutex.RLock()
	defer g.playersMutex.RUnlock()
	names := make([]string, 0, len(g.players))
	for player := range g.players {
		names = append(names, player.Name)
	}
	return names
}
