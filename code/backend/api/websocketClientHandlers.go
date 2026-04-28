package api

import (
	"digital-innovation/stratego/engine"
	"digital-innovation/stratego/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"
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
		c.handleRandomizeSetup(baseMsg.Data)
	case MsgTypeStartGame:
		c.handleStartGame(baseMsg.Data)
	case MsgTypeLoadSetup:
		c.handleLoadSetup(baseMsg.Data)
	case MsgTypePause:
		c.handlePause()
	case MsgTypeUnpause:
		c.handleUnpause()
	case MsgTypeSetSpeed:
		c.handleSetSpeed(baseMsg.Data)
	case MsgTypeStep:
		c.handleStep()
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

	moves, err := c.session.GetAvailableMoves(c.seatIndex, pos)
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
	// Let the validation happen below so we can infer the player ID

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

	playerID := c.seatIndex
	if playerID < 0 {
		if c.hub.gameType == models.AiVsAi {
			if pos1.Y >= 6 {
				playerID = 0 // Bottom rows belong to player 0 (Red)
			} else {
				playerID = 1 // Top rows belong to player 1 (Blue)
			}
		} else {
			c.sendError("Spectators cannot swap pieces")
			return
		}
	}

	if err := c.session.SwapSetupPieces(playerID, pos1, pos2); err != nil {
		c.sendError(fmt.Sprintf("Failed to swap pieces: %v", err))
		return
	}

	log.Printf("Pieces swapped: %v <-> %v", pos1, pos2)

	c.hub.BroadcastSetupBoard()
}

// handleRandomizeSetup processes a randomize setup message
func (c *WSClient) handleRandomizeSetup(data interface{}) {
	if c.seatIndex < 0 && c.hub.gameType != models.AiVsAi {
		c.sendError("Spectators cannot randomize setup")
		return
	}

	targetPlayer := c.seatIndex
	if c.hub.gameType == models.AiVsAi {
		// Spectators in AI vs AI can specify which target player to randomize
		dataBytes, _ := json.Marshal(data)
		var msg RandomizeSetupMessage
		if err := json.Unmarshal(dataBytes, &msg); err == nil && msg.PlayerID != nil {
			targetPlayer = *msg.PlayerID
		} else {
			c.sendError("Invalid randomize setup message")
			return
		}
	}

	if targetPlayer < 0 {
		c.sendError("Target player required for AI randomization")
		return
	}

	if err := c.session.RandomizeSetup(targetPlayer); err != nil {
		c.sendError(fmt.Sprintf("Failed to randomize setup: %v", err))
		return
	}

	log.Printf("Setup randomized for player %d", targetPlayer)

	c.hub.BroadcastSetupBoard()
}

// handleStartGame processes a start game message
func (c *WSClient) handleStartGame(data interface{}) {
	if c.seatIndex < 0 && c.hub.gameType != models.AiVsAi {
		c.sendError("Spectators cannot start game")
		return
	}

	headless := false
	if data != nil {
		dataBytes, _ := json.Marshal(data)
		var msg StartGameMessage
		if err := json.Unmarshal(dataBytes, &msg); err == nil {
			headless = msg.Headless
		}
	}

	if err := c.session.StartGameFromSetup(headless); err != nil {
		c.sendError(fmt.Sprintf("Failed to start game: %v", err))
		return
	}

	log.Printf("Game started (client: %d, headless: %v)", c.seatIndex, headless)

	c.hub.BroadcastGameTransition()
}

// handleLoadSetup processes a load setup message from saved board setups
func (c *WSClient) handleLoadSetup(data interface{}) {
	if c.seatIndex < 0 && c.hub.gameType != models.AiVsAi {
		c.sendError("Spectators cannot load setups")
		return
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		c.sendError("Invalid load setup message format")
		return
	}

	var loadMsg LoadSetupMessage
	if err := json.Unmarshal(dataBytes, &loadMsg); err != nil {
		c.sendError("Invalid load setup message")
		return
	}

	targetPlayer := c.seatIndex
	if c.hub.gameType == models.AiVsAi && loadMsg.PlayerID != nil {
		targetPlayer = *loadMsg.PlayerID
	}

	if targetPlayer < 0 {
		c.sendError("Target player required for AI setup loading")
		return
	}

	var setupData []byte
	if len(loadMsg.SetupData) == 40 {
		setupData = []byte(loadMsg.SetupData)
	} else {
		var err error
		setupData, err = base64.StdEncoding.DecodeString(loadMsg.SetupData)
		if err != nil {
			c.sendError(fmt.Sprintf("Invalid setup data (expected 40 chars or base64): %v", err))
			return
		}
	}

	if err := c.session.LoadSetup(targetPlayer, setupData); err != nil {
		c.sendError(fmt.Sprintf("Failed to load setup: %v", err))
		return
	}

	log.Printf("Setup loaded for player %d", targetPlayer)

	c.hub.BroadcastSetupBoard()
}

// handlePause processes a pause game message
func (c *WSClient) handlePause() {
	if c.seatIndex < 0 && c.hub.gameType != models.AiVsAi {
		c.sendError("Spectators cannot pause the game")
		return
	}
	c.session.Pause()
	log.Printf("Game paused (client seat: %d, type: %s)", c.seatIndex, c.hub.gameType)
	c.hub.BroadcastGameState()
}

// handleUnpause processes an unpause game message
func (c *WSClient) handleUnpause() {
	if c.seatIndex < 0 && c.hub.gameType != models.AiVsAi {
		c.sendError("Spectators cannot unpause the game")
		return
	}
	c.session.Unpause()
	log.Printf("Game unpaused (client seat: %d, type: %s)", c.seatIndex, c.hub.gameType)
	c.hub.BroadcastGameState()
}

// handleSetSpeed processes a set speed message
func (c *WSClient) handleSetSpeed(data interface{}) {
	if c.seatIndex < 0 && c.hub.gameType != models.AiVsAi {
		c.sendError("Spectators cannot change speed")
		return
	}

	dataBytes, _ := json.Marshal(data)
	var msg SetSpeedMessage
	if err := json.Unmarshal(dataBytes, &msg); err != nil {
		c.sendError("Invalid set speed message")
		return
	}

	// Range check: 500ms to 5000ms
	speed := msg.SpeedMs
	if speed < 500 {
		speed = 500
	} else if speed > 5000 {
		speed = 5000
	}

	c.session.SetTurnDelay(time.Duration(speed) * time.Millisecond)
	log.Printf("Game speed set to %dms", speed)
}

// handleStep processes a manual step message
func (c *WSClient) handleStep() {
	if c.seatIndex < 0 && c.hub.gameType != models.AiVsAi {
		c.sendError("Spectators cannot step the game")
		return
	}

	if c.session.StepAI() {
		log.Printf("Manual AI step executed")
		c.hub.BroadcastGameState()
	} else {
		c.sendError("Failed to execute step (maybe already running or not AI turn)")
	}
}
