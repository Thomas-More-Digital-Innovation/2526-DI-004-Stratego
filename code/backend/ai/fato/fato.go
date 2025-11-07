package fato

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
	"math"
	"math/rand/v2"
)

type FatoAI struct {
	player         *engine.Player
	memorizedField [10][10]*engine.Piece
}

func NewFatoAI(player *engine.Player) *FatoAI {
	return &FatoAI{
		memorizedField: [10][10]*engine.Piece{},
		player:         player,
	}
}

func (ai *FatoAI) GetPlayer() *engine.Player {
	return ai.player
}

func (ai *FatoAI) GetControllerType() engine.ControllerType {
	return engine.AIController
}

func (ai *FatoAI) IsPieceMemorized(pos engine.Position) bool {
	return ai.memorizedField[pos.Y][pos.X] != nil
}

func (ai *FatoAI) PickRandomPiece() *engine.Piece {
	pieces := ai.player.GetAlivePieces()
	if len(pieces) == 0 {
		return nil
	}
	random := rand.IntN(len(pieces))
	return pieces[random]
}

func (ai *FatoAI) MakeMove(board *engine.Board) engine.Move {
	// Not so random huh? :-)

	// 1. Try to attack a known enemy piece
	if move, found := ai.findAttackMove(board); found {
		return move
	}

	// 2. Try to explore toward enemy territory
	if move, found := ai.findExplorationMove(board); found {
		return move
	}

	// 3. Fallback: random valid move
	return ai.findRandomMove(board)
}

// findAttackMove looks for moves that attack known/visible enemy pieces
func (ai *FatoAI) findAttackMove(board *engine.Board) (engine.Move, bool) {
	pieces := ai.player.GetAlivePieces()
	for _, piece := range pieces {
		if !piece.CanMove() {
			continue
		}
		pos, exists := ai.player.GetPiecePosition(piece)
		if !exists {
			continue
		}

		moves, err := board.ListMoves(pos)
		if err != nil {
			continue
		}

		for _, move := range moves {
			target := board.GetPieceAt(move.GetTo())
			if target != nil && target.GetOwner() != ai.player {
				// Attack if we know the piece or it's worth trying
				memorized := ai.RecallPiece(move.GetTo())
				if memorized != nil || piece.GetRank() >= target.GetRank() {
					return move, true
				}
			}
		}
	}
	return engine.Move{}, false
}

// findExplorationMove moves toward enemy side
func (ai *FatoAI) findExplorationMove(board *engine.Board) (engine.Move, bool) {
	enemyY := 0
	if ai.player.GetID() == 1 {
		enemyY = 9
	}

	pieces := ai.player.GetAlivePieces()
	shuffled := make([]*engine.Piece, len(pieces))
	copy(shuffled, pieces)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	for _, piece := range shuffled {
		if !piece.CanMove() {
			continue
		}
		pos, exists := ai.player.GetPiecePosition(piece)
		if !exists {
			continue
		}

		moves, err := board.ListMoves(pos)
		if err != nil {
			continue
		}

		var bestMove *engine.Move
		bestDist := 100
		for _, move := range moves {
			if board.GetPieceAt(move.GetTo()) == nil {
				dist := int(math.Abs(float64(move.GetTo().Y - enemyY)))
				if dist < bestDist {
					bestDist = dist
					bestMove = &move
				}
			}
		}
		if bestMove != nil {
			return *bestMove, true
		}
	}
	return engine.Move{}, false
}

// findRandomMove picks any valid move as last resort
func (ai *FatoAI) findRandomMove(board *engine.Board) engine.Move {
	pieces := ai.player.GetAlivePieces()
	shuffled := make([]*engine.Piece, len(pieces))
	copy(shuffled, pieces)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	for _, piece := range shuffled {
		if !piece.CanMove() {
			continue
		}
		pos, exists := ai.player.GetPiecePosition(piece)
		if !exists {
			continue
		}

		moves, err := board.ListMoves(pos)
		if err != nil || len(moves) == 0 {
			continue
		}
		return moves[rand.IntN(len(moves))]
	}

	// No valid moves available - player has lost (only immobile pieces left)
	// Return empty move to signal defeat
	return engine.Move{}
}

func (ai *FatoAI) AnalyzeMove(opponentMove engine.Move, opponent *engine.Player) {
	if math.Abs(float64(opponentMove.GetFrom().X-opponentMove.GetTo().X)) > 1 || math.Abs(float64(opponentMove.GetFrom().Y-opponentMove.GetTo().Y)) > 1 {
		piece := engine.NewPiece(models.Scout, opponent)
		ai.MemorizePiece(opponentMove.GetTo(), piece)
	}
}

func (ai *FatoAI) MemorizePiece(pos engine.Position, piece *engine.Piece) {
	ai.memorizedField[pos.Y][pos.X] = piece
}

func (ai *FatoAI) RecallPiece(pos engine.Position) *engine.Piece {
	return ai.memorizedField[pos.Y][pos.X]
}

// small chance to forget a piece
func (ai *FatoAI) ForgetPiece(pos engine.Position) {
	ai.memorizedField[pos.Y][pos.X] = nil
}

// Reset clears the AI's memory for a new game
func (ai *FatoAI) Reset() {
	ai.memorizedField = [10][10]*engine.Piece{}
}
