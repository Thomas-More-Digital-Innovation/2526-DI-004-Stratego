package api

import (
	"digital-innovation/stratego/engine"
	"encoding/json"
	"fmt"
	"log"
)

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

// handleAnimationComplete processes animation complete message from client
func (c *WSClient) handleAnimationComplete() {
	log.Printf("Animation complete received from client %d", c.seatIndex)
	c.session.SignalAnimationComplete()
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

	if err := c.session.SwapSetupPieces(c.seatIndex, pos1, pos2); err != nil {
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
