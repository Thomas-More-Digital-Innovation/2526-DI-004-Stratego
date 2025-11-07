package api

import (
	"digital-innovation/stratego/models"
	"encoding/json"
	"log"
)

// BroadcastMessage sends a message to all connected clients
func (h *WSHub) BroadcastMessage(msgType string, data any) {
	msg := WSMessage{
		Type: msgType,
		Data: data,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.broadcast <- jsonData
}

// BroadcastSetupBoard sends the setup board state to all clients
func (h *WSHub) BroadcastSetupBoard() {
	boardMsg := h.setupBoard()

	h.BroadcastMessage(MsgTypeBoardState, boardMsg)
}

// BroadcastGameTransition broadcasts complete state after setup phase ends
// IsSetupPhase will be false now
func (h *WSHub) BroadcastGameTransition() {
	state := h.session.GetGameState()

	var winnerName string
	var winCause string
	if state.WinnerID != nil {
		winner := h.session.GetWinner()
		if winner != nil {
			winnerName = winner.GetName()
		}
		winCause = string(h.session.GetWinCause())
	}

	h.BroadcastMessage(MsgTypeGameState, GameStateMessage{
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
		MoveCount:          state.MoveCount,
		Player1AlivePieces: state.Player1AlivePieces,
		Player2AlivePieces: state.Player2AlivePieces,
		IsSetupPhase:       state.IsSetupPhase,
	})

	// Broadcast board state (pieces are now on the board)
	if h.gameType == models.AiVsAi {
		h.broadcastBoardStateRevealed()
	} else {
		h.broadcastBoardStatePerClient()
	}
}

// broadcastBoardStatePerClient sends personalized board to each client
func (h *WSHub) broadcastBoardStatePerClient() {
	h.mutex.RLock()
	clients := make([]*WSClient, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mutex.RUnlock()

	for _, client := range clients {
		h.sendBoardState(client)
	}
}

// broadcastBoardStateRevealed sends board with all pieces revealed
func (h *WSHub) broadcastBoardStateRevealed() {
	board := h.session.GetBoard()
	field := board.GetField()

	boardDTO := make([][]PieceDTO, 10)
	for y := 0; y < 10; y++ {
		boardDTO[y] = make([]PieceDTO, 10)
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			if piece != nil && piece.IsAlive() {
				dto := PieceToDTO(piece, piece.GetOwner().GetID())
				dto.Position = PositionDTO{X: x, Y: y}
				dto.Revealed = true
				boardDTO[y][x] = dto
			}
		}
	}

	boardMsg := BoardStateMessage{
		Board:  boardDTO,
		Width:  10,
		Height: 10,
	}

	h.BroadcastMessage(MsgTypeBoardState, boardMsg)
}
