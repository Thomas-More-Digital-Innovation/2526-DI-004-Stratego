// Package ai provides the base AI interfaces and shared logic.
package ai

import (
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/game/models"
)

// SimulateMove runs a copy-on-write move simulation on the board.
func SimulateMove(board *game.Board, move game.Move) *game.Board {
	nextBoard := board.FastClone()
	attacker := nextBoard.GetPieceAt(move.GetFrom())
	if attacker == nil {
		return nextBoard
	}

	target := nextBoard.GetPieceAt(move.GetTo())
	if target != nil {
		attackerRank := attacker.GetRank()
		defenderRank := target.GetRank()

		if defenderRank == models.Flag.GetRank() {
			nextBoard.SetPieceAt(move.GetFrom(), nil)
			nextBoard.SetPieceAt(move.GetTo(), attacker)
			return nextBoard
		}

		if attackerRank == models.Spy.GetRank() && defenderRank == models.Marshal.GetRank() {
			nextBoard.SetPieceAt(move.GetFrom(), nil)
			nextBoard.SetPieceAt(move.GetTo(), attacker)
			return nextBoard
		}

		if defenderRank == models.Bomb.GetRank() {
			if attacker.GetType().GetName() == "Miner" {
				nextBoard.SetPieceAt(move.GetFrom(), nil)
				nextBoard.SetPieceAt(move.GetTo(), attacker)
			} else {
				nextBoard.SetPieceAt(move.GetFrom(), nil)
			}
			return nextBoard
		}

		switch {
		case attackerRank > defenderRank:
			nextBoard.SetPieceAt(move.GetFrom(), nil)
			nextBoard.SetPieceAt(move.GetTo(), attacker)
		case attackerRank < defenderRank:
			nextBoard.SetPieceAt(move.GetFrom(), nil)
		default:
			nextBoard.SetPieceAt(move.GetFrom(), nil)
			nextBoard.SetPieceAt(move.GetTo(), nil)
		}
	} else {
		nextBoard.SetPieceAt(move.GetFrom(), nil)
		nextBoard.SetPieceAt(move.GetTo(), attacker)
	}
	return nextBoard
}

// EvaluateBoard calculates a static evaluation score for the current player's board.
func EvaluateBoard(board *game.Board, player *game.Player, memory *Memory, weights map[string]float64, aggression float64) float64 {
	score := 0.0
	playerID := player.GetID()

	for y := range 10 {
		for x := range 10 {
			pos := game.NewPosition(x, y)
			piece := board.GetPieceAt(pos)
			if piece == nil {
				continue
			}

			name := piece.GetType().GetName()
			val, ok := weights[name]
			if !ok {
				val = float64(piece.GetStrategicValue())
			}

			if piece.GetOwner().GetID() == playerID {
				score += val

				if piece.CanMove() {
					var dist int
					if playerID == 0 {
						dist = 9 - y
					} else {
						dist = y
					}
					explorationWeight := weights["explorationWeight"]
					if explorationWeight == 0 {
						explorationWeight = 5.0
					}
					score += float64(dist) * explorationWeight * aggression
				}
			} else {
				confidence := 0.0
				rememberedVal := val

				if piece.IsRevealed() {
					confidence = 1.0
				} else if memory != nil {
					entry := memory.Recall(pos)
					if entry != nil {
						confidence = entry.Confidence
						rememberedName := entry.Piece.GetType().GetName()
						if rVal, ok := weights[rememberedName]; ok {
							rememberedVal = rVal
						} else {
							rememberedVal = float64(entry.Piece.GetStrategicValue())
						}
					}
				}

				if confidence > 0 {
					score -= rememberedVal * confidence
				} else {
					score -= 40.0 * 0.5
				}
			}
		}
	}
	return score
}

// GetOpponent returns the opponent player on the board if found.
func GetOpponent(board *game.Board, ourPlayerID int) *game.Player {
	for y := range 10 {
		for x := range 10 {
			piece := board.GetPieceAt(game.NewPosition(x, y))
			if piece != nil && piece.GetOwner().GetID() != ourPlayerID {
				return piece.GetOwner()
			}
		}
	}
	return nil
}

// GetMoves gathers all valid moves for the specified player on the board.
func GetMoves(board *game.Board, player *game.Player) []game.Move {
	var allMoves []game.Move
	for y := range 10 {
		for x := range 10 {
			pos := game.NewPosition(x, y)
			piece := board.GetPieceAt(pos)
			if piece == nil || piece.GetOwner().GetID() != player.GetID() {
				continue
			}
			if !piece.CanMove() {
				continue
			}
			moves, err := board.ListMoves(pos)
			if err != nil {
				continue
			}
			for _, m := range moves {
				allMoves = append(allMoves, game.NewMove(m.GetFrom(), m.GetTo(), player))
			}
		}
	}
	return allMoves
}
