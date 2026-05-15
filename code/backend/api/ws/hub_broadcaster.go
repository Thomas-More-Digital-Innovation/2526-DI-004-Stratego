// Package ws provides WebSocket functionality for the game server.
package ws

import (
	"digital-innovation/gostrategy/api/dto"
	"digital-innovation/gostrategy/game"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/models"
	"encoding/json"
	"time"
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

	// Optimize: Build boards for each perspective (Red, Blue, Spectator)
	// seatIndex: 0 (Red), 1 (Blue), -1 (Spectator)
	perspectives := []int{0, 1, -1}
	cache := make(map[int][]byte)

	for _, seat := range perspectives {
		data := h.buildBoardStateJSON(seat)
		if data != nil {
			cache[seat] = data
		}
	}

	for _, client := range clients {
		if data, ok := cache[client.SeatIndex]; ok {
			client.Send(data, 500*time.Millisecond)
		} else {
			// Fallback for unexpected seat indices
			h.sendBoardState(client)
		}
	}
}

// buildBoardStateJSON helper to build and marshal board state for a perspective
func (h *Hub) buildBoardStateJSON(seatIndex int) []byte {
	board := h.Session.GetBoard()
	field := board.GetField()

	boardDTO := make([][]dto.PieceDTO, 10)
	for y := range 10 {
		boardDTO[y] = make([]dto.PieceDTO, 10)
		for x := range 10 {
			boardDTO[y][x] = dto.PieceDTO{OwnerID: -1}
			piece := field[y][x]
			if piece != nil {
				forceReveal := h.GameType == models.AiVsAi || h.Session.GetGameState().IsGameOver
				dtoPiece := dto.PieceToDTO(piece, seatIndex, forceReveal)
				dtoPiece.Position = dto.PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dtoPiece
			}
		}
	}

	lastMove := h.Session.GetLastHistoricalMove()
	var filteredLastMove *models.HistoricalMove
	if lastMove != nil {
		fm := h.filterHistoricalMove(*lastMove, seatIndex, false)
		filteredLastMove = &fm
	}

	boardMsg := dto.BoardStateMessage{
		Board:    boardDTO,
		Width:    10,
		Height:   10,
		LastMove: filteredLastMove,
	}

	msg := dto.WSMessage{
		Type: dto.MsgTypeBoardState,
		Data: boardMsg,
	}

	jsonData, _ := json.Marshal(msg)
	return jsonData
}

// broadcastBoardStateRevealed sends board with all pieces revealed
func (h *Hub) broadcastBoardStateRevealed() {
	board := h.Session.GetBoard()
	field := board.GetField()

	boardDTO := make([][]dto.PieceDTO, 10)
	for y := range 10 {
		boardDTO[y] = make([]dto.PieceDTO, 10)
		for x := range 10 {
			piece := field[y][x]
			if piece != nil && piece.IsAlive() {
				dtoPiece := dto.PieceToDTO(piece, piece.GetOwner().GetID(), true)
				dtoPiece.Position = dto.PositionDTO{X: x, Y: y}
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
