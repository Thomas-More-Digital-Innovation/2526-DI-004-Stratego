package api

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"digital-innovation/stratego/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// TODO: Restrict in production
		return true
	},
}

// WSClient represents a WebSocket client connection
type WSClient struct {
	conn     *websocket.Conn
	send     chan []byte
	session  *game.GameSession
	playerID int // -1 for spectator, 0 or 1 for player
	hub      *WSHub
}

// WSHub manages all WebSocket connections for a game
type WSHub struct {
	clients    map[*WSClient]bool
	broadcast  chan []byte
	register   chan *WSClient
	unregister chan *WSClient
	session    *game.GameSession
	gameType   string
	mutex      sync.RWMutex
}

func NewWSHub(session *game.GameSession, gameType string) *WSHub {
	return &WSHub{
		clients:    make(map[*WSClient]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		session:    session,
		gameType:   gameType,
	}
}

// Run starts the hub's main loop
func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()

			// Send initial game state to new client
			go h.sendGameState(client)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Unlock()

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

// BroadcastMessage sends a message to all connected clients
func (h *WSHub) BroadcastMessage(msgType string, data interface{}) {
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
	session := h.session

	// Create empty board
	boardDTO := make([][]PieceDTO, 10)
	for y := 0; y < 10; y++ {
		boardDTO[y] = make([]PieceDTO, 10)
	}

	// Place player 1 pieces in setup area (rows 6-9)
	player1Pieces := session.GetSetupPieces(0)
	idx := 0
	for y := 6; y <= 9; y++ {
		for x := 0; x < 10; x++ {
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
		for x := 0; x < 10; x++ {
			if idx < len(player2Pieces) {
				piece := player2Pieces[idx]
				dto := PieceToDTO(piece, -1) // Hide pieces
				dto.Position = PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dto
				idx++
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

// BroadcastGameTransition broadcasts complete state after setup phase ends
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

	// Broadcast updated game state (with isSetupPhase = false)
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
		IsSetupPhase:       state.IsSetupPhase, // This will be false now
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

	// Send personalized board state to each client
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

	// Also send board state
	h.sendBoardState(client)
}

// sendBoardState sends the current board state to a specific client
func (h *WSHub) sendBoardState(client *WSClient) {
	// Check if we're in setup phase
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
				dto := PieceToDTO(piece, client.playerID)
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
	session := h.session

	// Create empty board
	boardDTO := make([][]PieceDTO, 10)
	for y := 0; y < 10; y++ {
		boardDTO[y] = make([]PieceDTO, 10)
	}

	// Place player 1 pieces in setup area (rows 6-9)
	player1Pieces := session.GetSetupPieces(0)
	idx := 0
	for y := 6; y <= 9; y++ {
		for x := 0; x < 10; x++ {
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
		for x := 0; x < 10; x++ {
			if idx < len(player2Pieces) {
				piece := player2Pieces[idx]
				dto := PieceToDTO(piece, -1) // Hide pieces
				dto.Position = PositionDTO{X: x, Y: y}
				boardDTO[y][x] = dto
				idx++
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
		log.Printf("Error marshaling setup board state: %v", err)
		return
	}

	select {
	case client.send <- jsonData:
	case <-time.After(time.Second):
		log.Printf("Timeout sending setup board state to client")
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *WSClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		log.Printf("Error setting read deadline: %v", err)
		return
	}
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			log.Printf("Error setting read deadline: %v", err)
		}
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		c.handleMessage(message)
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *WSClient) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				log.Printf("Error setting write deadline: %v", err)
				return
			}
			if !ok {
				// Hub closed the channel
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("Error writing close message: %v", err)
				}
				return
			}

			// Send each message separately instead of concatenating
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				log.Printf("Error setting write deadline: %v", err)
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (c *WSClient) handleMessage(message []byte) {
	var baseMsg WSMessage
	if err := json.Unmarshal(message, &baseMsg); err != nil {
		c.sendError("Invalid message format")
		return
	}

	switch baseMsg.Type {
	case MsgTypeMove:
		c.handleMove(baseMsg.Data)
	case MsgTypeGetValidMoves:
		c.handleGetValidMoves(baseMsg.Data)
	case MsgTypePing:
		c.sendPong()
	case MsgTypeAnimationComplete:
		c.handleAnimationComplete()
	case MsgTypeSwapPieces:
		c.handleSwapPieces(baseMsg.Data)
	case MsgTypeRandomizeSetup:
		c.handleRandomizeSetup()
	case MsgTypeStartGame:
		c.handleStartGame()
	default:
		c.sendError("Unknown message type")
	}
}

// handleMove processes a move message from the client
func (c *WSClient) handleMove(data interface{}) {
	// Only players can move, not spectators
	if c.playerID < 0 {
		c.sendError("Spectators cannot make moves")
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.sendError("Invalid move data")
		return
	}

	var moveMsg MoveMessage
	if err := json.Unmarshal(jsonData, &moveMsg); err != nil {
		c.sendError("Invalid move format")
		return
	}

	// Convert to engine types
	from := engine.NewPosition(moveMsg.From.X, moveMsg.From.Y)
	to := engine.NewPosition(moveMsg.To.X, moveMsg.To.Y)

	// Get the player making the move
	g := c.session.GetGame()
	player := g.Players[c.playerID]
	move := engine.NewMove(from, to, player)

	// Submit move to game session
	err = c.session.SubmitMove(c.playerID, move)
	if err != nil {
		c.sendMoveResult(false, err.Error())
		return
	}

	// Move accepted - result will be broadcast when processed
	c.sendMoveResult(true, "")
}

// handleGetValidMoves processes a request for valid moves for a piece
func (c *WSClient) handleGetValidMoves(data interface{}) {
	// Only players can request moves for their pieces
	if c.playerID < 0 {
		c.sendError("Spectators cannot request valid moves")
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.sendError("Invalid request data")
		return
	}

	var reqMsg GetValidMovesMessage
	if err := json.Unmarshal(jsonData, &reqMsg); err != nil {
		c.sendError("Invalid request format")
		return
	}

	// Convert to engine type
	pos := engine.NewPosition(reqMsg.Position.X, reqMsg.Position.Y)

	// Get available moves
	moves, err := c.session.GetAvailableMoves(pos)
	if err != nil {
		c.sendError(err.Error())
		return
	}

	// Convert to DTOs
	validMoveDTOs := make([]PositionDTO, len(moves))
	for i, move := range moves {
		validMoveDTOs[i] = PositionToDTO(move.GetTo())
	}

	// Send response
	msg := WSMessage{
		Type: MsgTypeValidMoves,
		Data: ValidMovesMessage{
			Position:   reqMsg.Position,
			ValidMoves: validMoveDTOs,
		},
	}

	jsonResponse, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling valid moves: %v", err)
		return
	}

	c.send <- jsonResponse
}

// sendMoveResult sends a move result message
func (c *WSClient) sendMoveResult(success bool, error string) {
	result := MoveResultMessage{
		Success: success,
		Error:   error,
	}

	msg := WSMessage{
		Type: MsgTypeMoveResult,
		Data: result,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling move result: %v", err)
		return
	}

	c.send <- jsonData
}

// handleAnimationComplete processes animation complete message from client
func (c *WSClient) handleAnimationComplete() {
	log.Printf("Animation complete received from client %d", c.playerID)
	c.session.SignalAnimationComplete()
}

// sendError sends an error message
func (c *WSClient) sendError(errMsg string) {
	msg := WSMessage{
		Type: MsgTypeError,
		Data: ErrorMessage{Error: errMsg},
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling error message: %v", err)
		return
	}

	c.send <- jsonData
}

// sendPong sends a pong response
func (c *WSClient) sendPong() {
	msg := WSMessage{
		Type: MsgTypePong,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling pong: %v", err)
		return
	}

	c.send <- jsonData
}

// handleSwapPieces processes a swap pieces message during setup
func (c *WSClient) handleSwapPieces(data interface{}) {
	if c.playerID < 0 {
		c.sendError("Spectators cannot swap pieces")
		return
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		c.sendError("Invalid swap message format")
		return
	}

	var swapMsg SwapPiecesMessage
	if err := json.Unmarshal(dataBytes, &swapMsg); err != nil {
		c.sendError("Invalid swap message")
		return
	}

	pos1 := engine.NewPosition(swapMsg.Pos1.X, swapMsg.Pos1.Y)
	pos2 := engine.NewPosition(swapMsg.Pos2.X, swapMsg.Pos2.Y)

	if err := c.session.SwapPieces(c.playerID, pos1, pos2); err != nil {
		c.sendError(fmt.Sprintf("Failed to swap pieces: %v", err))
		return
	}

	log.Printf("Pieces swapped: %v <-> %v", pos1, pos2)

	// Broadcast updated setup board to all clients
	c.hub.BroadcastSetupBoard()
}

// handleRandomizeSetup processes a randomize setup message
func (c *WSClient) handleRandomizeSetup() {
	if c.playerID < 0 {
		c.sendError("Spectators cannot randomize setup")
		return
	}

	if err := c.session.RandomizeSetup(c.playerID); err != nil {
		c.sendError(fmt.Sprintf("Failed to randomize setup: %v", err))
		return
	}

	log.Printf("Setup randomized for player %d", c.playerID)

	// Broadcast updated setup board to all clients
	c.hub.BroadcastSetupBoard()
}

// handleStartGame processes a start game message
func (c *WSClient) handleStartGame() {
	if c.playerID < 0 {
		c.sendError("Spectators cannot start game")
		return
	}

	if err := c.session.StartGameFromSetup(); err != nil {
		c.sendError(fmt.Sprintf("Failed to start game: %v", err))
		return
	}

	log.Printf("Game started by player %d", c.playerID)

	// Broadcast the transition out of setup phase to all clients
	c.hub.BroadcastGameTransition()
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request, session *game.GameSession, hub *WSHub, playerID int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &WSClient{
		conn:     conn,
		send:     make(chan []byte, 256),
		session:  session,
		playerID: playerID,
		hub:      hub,
	}

	hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}
