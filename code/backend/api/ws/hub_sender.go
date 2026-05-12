// Package ws provides WebSocket functionality for the game server.
package ws

import (
	"digital-innovation/gostrategy/api/dto"
	"digital-innovation/gostrategy/logging"
	"digital-innovation/gostrategy/models"
	"encoding/json"
	"time"
)

// sendGameState sends the current game state to a specific client
func (h *Hub) sendGameState(client *Client) {
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

	stateMsg := dto.GameStateMessage{
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
	}

	msg := dto.WSMessage{
		Type: dto.MsgTypeGameState,
		Data: stateMsg,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		logging.ErrorWithUser("Error marshaling game state", client.Username, client.UserID, err)
		return
	}

	if !client.Send(jsonData, time.Second) {
		logging.DebugWithUser(logging.TagWeb, client.Username, client.UserID, "Timeout sending game state")
	}

	h.sendBoardState(client)

	if !state.IsSetupPhase {
		h.sendMoveHistory(client)
	}
}

// sendBoardState sends the current board state to a specific client
func (h *Hub) sendBoardState(client *Client) {
	if h.Session.IsSetupPhase() {
		h.sendSetupBoard(client)
		return
	}

	board := h.Session.GetBoard()
	field := board.GetField()

	boardDTO := make([][]dto.PieceDTO, 10)
	for y := range 10 {
		boardDTO[y] = make([]dto.PieceDTO, 10)
		for x := range 10 {
			boardDTO[y][x] = dto.PieceDTO{OwnerID: -1}
			piece := field[y][x]
			if piece != nil {
				// Force reveal all pieces for AI vs AI spectators or when game is over
				forceReveal := h.GameType == models.AiVsAi || h.Session.GetGameState().IsGameOver
				dtoPiece := dto.PieceToDTO(piece, client.SeatIndex, forceReveal)
				dtoPiece.Position = dto.PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dtoPiece
			}
		}
	}

	lastMove := h.Session.GetLastHistoricalMove()
	var filteredLastMove *models.HistoricalMove
	if lastMove != nil {
		// For the board state's LastMove, we don't force filter combat so the visualization works
		fm := h.filterHistoricalMove(*lastMove, client.SeatIndex, false)
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

	jsonData, err := json.Marshal(msg)
	if err != nil {
		logging.ErrorWithUser("Error marshaling board state", client.Username, client.UserID, err)
		return
	}

	if !client.Send(jsonData, time.Second) {
		logging.DebugWithUser(logging.TagWeb, client.Username, client.UserID, "Timeout sending board state")
	}
}

// sendSetupBoard sends the setup board state to a specific client
func (h *Hub) sendSetupBoard(client *Client) {
	boardMsg := h.setupBoard()
	msg := dto.WSMessage{
		Type: dto.MsgTypeBoardState,
		Data: boardMsg,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		logging.ErrorWithUser("Error marshaling setup board state", client.Username, client.UserID, err)
		return
	}

	if !client.Send(jsonData, time.Second) {
		logging.DebugWithUser(logging.TagWeb, client.Username, client.UserID, "Timeout sending setup board state")
	}
}

// sendMoveHistory sends the move history to a specific client
func (h *Hub) sendMoveHistory(client *Client) {
	g := h.Session.GetGame()
	moveHistory := g.MoveHistory

	moveDTOs := make([]dto.MoveDTO, len(moveHistory))
	for i, move := range moveHistory {
		moveDTOs[i] = dto.MoveToDTO(move)
	}

	// Filter history if not AI vs AI and game is not over
	fullHistory := g.HistoricalHistory
	initialState := g.InitialState

	if h.GameType != models.AiVsAi && !g.IsGameOver() {
		// Filter initial state
		initialState = make([][]models.PieceData, len(g.InitialState))
		for y, row := range g.InitialState {
			initialState[y] = make([]models.PieceData, len(row))
			for x, piece := range row {
				p := piece
				if p.OwnerID != client.SeatIndex && p.OwnerID != -1 && p.Type != "" {
					p.Type = ""
					p.Rank = ""
				}
				initialState[y][x] = p
			}
		}

		fullHistory = make([]models.HistoricalMove, len(g.HistoricalHistory))
		for i, m := range g.HistoricalHistory {
			// For history, we force filter combat to prevent leaking piece ranks in a live game
			fullHistory[i] = h.filterHistoricalMove(m, client.SeatIndex, true)
		}
	}

	historyMsg := dto.MoveHistoryMessage{
		Moves:        moveDTOs,
		FullHistory:  fullHistory,
		InitialState: initialState,
	}

	msg := dto.WSMessage{
		Type: dto.MsgTypeMoveHistory,
		Data: historyMsg,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		logging.ErrorWithUser("Error marshaling move history", client.Username, client.UserID, err)
		return
	}

	if !client.Send(jsonData, time.Second) {
		logging.DebugWithUser(logging.TagWeb, client.Username, client.UserID, "Timeout sending move history")
	}
}

func (h *Hub) filterHistoricalMove(m models.HistoricalMove, seatIndex int, forceFilterCombat bool) models.HistoricalMove {
	if h.GameType == models.AiVsAi || h.Session.GetGameState().IsGameOver {
		return m
	}

	move := m
	if move.Attacker != nil && move.Attacker.OwnerID != seatIndex {
		if move.Result == models.ResultMove || forceFilterCombat {
			move.Attacker = &models.PieceData{
				OwnerID: move.Attacker.OwnerID,
				Type:    "",
				Rank:    "",
			}
		}
	}
	if move.Defender != nil && move.Defender.OwnerID != seatIndex {
		if move.Result == models.ResultMove || forceFilterCombat {
			move.Defender = &models.PieceData{
				OwnerID: move.Defender.OwnerID,
				Type:    "",
				Rank:    "",
			}
		}
	}
	return move
}
