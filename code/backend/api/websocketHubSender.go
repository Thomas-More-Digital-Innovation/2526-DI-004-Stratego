package api

import (
	"digital-innovation/stratego/models"
	"encoding/json"
	"log"
	"time"
)

// sendGameState sends the current game state to a specific client
func (h *WSHub) sendGameState(client *WSClient) {
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

	stateMsg := GameStateMessage{
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
	}

	msg := WSMessage{
		Type: MsgTypeGameState,
		Data: stateMsg,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling game state: %v", err)
		return
	}

	select {
	case client.send <- jsonData:
	case <-time.After(time.Second):
		log.Printf("Timeout sending game state to client")
	}

	h.sendBoardState(client)

	if !state.IsSetupPhase {
		h.sendMoveHistory(client)
	}
}

// sendBoardState sends the current board state to a specific client
func (h *WSHub) sendBoardState(client *WSClient) {
	if h.session.IsSetupPhase() {
		h.sendSetupBoard(client)
		return
	}

	board := h.session.GetBoard()
	field := board.GetField()

	boardDTO := make([][]PieceDTO, 10)
	for y := range 10 {
		boardDTO[y] = make([]PieceDTO, 10)
		for x := range 10 {
			boardDTO[y][x] = PieceDTO{OwnerID: -1}
			piece := field[y][x]
			if piece != nil {
				dto := PieceToDTO(piece, client.seatIndex)
				// Force reveal all pieces for AI vs AI spectators or when game is over
				if h.gameType == models.AiVsAi || h.session.GetGameState().IsGameOver {
					pieceType := piece.GetType()
					dto.Type = pieceType.GetName()
					dto.Rank = string(pieceType.GetRank())
					dto.Icon = pieceType.GetIcon()
					dto.Revealed = true
				}
				dto.Position = PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dto
			}
		}
	}

	boardMsg := BoardStateMessage{
		Board:  boardDTO,
		Width:  10,
		Height: 10,
	}

	msg := WSMessage{
		Type: MsgTypeBoardState,
		Data: boardMsg,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling board state: %v", err)
		return
	}

	select {
	case client.send <- jsonData:
	case <-time.After(time.Second):
		log.Printf("Timeout sending board state to client")
	}
}

// sendSetupBoard sends the setup board state to a specific client
func (h *WSHub) sendSetupBoard(client *WSClient) {
	boardMsg := h.setupBoard()
	msg := WSMessage{
		Type: MsgTypeBoardState,
		Data: boardMsg,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling setup board state: %v", err)
		return
	}

	select {
	case client.send <- jsonData:
	case <-time.After(time.Second):
		log.Printf("Timeout sending setup board state to client")
	}
}

// sendMoveHistory sends the move history to a specific client
func (h *WSHub) sendMoveHistory(client *WSClient) {
	g := h.session.GetGame()
	moveHistory := g.MoveHistory

	moveDTOs := make([]MoveDTO, len(moveHistory))
	for i, move := range moveHistory {
		moveDTOs[i] = MoveToDTO(move)
	}

	// Filter history if not AI vs AI and game is not over
	fullHistory := g.HistoricalHistory
	initialState := g.InitialState

	if h.gameType != models.AiVsAi && !g.IsGameOver() {
		// Filter initial state
		initialState = make([][]models.PieceData, len(g.InitialState))
		for y, row := range g.InitialState {
			initialState[y] = make([]models.PieceData, len(row))
			for x, piece := range row {
				p := piece
				if p.OwnerID != client.seatIndex && p.OwnerID != -1 && p.Type != "" {
					p.Type = ""
					p.Rank = ""
				}
				initialState[y][x] = p
			}
		}

		fullHistory = make([]models.HistoricalMove, len(g.HistoricalHistory))
		for i, m := range g.HistoricalHistory {
			move := m
			if move.Attacker != nil && move.Attacker.OwnerID != client.seatIndex && move.Result == models.ResultMove {
				move.Attacker = &models.PieceData{
					OwnerID: move.Attacker.OwnerID,
					Type:    "",
					Rank:    "",
				}
			}
			if move.Defender != nil && move.Defender.OwnerID != client.seatIndex && move.Result == models.ResultMove {
				move.Defender = &models.PieceData{
					OwnerID: move.Defender.OwnerID,
					Type:    "",
					Rank:    "",
				}
			}
			fullHistory[i] = move
		}
	}

	historyMsg := MoveHistoryMessage{
		Moves:        moveDTOs,
		FullHistory:  fullHistory,
		InitialState: initialState,
	}

	msg := WSMessage{
		Type: MsgTypeMoveHistory,
		Data: historyMsg,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling move history: %v", err)
		return
	}

	select {
	case client.send <- jsonData:
	case <-time.After(time.Second):
		log.Printf("Timeout sending move history to client")
	}
}
