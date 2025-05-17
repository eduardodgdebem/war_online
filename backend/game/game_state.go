package game

import "sync"

const (
	ReinforcePhase = "reinforcePhase"
	AttackPhase    = "attackPhase"
)

type gameRound struct {
	currentPlayerId uint32
	phase           string
}

func newGameRound() *gameRound {
	return &gameRound{}
}

type gameState struct {
	players      map[*Player]bool
	playersMutex sync.RWMutex
	gameStarted  bool
	gameMap      Map
	round        gameRound
}

func newGameState() *gameState {
	return &gameState{
		players:     make(map[*Player]bool),
		gameMap:     *NewMap(),
		gameStarted: false,
		round:       *newGameRound(),
	}
}

func (gs *gameState) formatGameState() map[string]any {
	gs.playersMutex.RLock()
	defer gs.playersMutex.RUnlock()

	players := make([]Player, 0, len(gs.players))
	for p := range gs.players {
		players = append(players, *p)
	}

	return map[string]any{
		"map":     gs.gameMap,
		"players": players,
	}
}

func (gs *gameState) CanStartGame() bool {
	gs.playersMutex.RLock()
	defer gs.playersMutex.RUnlock()

	if len(gs.players) < 2 { // Minimum 2 players to start
		return false
	}

	for player := range gs.players {
		if !player.IsReady {
			return false
		}
	}
	return true
}
