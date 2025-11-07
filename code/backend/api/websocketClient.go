package api

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/game"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// WSClient represents a WebSocket client connection
type WSClient struct {
	conn      *websocket.Conn
	send      chan []byte
	session   *game.GameSession
	seatIndex int // -1 for spectator, 0 or 1 for player
	hub       *WSHub
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
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("Error writing close message: %v", err)
				}
				return
			}

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
func (c *WSClient) handleMove(data any) {
	if c.seatIndex < 0 {
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

	from := engine.NewPosition(moveMsg.From.X, moveMsg.From.Y)
	to := engine.NewPosition(moveMsg.To.X, moveMsg.To.Y)

	g := c.session.GetGame()
	player := g.Players[c.seatIndex]
	move := engine.NewMove(from, to, player)

	err = c.session.SubmitMove(c.seatIndex, move)
	if err != nil {
		c.sendMoveResult(false, err.Error())
		return
	}

	c.sendMoveResult(true, "")
}

// handleGetValidMoves processes a request for valid moves for a piece
func (c *WSClient) handleGetValidMoves(data any) {
	if c.seatIndex < 0 {
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

	pos := engine.NewPosition(reqMsg.Position.X, reqMsg.Position.Y)

	moves, err := c.session.GetAvailableMoves(pos)
	if err != nil {
		c.sendError(err.Error())
		return
	}

	validMoveDTOs := make([]PositionDTO, len(moves))
	for i, move := range moves {
		validMoveDTOs[i] = PositionToDTO(move.GetTo())
	}

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
	log.Printf("Animation complete received from client %d", c.seatIndex)
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
	if c.seatIndex < 0 {
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

	if err := c.session.SwapPieces(c.seatIndex, pos1, pos2); err != nil {
		c.sendError(fmt.Sprintf("Failed to swap pieces: %v", err))
		return
	}

	log.Printf("Pieces swapped: %v <-> %v", pos1, pos2)

	c.hub.BroadcastSetupBoard()
}

// handleRandomizeSetup processes a randomize setup message
func (c *WSClient) handleRandomizeSetup() {
	if c.seatIndex < 0 {
		c.sendError("Spectators cannot randomize setup")
		return
	}

	if err := c.session.RandomizeSetup(c.seatIndex); err != nil {
		c.sendError(fmt.Sprintf("Failed to randomize setup: %v", err))
		return
	}

	log.Printf("Setup randomized for player %d", c.seatIndex)

	c.hub.BroadcastSetupBoard()
}

// handleStartGame processes a start game message
func (c *WSClient) handleStartGame() {
	if c.seatIndex < 0 {
		c.sendError("Spectators cannot start game")
		return
	}

	if err := c.session.StartGameFromSetup(); err != nil {
		c.sendError(fmt.Sprintf("Failed to start game: %v", err))
		return
	}

	log.Printf("Game started by player %d", c.seatIndex)

	c.hub.BroadcastGameTransition()
}
