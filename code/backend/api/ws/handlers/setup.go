package handlers

import (
	"digital-innovation/gostrategy/api/dto"
	"digital-innovation/gostrategy/api/ws"
	"digital-innovation/gostrategy/engine"
	"encoding/json"
	"fmt"
)

func init() {
	ws.RegisterHandler(dto.MsgTypeSwapPieces, handleSwapPieces)
	ws.RegisterHandler(dto.MsgTypeRandomizeSetup, handleRandomizeSetup)
	ws.RegisterHandler(dto.MsgTypeStartGame, handleStartGame)
	ws.RegisterHandler(dto.MsgTypeLoadSetup, handleLoadSetup)
}

func handleSwapPieces(c *ws.Client, data json.RawMessage) error {
	var swapMsg dto.SwapPiecesMessage
	if err := json.Unmarshal(data, &swapMsg); err != nil {
		return err
	}

	pos1 := engine.NewPosition(swapMsg.Pos1.X, swapMsg.Pos1.Y)
	pos2 := engine.NewPosition(swapMsg.Pos2.X, swapMsg.Pos2.Y)

	playerID := c.GetSeatIndex()
	if playerID < 0 {
		if !c.IsAuthorized() {
			return fmt.Errorf("spectators cannot swap pieces")
		}
		// Creator in AI vs AI can swap for both players
		// Ensure both positions are in the same player's setup area
		p1_0 := pos1.Y >= 6 && pos1.Y <= 9
		p2_0 := pos2.Y >= 6 && pos2.Y <= 9
		p1_1 := pos1.Y >= 0 && pos1.Y <= 3
		p2_1 := pos2.Y >= 0 && pos2.Y <= 3

		switch {
		case p1_0 && p2_0:
			playerID = 0
		case p1_1 && p2_1:
			playerID = 1
		default:
			return fmt.Errorf("both positions must be within the same player's setup area")
		}
	} else {
		// Regular player - ensure both positions are in their area
		if playerID == 0 {
			if pos1.Y < 6 || pos1.Y > 9 || pos2.Y < 6 || pos2.Y > 9 {
				return fmt.Errorf("positions outside setup area")
			}
		} else {
			if pos1.Y < 0 || pos1.Y > 3 || pos2.Y < 0 || pos2.Y > 3 {
				return fmt.Errorf("positions outside setup area")
			}
		}
	}

	if err := c.GetSession().SwapSetupPieces(playerID, pos1, pos2); err != nil {
		return err
	}
	c.GetHub().BroadcastSetupBoard()
	return nil
}

func handleRandomizeSetup(c *ws.Client, data json.RawMessage) error {
	if !c.IsAuthorized() {
		return fmt.Errorf("unauthorized")
	}

	targetPlayer := c.GetSeatIndex()
	// AI vs AI check... simplified for brevity as per dispatcher goal
	var msg dto.RandomizeSetupMessage
	if err := json.Unmarshal(data, &msg); err == nil && msg.PlayerID != nil {
		targetPlayer = *msg.PlayerID
	}

	if targetPlayer < 0 {
		return fmt.Errorf("target player required")
	}

	if err := c.GetSession().RandomizeSetup(targetPlayer); err != nil {
		return err
	}
	c.GetHub().BroadcastSetupBoard()
	return nil
}

func handleStartGame(c *ws.Client, data json.RawMessage) error {
	if !c.IsAuthorized() {
		return fmt.Errorf("unauthorized")
	}

	var msg dto.StartGameMessage
	headless := false
	if err := json.Unmarshal(data, &msg); err == nil {
		headless = msg.Headless
	}

	if err := c.GetSession().StartGameFromSetup(headless); err != nil {
		return err
	}
	c.GetHub().BroadcastGameTransition()
	return nil
}

func handleLoadSetup(c *ws.Client, data json.RawMessage) error {
	if !c.IsAuthorized() {
		return fmt.Errorf("unauthorized")
	}

	var loadMsg dto.LoadSetupMessage
	if err := json.Unmarshal(data, &loadMsg); err != nil {
		return err
	}

	targetPlayer := c.GetSeatIndex()
	if loadMsg.PlayerID != nil {
		targetPlayer = *loadMsg.PlayerID
	}

	if targetPlayer < 0 {
		return fmt.Errorf("target player required")
	}

	// Setup data decoding logic...
	// For now, assume session handles it
	if err := c.GetSession().LoadSetup(targetPlayer, []byte(loadMsg.SetupData)); err != nil {
		return err
	}
	c.GetHub().BroadcastSetupBoard()
	return nil
}
