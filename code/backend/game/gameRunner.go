package game

import (
	"digital-innovation/stratego/engine"
	"fmt"
	"log"
	"time"
)

// GameRunner handles the execution of a game turn-by-turn
type GameRunner struct {
	game            *Game
	turnDelay       time.Duration // Optional delay between AI turns for visualization
	maxTurns        int           // Safety limit to prevent infinite games
	waitingForInput bool          // True when waiting for human player input
}

func NewGameRunner(game *Game, turnDelay time.Duration, maxTurns int) *GameRunner {
	if maxTurns <= 0 {
		maxTurns = 1000 // Default safety limit
	}
	return &GameRunner{
		game:      game,
		turnDelay: turnDelay,
		maxTurns:  maxTurns,
	}
}

// RunToCompletion runs the game until it's over (for AI vs AI)
func (gr *GameRunner) RunToCompletion() *engine.Player {
	turnCount := 0
	log.Printf("GameRunner: Starting RunToCompletion loop")

	for !gr.game.IsGameOver() && turnCount < gr.maxTurns {
		executed := gr.ExecuteTurn()

		if executed {
			// Turn was executed, increment counter
			turnCount++
			log.Printf("GameRunner: Turn %d executed, currentPlayer=%s", turnCount, gr.game.CurrentPlayer.GetName())
		} else {
			// ExecuteTurn returned false - check why
			if gr.game.IsGameOver() {
				log.Printf("GameRunner: Game ended during ExecuteTurn")
				break
			}
			// Still waiting for human input, continue polling
			time.Sleep(100 * time.Millisecond)
			continue
		}

		// Optional delay for visualization
		if gr.turnDelay > 0 {
			time.Sleep(gr.turnDelay)
		}
	}

	if turnCount >= gr.maxTurns {
		fmt.Println("Game ended: Maximum turns reached")
		// Set winner to player with higher score, or nil for draw
		if gr.game.Players[0].GetPieceScore() > gr.game.Players[1].GetPieceScore() {
			gr.game.SetWinner(gr.game.Players[0], WinCauseMaxTurns)
		} else if gr.game.Players[1].GetPieceScore() > gr.game.Players[0].GetPieceScore() {
			gr.game.SetWinner(gr.game.Players[1], WinCauseMaxTurns)
		}
		// Otherwise winner remains nil (draw)
		return gr.game.GetWinner()
	}

	return gr.game.GetWinner()
}

// ExecuteTurn executes a single turn. Returns false if waiting for human input.
func (gr *GameRunner) ExecuteTurn() bool {
	if gr.game.IsGameOver() {
		log.Printf("GameRunner.ExecuteTurn: Game is over")
		return false
	}

	controller := gr.game.GetCurrentController()
	log.Printf("GameRunner.ExecuteTurn: Current player=%s, controllerType=%d",
		gr.game.CurrentPlayer.GetName(), controller.GetControllerType())

	// Check if human controller and if it has a pending move
	if controller.GetControllerType() == engine.HumanController {
		humanController, ok := controller.(*engine.HumanPlayerController)
		if !ok || !humanController.HasPendingMove() {
			gr.waitingForInput = true
			log.Printf("GameRunner.ExecuteTurn: Waiting for human input")
			return false // Wait for human input
		}

		// Get the pending move
		move := humanController.GetPendingMove()
		if move == nil {
			return false
		}

		// Execute the move
		piece := gr.game.Board.GetPieceAt(move.GetFrom())
		if piece == nil {
			fmt.Println("Invalid move: no piece at from position")
			return false
		}

		gr.game.MakeMove(move, piece)
		gr.waitingForInput = false
		return true
	}

	// AI controller - make move immediately
	move := controller.MakeMove(gr.game.Board)

	// Validate move - check if piece exists at from position
	piece := gr.game.Board.GetPieceAt(move.GetFrom())
	if piece == nil || piece.GetOwner() != gr.game.CurrentPlayer {
		// No piece at from position or wrong owner = AI has no valid moves
		log.Printf("AI %s has no valid moves (no piece at %v or wrong owner)",
			gr.game.CurrentPlayer.GetName(), move.GetFrom())
		opponent := gr.getOpponent(gr.game.CurrentPlayer)
		opponent.SetWinner()
		gr.game.SetWinner(opponent, WinCauseNoMovablePieces)
		fmt.Printf("%s has no valid moves remaining - %s wins!\n",
			gr.game.CurrentPlayer.GetName(), opponent.GetName())
		return false
	}

	// Validate the move is legal
	if !gr.game.Board.IsValidMove(&move) {
		log.Printf("AI %s provided invalid move: %v", gr.game.CurrentPlayer.GetName(), move)
		opponent := gr.getOpponent(gr.game.CurrentPlayer)
		opponent.SetWinner()
		gr.game.SetWinner(opponent, WinCauseNoMovablePieces)
		return false
	}

	gr.game.MakeMove(&move, piece)
	return true
}

// getOpponent returns the opponent of the given player
func (gr *GameRunner) getOpponent(player *engine.Player) *engine.Player {
	if gr.game.Players[0] == player {
		return gr.game.Players[1]
	}
	return gr.game.Players[0]
}

// IsWaitingForInput returns true if the game is waiting for human input
func (gr *GameRunner) IsWaitingForInput() bool {
	return gr.waitingForInput
}

// GetGame returns the underlying game
func (gr *GameRunner) GetGame() *Game {
	return gr.game
}

// SubmitHumanMove allows external code to submit a human player's move
func (gr *GameRunner) SubmitHumanMove(move engine.Move) error {
	if !gr.waitingForInput {
		return fmt.Errorf("not waiting for input")
	}

	controller := gr.game.GetCurrentController()
	if controller.GetControllerType() != engine.HumanController {
		return fmt.Errorf("current player is not human")
	}

	humanController, ok := controller.(*engine.HumanPlayerController)
	if !ok {
		return fmt.Errorf("invalid controller type")
	}

	// Validate the move
	if !gr.game.Board.IsValidMove(&move) {
		return fmt.Errorf("invalid move")
	}

	humanController.SetPendingMove(move)

	// Execute the turn
	gr.ExecuteTurn()

	return nil
}
