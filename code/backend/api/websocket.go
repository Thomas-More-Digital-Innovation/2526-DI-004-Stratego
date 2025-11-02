package api

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

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
	mutex    sync.RWMutex
}

// WSHub manages all WebSocket connections for a game
type WSHub struct {
	clients    map[*WSClient]bool
	broadcast  chan []byte
	register   chan *WSClient
	unregister chan *WSClient
	session    *game.GameSession
	mutex      sync.RWMutex
}

func NewWSHub(session *game.GameSession) *WSHub {
	return &WSHub{
		clients:    make(map[*WSClient]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		session:    session,
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

// readPump pumps messages from the websocket connection to the hub
func (c *WSClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
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
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Send each message separately instead of concatenating
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
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
