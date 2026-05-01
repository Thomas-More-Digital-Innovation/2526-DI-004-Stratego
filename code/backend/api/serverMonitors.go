package api

import (
	"digital-innovation/stratego/logging"
	"digital-innovation/stratego/utils"
	"fmt"
	"time"
)

// monitorGame watches for game events and broadcasts them
func (s *GameServer) monitorGame(handler *GameSessionHandler, gameType string) {
	session := handler.Session
	hub := handler.Hub

	// Send initial state to all connected clients
	time.Sleep(100 * time.Millisecond) // Brief delay for clients to connect
	s.broadcastFullState(hub, gameType)

	// WAIT IN SETUP PHASE - WebSocket handlers will broadcast when user acts
	select {
	case <-session.GetSetupCompleteChan():
		// Game starting
	case <-session.IsAbortedChan():
		logging.ErrorWith2Users(fmt.Sprintf("Game aborted during setup: %s", session.ID), session.Player1Username, utils.GetIntSafe(session.Player1UserID), session.Player2Username, utils.GetIntSafe(session.Player2UserID), nil)
		s.RemoveSession(session.ID)
		return
	}

	// NOW we enter the game loop
	for {
		// Wait for a move notification with timeout
		if !session.WaitForMoveNotification(5 * time.Second) {
			// Check if session was aborted while waiting
			if session.IsAborted() {
				logging.ErrorWith2Users(fmt.Sprintf("Game aborted during gameplay: %s", session.ID), session.Player1Username, utils.GetIntSafe(session.Player1UserID), session.Player2Username, utils.GetIntSafe(session.Player2UserID), nil)
				s.RemoveSession(session.ID)
				return
			}

			// Timeout - check if game is over
			if !session.IsRunning() && session.GetGameState().IsGameOver {
				s.handleGameOver(session, hub)
				return
			}
			continue
		}

		// Move was executed
		isHeadless := session.IsHeadless()
		if !isHeadless {
			logging.Debug(logging.TagGame, "Move executed in game %s", session.ID)
		}

		if isHeadless {
			state := session.GetGameState()
			if state.IsGameOver {
				time.Sleep(100 * time.Millisecond) // Brief delay to ensure runner completes
				s.handleGameOver(session, hub)
				return
			}
			session.AckMoveProcessed()
			continue
		}

		// Check if combat occurred
		combat := session.GetLastCombat()
		hasCombat := combat != nil && combat.Occurred

		if hasCombat {
			logging.Debug(logging.TagGame, "Combat detected! Broadcasting combat data and waiting for animation")

			// Broadcast combat message (with piece info)
			s.broadcastCombat(hub, combat, gameType)

			// Wait for frontend animation to complete (3 second timeout)
			session.WaitForAnimationComplete(3 * time.Second)

			logging.Debug(logging.TagGame, "Animation complete, broadcasting updated state")

			// Clear combat after animation
			session.ClearLastCombat()

			// NOW broadcast state after animation (winner moves to position, loser removed)
			s.broadcastFullState(hub, gameType)
		} else {
			// No combat - broadcast state immediately
			s.broadcastFullState(hub, gameType)
		}

		// Signal that move has been processed - GameRunner can continue
		session.AckMoveProcessed()

		// Check if game is over
		state := session.GetGameState()
		if state.IsGameOver {
			time.Sleep(500 * time.Millisecond) // Brief delay before game over message
			s.handleGameOver(session, hub)
			return
		}
	}
}
