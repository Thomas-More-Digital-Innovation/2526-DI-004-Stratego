package ws

import (
	"digital-innovation/gostrategy/internal/api/dto"
	"digital-innovation/gostrategy/internal/models"
	"digital-innovation/gostrategy/pkg/game"
)

// BroadcastCombat sends combat information to all clients
func (h *Hub) BroadcastCombat(combat *game.CombatResult) {
	if combat == nil || !combat.Occurred {
		return
	}

	attacker := combat.AttackerPiece
	defender := combat.DefenderPiece

	attackerDTO := dto.PieceToDTO(attacker, attacker.GetOwner().GetID(), true)
	attackerDTO.Position = dto.PositionToDTO(combat.AttackerPosition)

	defenderDTO := dto.PieceToDTO(defender, defender.GetOwner().GetID(), true)
	defenderDTO.Position = dto.PositionToDTO(combat.DefenderPosition)

	combatMsg := dto.CombatMessage{
		Attacker:     attackerDTO,
		Defender:     defenderDTO,
		AttackerWon:  attacker.IsAlive(),
		DefenderWon:  defender.IsAlive(),
		AttackerDied: !attacker.IsAlive(),
		DefenderDied: !defender.IsAlive(),
	}

	h.BroadcastMessage(dto.MsgTypeCombat, combatMsg)
}

// BroadcastMoveHistory sends personalized move history to each client
func (h *Hub) BroadcastMoveHistory() {
	h.mutex.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mutex.RUnlock()

	for _, client := range clients {
		h.sendMoveHistory(client)
	}
}

func (h *Hub) setupBoard() dto.BoardStateMessage {
	session := h.Session

	boardDTO := make([][]dto.PieceDTO, 10)
	for y := range 10 {
		boardDTO[y] = make([]dto.PieceDTO, 10)
	}

	// Place player 1 pieces in setup area (rows 6-9)
	player1Pieces := session.GetSetupPieces(0)
	idx := 0
	for y := 6; y <= 9; y++ {
		for x := range 10 {
			if idx < len(player1Pieces) {
				piece := player1Pieces[idx]
				forceReveal := h.GameType == models.AiVsAi
				dtoPiece := dto.PieceToDTO(piece, 0, forceReveal)
				dtoPiece.Position = dto.PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dtoPiece
				idx++
			}
		}
	}

	// Place player 2 pieces in setup area (rows 0-3)
	// Hide opponent pieces during setup
	player2Pieces := session.GetSetupPieces(1)
	idx = 0
	for y := 0; y <= 3; y++ {
		for x := range 10 {
			if idx < len(player2Pieces) {
				piece := player2Pieces[idx]
				viewerID := -1
				forceReveal := h.GameType == models.AiVsAi
				if forceReveal {
					viewerID = 1 // Show all pieces in AI vs AI
				}
				dtoPiece := dto.PieceToDTO(piece, viewerID, forceReveal)
				dtoPiece.Position = dto.PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dtoPiece
				idx++
			}
		}
	}

	return dto.BoardStateMessage{
		Board:  boardDTO,
		Width:  10,
		Height: 10,
	}
}
