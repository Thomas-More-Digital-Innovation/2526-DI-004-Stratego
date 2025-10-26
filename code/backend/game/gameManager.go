package game

import (
	"digital-innovation/stratego/engine"
)

type GameManager struct {
	ActiveGame int
	Games      []*Game
}

func NewGameManager() *GameManager {
	return &GameManager{
		Games: []*Game{},
	}
}

func (gm *GameManager) RunMatch(controller1, controller2 engine.PlayerController) *Game {
	game := NewGame(controller1, controller2)
	gm.Games = append(gm.Games, game)
	gm.ActiveGame = len(gm.Games) - 1
	return game
}

func (gm *GameManager) RunTournament(rounds int, controller1, controller2 engine.PlayerController) []*Game {
	if rounds == 0 {
		return gm.Games
	}

	gm.ActiveGame = max(len(gm.Games)-1, 0)

	games := []*Game{}
	for range rounds {
		games = append(games, gm.RunMatch(controller1, controller2))
	}

	return games
}

func (gm *GameManager) GetActiveGame() *Game {
	if len(gm.Games) == 0 {
		return nil
	}
	return gm.Games[gm.ActiveGame]
}

func (gm *GameManager) NextGame(index int) {
	if index < 0 || index >= len(gm.Games) {
		return
	}
	gm.ActiveGame = index
}

func (gm *GameManager) GetGameCount() int {
	return len(gm.Games)
}
