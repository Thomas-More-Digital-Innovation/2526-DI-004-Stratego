package game

import (
	"digital-innovation/stratego/engine"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// GameRunner handles the execution of a game turn-by-turn
type GameRunner struct {
	game            *Game
	turnDelay       time.Duration // Optional delay between AI turns for visualization, can be 0 to remove the delay
	maxTurns        int           // Safety limit to prevent infinite games
	waitingForInput bool          // True when waiting for human player input
	onMoveExecuted  func()        // Callback when a move is executed
	stopChan        chan bool     // Channel to signal stopping the game
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

// SetMoveCallback sets the callback to be called when a move is executed
func (gr *GameRunner) SetMoveCallback(callback func()) {
	gr.onMoveExecuted = callback
}

// RunToCompletion runs the game until it's over (for AI vs AI)
func (gr *GameRunner) RunToCompletion(logging bool) *engine.Player {
	turnCount := 0
	if logging {
		log.Printf("GameRunner: Starting RunToCompletion loop")
	}

	for !gr.game.IsGameOver() && turnCount < gr.maxTurns {
		// Check for stop signal
		select {
		case <-gr.stopChan:
			if logging {
				log.Printf("GameRunner: Stop signal received, ending game")
			}
			return nil
		default:
			// No stop signal, continue
		}
		
		executed := gr.ExecuteTurn(logging)

		if executed {
			// Turn was executed, increment counter
			turnCount++
			if logging {
				log.Printf("GameRunner: Turn %d executed, currentPlayer=%s", turnCount, gr.game.CurrentPlayer.GetName())
			}
		} else {
			// ExecuteTurn returned false - check why
			if gr.game.IsGameOver() {
				if logging {
					log.Printf("GameRunner: Game ended during ExecuteTurn")
				}
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
		return gr.calculateWinnerOnMaxTurnsExceeded()
	}

	return gr.game.GetWinner()
}

func (gr *GameRunner) calculateWinnerOnMaxTurnsExceeded() *engine.Player {
	if float64(gr.game.Players[0].GetPieceScore())/float64(gr.game.Players[1].GetPieceScore()) > 1.15 {
		gr.game.SetWinner(gr.game.Players[0], WinCauseMaxTurns)
	} else if float64(gr.game.Players[1].GetPieceScore())/float64(gr.game.Players[0].GetPieceScore()) > 1.15 {
		gr.game.SetWinner(gr.game.Players[1], WinCauseMaxTurns)
	}
	return gr.game.GetWinner()
}

// ExecuteTurn executes a single turn. Returns false if waiting for human input.
func (gr *GameRunner) ExecuteTurn(logging bool) bool {
	if gr.game.IsGameOver() {
		log.Printf("GameRunner.ExecuteTurn: Game is over")
		return false
	}

	controller := gr.game.GetCurrentController()
	if logging {
		log.Printf("GameRunner.ExecuteTurn: Current player=%s, controllerType=%d",
			gr.game.CurrentPlayer.GetName(), controller.GetControllerType())
	}

	// Check if human controller and if it has a pending move
	if controller.GetControllerType() == engine.HumanController {
		humanController, ok := controller.(*engine.HumanPlayerController)
		if !ok || !humanController.HasPendingMove() {
			if !gr.waitingForInput {
				// Only log once when we start waiting
				if logging {
					log.Printf("GameRunner.ExecuteTurn: Waiting for human input")
				}
				gr.waitingForInput = true
			}
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

		// Notify that a move was executed
		if gr.onMoveExecuted != nil {
			gr.onMoveExecuted()
		}
		return true
	}

	// AI controller - make move
	// Add delay between 500ms and 1000ms before AI moves (for pacing)
	if gr.turnDelay > 0 {
		aiDelay := time.Duration(500+rand.Intn(500)) * time.Millisecond
		time.Sleep(aiDelay)
	}

	move := controller.MakeMove(gr.game.Board)

	// Validate move - check if piece exists at from position
	piece := gr.game.Board.GetPieceAt(move.GetFrom())
	if piece == nil || piece.GetOwner() != gr.game.CurrentPlayer {
		// No piece at from position or wrong owner = AI has no valid moves
		if logging {
			log.Printf("AI %s has no valid moves (no piece at %v or wrong owner)",
				gr.game.CurrentPlayer.GetName(), move.GetFrom())
		}
		opponent := gr.getOpponent(gr.game.CurrentPlayer)
		opponent.SetWinner()
		gr.game.SetWinner(opponent, WinCauseNoMovablePieces)
		fmt.Printf("%s has no valid moves remaining - %s wins!\n",
			gr.game.CurrentPlayer.GetName(), opponent.GetName())
		return false
	}

	// Validate the move is legal
	if !gr.game.Board.IsValidMove(&move) {
		if logging {
			log.Printf("AI %s provided invalid move: %v", gr.game.CurrentPlayer.GetName(), move)
		}
		opponent := gr.getOpponent(gr.game.CurrentPlayer)
		opponent.SetWinner()
		gr.game.SetWinner(opponent, WinCauseNoMovablePieces)
		return false
	}

	gr.game.MakeMove(&move, piece)

	// Notify that a move was executed
	if gr.onMoveExecuted != nil {
		gr.onMoveExecuted()
	}
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

func (gr *GameRunner) DebugSetWaitingForInput(value bool) {
	gr.waitingForInput = value
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
	gr.ExecuteTurn(true) // TODO assuming logging is true for human moves

	return nil
}
