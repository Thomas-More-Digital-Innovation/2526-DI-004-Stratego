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

type Game struct {
	Players           []*engine.Player
	PlayerControllers []engine.PlayerController // AI or Human controllers
	Board             *engine.Board
	CurrentPlayer     *engine.Player
	CurrentController engine.PlayerController
	MoveHistory       []engine.Move
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
}

func (g *Game) MakeMove(move *engine.Move, piece *engine.Piece) []*engine.Piece {
	target := g.Board.GetPieceAt(move.GetTo())
	if target != nil {
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
		g.Board.MovePiece(move, piece)
		piece.GetOwner().UpdatePiecePosition(piece, move.GetTo())
	}
	g.MoveHistory = append(g.MoveHistory, *move)
	g.NextTurn()
	return []*engine.Piece{piece, target}
}

// InitializePieces scans board and tracks all pieces for both players (call once at game start)
func (g *Game) InitializePieces() {
	field := g.Board.GetField()
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			if piece != nil {
				pos := engine.NewPosition(x, y)
				piece.GetOwner().AddPiece(piece, pos)
			}
		}
	}
}
