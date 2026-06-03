// Package ws provides WebSocket functionality for the game server.
package ws

import (
	"digital-innovation/gostrategy/internal/api/dto"
	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/models"
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

	h.BroadcastMessage(dto.MsgTypeGameState, dto.BuildGameStateMessage(state, winnerName, winCause))
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

	isGameOver := h.Session.GetGameState().IsGameOver

	// Optimize: Build boards for each perspective (Red, Blue, Spectator)
	// seatIndex: 0 (Red), 1 (Blue), -1 (Spectator)
	perspectives := []int{0, 1, -1}
	cache := make(map[int][]byte)

	for _, seat := range perspectives {
		data := h.buildBoardStateJSON(seat, isGameOver)
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
func (h *Hub) buildBoardStateJSON(seatIndex int, isGameOver bool) []byte {
	board := h.Session.GetBoard()
	lastMove := h.Session.GetLastHistoricalMove()

	h.Session.RLock()
	field := board.GetField()

	forceReveal := h.GameType == models.AiVsAi || isGameOver
	boardDTO := dto.MapBoardToDTO(field, seatIndex, forceReveal)

	var filteredLastMove *models.HistoricalMove
	if lastMove != nil {
		fm := h.filterHistoricalMove(*lastMove, seatIndex, false, isGameOver)
		filteredLastMove = &fm
	}
	h.Session.RUnlock()

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
	lastMove := h.Session.GetLastHistoricalMove()

	h.Session.RLock()
	field := board.GetField()

	boardDTO := dto.MapBoardToDTO(field, 0, true)

	h.Session.RUnlock()

	boardMsg := dto.BoardStateMessage{
		Board:    boardDTO,
		Width:    10,
		Height:   10,
		LastMove: lastMove,
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
