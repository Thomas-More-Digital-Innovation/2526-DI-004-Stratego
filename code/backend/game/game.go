package game

import (
	"digital-innovation/stratego/engine"
)

type Game struct {
	Players       []*engine.Player
	Board         *engine.Board
	CurrentPlayer *engine.Player
	MoveHistory   []engine.Move
}

func NewGame(player1 *engine.Player, player2 *engine.Player) *Game {
	board := engine.NewBoard()
	return &Game{
		Players:       []*engine.Player{player1, player2},
		Board:         board,
		CurrentPlayer: player1,
		MoveHistory:   []engine.Move{},
	}
}

func (g *Game) NextTurn() {
	if g.CurrentPlayer == g.Players[0] {
		g.CurrentPlayer = g.Players[1]
	} else {
		g.CurrentPlayer = g.Players[0]
	}
}

func (g *Game) MakeMove(move *engine.Move, piece *engine.Piece) []*engine.Piece {
	target := g.Board.GetPieceAt(move.GetTo())
	if target != nil {
		result := piece.Attack(target)
		piece, target = result[0], result[1]
		switch {
		case !piece.IsAlive() && !target.IsAlive():
			g.Board.RemovePieceAt(move.GetFrom())
			g.Board.RemovePieceAt(move.GetTo())
		case !piece.IsAlive():
			g.Board.RemovePieceAt(move.GetFrom())
		default:
			g.Board.MovePiece(move, piece)
		}
	} else {
		g.Board.MovePiece(move, piece)
	}
	g.MoveHistory = append(g.MoveHistory, *move)
	g.NextTurn()
	return []*engine.Piece{piece, target}
}
