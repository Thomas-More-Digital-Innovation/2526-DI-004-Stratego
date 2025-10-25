package game

import (
	"digital-innovation/stratego/engine"
)

type Game struct {
	Players       []*engine.Player
	Board         *engine.Board
	CurrentPlayer *engine.Player
	MoveHistory   []engine.Move
	round         int
	winner        *engine.Player
}

func NewGame(player1 *engine.Player, player2 *engine.Player) *Game {
	board := engine.NewBoard()
	return &Game{
		Players:       []*engine.Player{player1, player2},
		Board:         board,
		CurrentPlayer: player1,
		MoveHistory:   []engine.Move{},
		round:         1,
	}
}

func (g *Game) NextTurn() {
	switch {
	case g.Players[0].HasWon():
		g.winner = g.Players[0]
	case g.Players[1].HasWon():
		g.winner = g.Players[1]
	case g.CurrentPlayer == g.Players[0]:
		g.CurrentPlayer = g.Players[1]
	default:
		g.CurrentPlayer = g.Players[0]
		g.round++
	}
}

func (g *Game) GetRound() int {
	return g.round
}

func (g *Game) GetWinner() *engine.Player {
	return g.winner
}

func (g *Game) SetWinner(player *engine.Player) {
	g.winner = player
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

			if !target.IsAlive() {
				err = g.Board.RemovePieceAt(move.GetTo())
				if err != nil {
					panic(err) // errors should not happen here if function is used correctly
				}
			}
		} else {
			g.Board.MovePiece(move, piece)
		}
	} else {
		g.Board.MovePiece(move, piece)
	}
	g.MoveHistory = append(g.MoveHistory, *move)
	g.NextTurn()
	return []*engine.Piece{piece, target}
}
