package api

import (
	"log"
	"time"
)

// monitorGame watches for game events and broadcasts them
func (s *GameServer) monitorGame(handler *GameSessionHandler, gameType string) {
	session := handler.Session
	hub := handler.Hub
	log.Printf("Starting game monitor for %s (type: %s)", handler.Session.ID, gameType)

	// Send initial state to all connected clients
	time.Sleep(100 * time.Millisecond) // Brief delay for clients to connect
	s.broadcastFullState(hub, gameType)

	// WAIT IN SETUP PHASE - WebSocket handlers will broadcast when user acts
	for session.IsSetupPhase() {
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("GameMonitor %s: Exiting setup phase, game starting", session.ID)

	// NOW we enter the game loop
	for {
		// Wait for a move notification with timeout
		if !session.WaitForMoveNotification(5 * time.Second) {
			// Timeout - check if game is over
			if !session.IsRunning() && session.GetGameState().IsGameOver {
				s.handleGameOver(session, hub)
				return
			}
			continue
		}

		// Move was executed
		log.Printf("Move executed in game %s", session.ID)

		// Check if combat occurred
		combat := session.GetLastCombat()
		hasCombat := combat != nil && combat.Occurred

		if hasCombat {
			log.Printf("Combat detected! Broadcasting combat data and waiting for animation")

			// Broadcast combat message (with piece info)
			s.broadcastCombat(hub, combat, gameType)

			// Wait for frontend animation to complete (3 second timeout)
			session.WaitForAnimationComplete(3 * time.Second)

			log.Printf("Animation complete, broadcasting updated state")

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
