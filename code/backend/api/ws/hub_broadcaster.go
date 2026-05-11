// Package ws provides WebSocket functionality for the game server.
package ws

import (
	"digital-innovation/gostrategy/api/dto"
	"digital-innovation/gostrategy/game"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/models"
	"encoding/json"
)

// BroadcastMessage sends a message to all connected clients
func (h *Hub) BroadcastMessage(msgType string, data any) {
	msg := dto.WSMessage{
		Type: msgType,
		Data: data,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		logging.Error("Error marshaling broadcast message", err)
		return
	}

	if h.IsStopped() {
		return
	}

	select {
	case h.broadcast <- jsonData:
	default:
		logging.Debug(logging.TagWeb, "Warning: Broadcast dropped for game %s (buffer full or hub stopping)", h.Session.ID)
	}
}

// BroadcastSetupBoard sends the setup board state to all clients
func (h *Hub) BroadcastSetupBoard() {
	boardMsg := h.setupBoard()

	h.BroadcastMessage(dto.MsgTypeBoardState, boardMsg)
}

// BroadcastGameState broadcasts the current game state to all clients
func (h *Hub) BroadcastGameState() {
	state := h.Session.GetGameState()

	var winnerName string
	var winCause string
	if state.WinnerID != nil {
		winner := h.Session.GetWinner()
		if winner != nil {
			winnerName = winner.GetName()
		}
		winCause = string(h.Session.GetWinCause())
	}

	h.BroadcastMessage(dto.MsgTypeGameState, dto.GameStateMessage{
		Round:              state.Round,
		CurrentPlayerID:    state.CurrentPlayerID,
		CurrentPlayerName:  state.CurrentPlayerName,
		IsGameOver:         state.IsGameOver,
		WinnerID:           state.WinnerID,
		WinnerName:         winnerName,
		WinCause:           winCause,
		Player1Score:       state.Player1Score,
		Player2Score:       state.Player2Score,
		WaitingForInput:    state.WaitingForInput,
		Paused:             state.Paused,
		MoveCount:          state.MoveCount,
		Player1AlivePieces: state.Player1AlivePieces,
		Player2AlivePieces: state.Player2AlivePieces,
		IsSetupPhase:       state.IsSetupPhase,
		Headless:           state.Headless,
		SetupRemainingSecs: state.SetupRemainingSecs,
		Player1Username:    state.Player1Username,
		Player2Username:    state.Player2Username,
	})
}

// BroadcastGameTransition broadcasts complete state after setup phase ends
// IsSetupPhase will be false now
func (h *Hub) BroadcastGameTransition() {
	h.BroadcastGameState()

	// Broadcast board state (pieces are now on the board)
	if h.GameType == models.AiVsAi {
		h.broadcastBoardStateRevealed()
	} else {
		h.broadcastBoardStatePerClient()
	}
}

// broadcastBoardStatePerClient sends personalized board to each client
func (h *Hub) broadcastBoardStatePerClient() {
	h.mutex.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mutex.RUnlock()

	for _, client := range clients {
		h.sendBoardState(client)
	}
}

// broadcastBoardStateRevealed sends board with all pieces revealed
func (h *Hub) broadcastBoardStateRevealed() {
	board := h.Session.GetBoard()
	field := board.GetField()

	boardDTO := make([][]dto.PieceDTO, 10)
	for y := 0; y < 10; y++ {
		boardDTO[y] = make([]dto.PieceDTO, 10)
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			if piece != nil && piece.IsAlive() {
				dtoPiece := dto.PieceToDTO(piece, piece.GetOwner().GetID())
				dtoPiece.Position = dto.PositionDTO{X: x, Y: y}
				dtoPiece.Revealed = true
				boardDTO[y][x] = dtoPiece
			}
		}
	}

	boardMsg := dto.BoardStateMessage{
		Board:    boardDTO,
		Width:    10,
		Height:   10,
		LastMove: h.Session.GetLastHistoricalMove(),
	}

	h.BroadcastMessage(dto.MsgTypeBoardState, boardMsg)
}

// BroadcastFullState sends complete game state and board to all clients
func (h *Hub) BroadcastFullState() {
	state := h.Session.GetGameState()
	h.BroadcastGameState()

	switch {
	case state.IsSetupPhase:
		h.BroadcastSetupBoard()
	case h.GameType == models.AiVsAi:
		h.broadcastBoardStateRevealed()
	default:
		h.broadcastBoardStatePerClient()
	}
}

// BroadcastCombat sends combat information to all clients
func (h *Hub) BroadcastCombat(combat *game.CombatResult) {
	if combat == nil || !combat.Occurred {
		return
	}

	attacker := combat.AttackerPiece
	defender := combat.DefenderPiece

	attackerDTO := dto.PieceToDTO(attacker, attacker.GetOwner().GetID())
	attackerDTO.Position = dto.PositionToDTO(combat.AttackerPosition)
	attackerDTO.Revealed = true

	defenderDTO := dto.PieceToDTO(defender, defender.GetOwner().GetID())
	defenderDTO.Position = dto.PositionToDTO(combat.DefenderPosition)
	defenderDTO.Revealed = true

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
