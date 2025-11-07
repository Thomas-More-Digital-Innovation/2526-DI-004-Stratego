package api

import (
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"encoding/json"
	"log"
	"sync"
	"time"
)

// WSHub manages all WebSocket connections for a game
type WSHub struct {
	clients       map[*WSClient]bool
	broadcast     chan []byte
	register      chan *WSClient
	unregister    chan *WSClient
	session       *game.GameSession
	gameType      string
	mutex         sync.RWMutex
	cleanupTimer  *time.Timer
	timerMutex    sync.Mutex
	cleanupPeriod time.Duration
}

func NewWSHub(session *game.GameSession, gameType string) *WSHub {
	return &WSHub{
		clients:       make(map[*WSClient]bool),
		broadcast:     make(chan []byte, 256),
		register:      make(chan *WSClient),
		unregister:    make(chan *WSClient),
		session:       session,
		gameType:      gameType,
		cleanupPeriod: 1 * time.Minute, // 1 minute grace period for reconnection
	}
}

// Run starts the hub's main loop
// If user disconnects, the hub will stop the game after 1 minute, if human is playing
func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			clientCount := len(h.clients)
			h.mutex.Unlock()

			if clientCount > 0 {
				h.cancelCleanupTimer()
			}

			go h.sendGameState(client)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			clientCount := len(h.clients)
			h.mutex.Unlock()

			if clientCount == 0 {
				switch h.gameType {
				case models.AiVsAi:
					// Stop AI vs AI games immediately - no point running without observers
					log.Printf("WSHub: All clients disconnected from AI vs AI game, stopping game immediately")
					h.session.Stop()

				case models.HumanVsAi:
					// Start cleanup timer for Human vs AI - allow reconnection grace period
					log.Printf("WSHub: All clients disconnected from Human vs AI game, starting cleanup timer")
					h.startCleanupTimer()

				case models.HumanVsHuman:
					// For Human vs Human, start timer to allow reconnection if both players leave
					log.Printf("WSHub: All clients disconnected from Human vs Human game, starting cleanup timer")
					h.startCleanupTimer()
				}
			}

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client is slow or disconnected
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// startCleanupTimer starts a timer to stop the game after the cleanup period
func (h *WSHub) startCleanupTimer() {
	h.timerMutex.Lock()
	defer h.timerMutex.Unlock()

	if h.cleanupTimer != nil {
		h.cleanupTimer.Stop()
	}

	log.Printf("WSHub: Starting cleanup timer for %s game (will stop in %v)", h.gameType, h.cleanupPeriod)

	h.cleanupTimer = time.AfterFunc(h.cleanupPeriod, func() {
		log.Printf("WSHub: Cleanup timer expired for %s game, stopping game", h.gameType)
		h.session.Stop()
	})
}

// cancelCleanupTimer cancels the cleanup timer if it's running
func (h *WSHub) cancelCleanupTimer() {
	h.timerMutex.Lock()
	defer h.timerMutex.Unlock()

	if h.cleanupTimer != nil {
		wasActive := h.cleanupTimer.Stop()
		if wasActive {
			log.Printf("WSHub: Cleanup timer cancelled for %s game (client reconnected)", h.gameType)
		}
		h.cleanupTimer = nil
	}
}

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
	for y := 0; y < 10; y++ {
		boardDTO[y] = make([]PieceDTO, 10)
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			if piece != nil {
				dto := PieceToDTO(piece, client.seatIndex)
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

func (h *WSHub) setupBoard() BoardStateMessage {
	session := h.session

	boardDTO := make([][]PieceDTO, 10)
	for y := range 10 {
		boardDTO[y] = make([]PieceDTO, 10)
	}

	// Place player 1 pieces in setup area (rows 6-9)
	player1Pieces := session.GetSetupPieces(0)
	idx := 0
	for y := 6; y <= 9; y++ {
		for x := range 10 {
			if idx < len(player1Pieces) {
				piece := player1Pieces[idx]
				dto := PieceToDTO(piece, 0) // Player 0 can see their own pieces
				dto.Position = PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dto
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
				dto := PieceToDTO(piece, -1) // Hide pieces
				dto.Position = PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dto
				idx++
			}
		}
	}

	return BoardStateMessage{
		Board:  boardDTO,
		Width:  10,
		Height: 10,
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

	historyMsg := MoveHistoryMessage{
		Moves: moveDTOs,
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
