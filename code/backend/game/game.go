package game

import (
	"digital-innovation/stratego/engine"
)

type WinCause string

const (
	WinCauseFlagCaptured    WinCause = "flag_captured"
	WinCauseNoMovablePieces WinCause = "no_movable_pieces"
	WinCauseMaxTurns        WinCause = "max_turns"
)

type CombatResult struct {
	Occurred         bool
	AttackerPiece    *engine.Piece
	DefenderPiece    *engine.Piece
	AttackerPosition engine.Position
	DefenderPosition engine.Position
}

type Game struct {
	Players           []*engine.Player
	PlayerControllers []engine.PlayerController // AI or Human controllers
	Board             *engine.Board
	CurrentPlayer     *engine.Player
	CurrentController engine.PlayerController
	MoveHistory       []engine.Move
	LastCombat        *CombatResult // Track last combat for broadcasting
	round             int
	winner            *engine.Player
	winCause          WinCause
	gameOver          bool
}

func NewGame(controller1, controller2 engine.PlayerController) *Game {
	board := engine.NewBoard()
	player1 := controller1.GetPlayer()
	player2 := controller2.GetPlayer()

	return &Game{
		Players:           []*engine.Player{player1, player2},
		PlayerControllers: []engine.PlayerController{controller1, controller2},
		Board:             board,
		CurrentPlayer:     player1,
		CurrentController: controller1,
		MoveHistory:       []engine.Move{},
		round:             1,
		gameOver:          false,
	}
}

func (g *Game) NextTurn() {
	switch {
	case g.Players[0].HasWon():
		g.winner = g.Players[0]
		g.winCause = WinCauseFlagCaptured
		g.gameOver = true
	case g.Players[1].HasWon():
		g.winner = g.Players[1]
		g.winCause = WinCauseFlagCaptured
		g.gameOver = true
	case g.CurrentPlayer == g.Players[0]:
		g.CurrentPlayer = g.Players[1]
		g.CurrentController = g.PlayerControllers[1]
	default:
		g.CurrentPlayer = g.Players[0]
		g.CurrentController = g.PlayerControllers[0]
		g.round++
		// Hide all revealed pieces at the start of a new round
		g.HideAllRevealedPieces()
	}
}

func (g *Game) IsGameOver() bool {
	return g.gameOver
}

func (g *Game) GetCurrentController() engine.PlayerController {
	return g.CurrentController
}

func (g *Game) GetRound() int {
	return g.round
}

func (g *Game) GetWinner() *engine.Player {
	return g.winner
}

func (g *Game) GetWinCause() WinCause {
	return g.winCause
}

func (g *Game) SetWinner(player *engine.Player, cause WinCause) {
	g.winner = player
	g.winCause = cause
	g.gameOver = true
	player.SetWinner()
}

// MakeMove makes a move on the game board and resolves any combat that may occur.
// If the move results in combat, the attacker and defender pieces are revealed.
// The function returns a slice of two pieces: the attacker and defender pieces in the combat.
// If no combat occurs, the slice will contain only the attacker piece.
// The game state is updated after the move, and all observers (AI) are notified of the move.
// The observers are given the opportunity to analyze the move and observe any combat that may have occurred.
func (g *Game) MakeMove(move *engine.Move, piece *engine.Piece) []*engine.Piece {
	target := g.Board.GetPieceAt(move.GetTo())
	if target != nil {
		piece.Reveal()
		target.Reveal()

		g.LastCombat = &CombatResult{
			Occurred:         true,
			AttackerPiece:    piece,
			DefenderPiece:    target,
			AttackerPosition: move.GetFrom(),
			DefenderPosition: move.GetTo(),
		}

		result := piece.Attack(target)
		piece, target = result[0], result[1]
		if !piece.IsAlive() {
			err := g.Board.RemovePieceAt(move.GetFrom())
			if err != nil {
				panic(err) // errors should not happen here if function is used correctly
			}

			if target != nil && !target.IsAlive() {
				err = g.Board.RemovePieceAt(move.GetTo())
				if err != nil {
					panic(err) // errors should not happen here if function is used correctly
				}
			}
		} else {
			g.Board.MovePiece(move, piece)
			piece.GetOwner().UpdatePiecePosition(piece, move.GetTo())
		}
	} else {
		// No combat - clear any previous combat result
		g.LastCombat = nil
		g.Board.MovePiece(move, piece)
		piece.GetOwner().UpdatePiecePosition(piece, move.GetTo())
	}
	g.MoveHistory = append(g.MoveHistory, *move)

	// Notify all observers (AI)
	round := g.GetRound()
	for _, ctrl := range g.PlayerControllers {
		if ctrl.GetPlayer() == move.GetPlayer() {
			continue
		}

		if analyzer, ok := ctrl.(interface {
			AnalyzeMove(engine.Move, *engine.Player, int)
		}); ok {
			analyzer.AnalyzeMove(*move, move.GetPlayer(), round)
		}

		if g.LastCombat != nil && g.LastCombat.Occurred {
			if observer, ok := ctrl.(interface {
				ObserveCombat(engine.Position, engine.Position, *engine.Piece, *engine.Piece, int)
			}); ok {
				observer.ObserveCombat(
					g.LastCombat.AttackerPosition,
					g.LastCombat.DefenderPosition,
					g.LastCombat.AttackerPiece,
					g.LastCombat.DefenderPiece,
					round,
				)
			}
		}
	}

	g.NextTurn()
	return []*engine.Piece{piece, target}
}

// GetLastCombat returns the last combat result if any
func (g *Game) GetLastCombat() *CombatResult {
	return g.LastCombat
}

// ClearLastCombat clears the last combat result (called after broadcast)
func (g *Game) ClearLastCombat() {
	g.LastCombat = nil
}

// HideCombatPieces hides the pieces involved in the last combat
func (g *Game) HideCombatPieces() {
	if g.LastCombat != nil && g.LastCombat.Occurred {
		if g.LastCombat.AttackerPiece != nil && g.LastCombat.AttackerPiece.IsAlive() {
			g.LastCombat.AttackerPiece.Hide()
		}
		if g.LastCombat.DefenderPiece != nil && g.LastCombat.DefenderPiece.IsAlive() {
			g.LastCombat.DefenderPiece.Hide()
		}
	}
}

// HideAllRevealedPieces hides all revealed pieces on the board
// Called at the start of each new round to reset piece visibility
func (g *Game) HideAllRevealedPieces() {
	field := g.Board.GetField()
	for y := range 10 {
		for x := range 10 {
			piece := field[y][x]
			if piece != nil && piece.IsAlive() && piece.IsRevealed() {
				piece.Hide()
			}
		}
	}
}

// InitializePieces scans board and tracks all pieces for both players (call once at game start)
func (g *Game) InitializePieces() {
	field := g.Board.GetField()
	for y := range 10 {
		for x := range 10 {
			piece := field[y][x]
			if piece != nil {
				pos := engine.NewPosition(x, y)
				piece.GetOwner().AddPiece(piece, pos)
			}
		}
	}
}
