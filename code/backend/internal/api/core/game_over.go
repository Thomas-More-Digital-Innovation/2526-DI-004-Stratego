package core

import (
	"context"
	"digital-innovation/gostrategy/internal/api/ws"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/game"
	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/models"
	"fmt"
	"time"
)

// SaveGameResult persists a game session's final outcome to the database and updates user stats.
// AI vs AI games do not declare a winner and are skipped from the stats update.
func (s *GameServer) SaveGameResult(ctx context.Context, session *game.Session, gameType string) error {
	state := session.GetGameState()
	g := session.GetGame()
	now := time.Now()

	if err := db.SaveGame(ctx, session.ID, session.Player1UserID, session.Player2UserID, gameType, g.InitialState, state.WinnerID, session.StartTime, now); err != nil {
		logging.Error("Failed to save game to database", err)
		return err
	}

	if err := db.SaveMoves(ctx, session.ID, g.HistoricalHistory); err != nil {
		logging.Error("Failed to save game moves to database", err)
	}

	if gameType != models.AiVsAi {
		durationSecs := now.Sub(session.StartTime).Seconds()
		moveCount := len(g.HistoricalHistory)

		updateStats := func(userID *int, playerIndex int) {
			if userID != nil {
				won := state.WinnerID != nil && *state.WinnerID == playerIndex
				if err := db.UpdateUserStats(ctx, *userID, won, moveCount, durationSecs); err != nil {
					logging.Error(fmt.Sprintf("Failed to update stats for Player %d", playerIndex+1), err)
				}
			}
		}

		updateStats(session.Player1UserID, 0)
		updateStats(session.Player2UserID, 1)
	}

	return nil
}

// resolveUserSeat maps a userID to the user's player seat index and their opponent's player seat index.
// NOTE: This method internally acquires a session read lock (session.GetPlayerIDs()).
// To avoid self-deadlocks, do not call this function while already holding a session lock.
func resolveUserSeat(session *game.Session, userID int) (playerIndex int, opponentIndex int, err error) {
	p1, p2 := session.GetPlayerIDs()
	switch {
	case p1 != nil && *p1 == userID:
		return 0, 1, nil
	case p2 != nil && *p2 == userID:
		return 1, 0, nil
	default:
		return -1, -1, fmt.Errorf("user is not a participant in this game session")
	}
}

// getPlayerAndOpponent safely retrieves the player and opponent objects from the game engine, handling out-of-bound indices.
func getPlayerAndOpponent(gameObj *game.Game, playerIndex int, opponentIndex int) (player *game.Player, opponent *game.Player) {
	resolvedPlayerIndex := 0
	if playerIndex >= 0 && playerIndex < len(gameObj.Players) {
		resolvedPlayerIndex = playerIndex
	}

	resolvedOpponentIndex := 1 - resolvedPlayerIndex
	if opponentIndex >= 0 && opponentIndex < len(gameObj.Players) {
		resolvedOpponentIndex = opponentIndex
	}

	return gameObj.Players[resolvedPlayerIndex], gameObj.Players[resolvedOpponentIndex]
}

// ResignGame handles a player resigning from a game session.
func (s *GameServer) ResignGame(gameID string, resigningUserID int) error {
	handler, exists := s.GetSession(gameID)
	if !exists {
		return fmt.Errorf("game session %s not found", gameID)
	}

	resigningPlayerIndex, opponentIndex, err := resolveUserSeat(handler.Session, resigningUserID)
	if err != nil {
		return err
	}

	gameObj := handler.Session.GetGame()
	if gameObj != nil && !gameObj.IsGameOver() {
		if !handler.Session.IsSetupPhase() {
			loser, opponent := getPlayerAndOpponent(gameObj, resigningPlayerIndex, opponentIndex)

			logging.Debug(logging.TagGame, "Declaring Player %s as winner and Player %s as loser via resignation", opponent.GetName(), loser.GetName())
			handler.Session.SetWinner(opponent, game.WinCause("resigned"))

			s.handleGameOver(handler.Session, handler.Hub)
		} else {
			logging.Debug(logging.TagGame, "Session %s resigned during setup phase. Aborting game without win/loss attribution.", gameID)
			handler.Session.Stop("Setup phase aborted")
		}
	}

	s.RemoveSession(gameID, "Player resigned")

	return nil
}

// WinGame declares a winner for a game session by their user ID.
func (s *GameServer) WinGame(gameID string, winnerUserID int, cause string) error {
	handler, exists := s.GetSession(gameID)
	if !exists {
		return fmt.Errorf("game session %s not found", gameID)
	}

	winnerPlayerIndex, loserPlayerIndex, err := resolveUserSeat(handler.Session, winnerUserID)
	if err != nil {
		return err
	}

	gameObj := handler.Session.GetGame()
	if gameObj == nil {
		return fmt.Errorf("game instance not found in session %s", gameID)
	}

	if gameObj.IsGameOver() {
		return fmt.Errorf("game %s is already over", gameID)
	}

	winner, loser := getPlayerAndOpponent(gameObj, winnerPlayerIndex, loserPlayerIndex)

	logging.Debug(logging.TagGame, "Declaring Player %s as winner and Player %s as loser via WinGame", winner.GetName(), loser.GetName())
	handler.Session.SetWinner(winner, game.WinCause(cause))

	s.handleGameOver(handler.Session, handler.Hub)

	return nil
}

// handleGameOver handles the end of a game
func (s *GameServer) handleGameOver(session *game.Session, hub *ws.Hub) {
	state := session.GetGameState()
	logging.Debug(logging.TagGame, "Game %s over. Winner: %v, Cause: %s", session.ID, state.WinnerID, session.GetWinCause())

	ctx := s.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.SaveGameResult(dbCtx, session, hub.GetGameType()); err != nil {
		logging.Error("Failed to save game result", err)
	}

	hub.BroadcastFullState()
	hub.BroadcastMoveHistory()

	hub.StartGameOverCleanup(5 * time.Minute)
}

// monitorGame watches for game events and broadcasts them
func (s *GameServer) monitorGame(handler *SessionHandler) {
	session := handler.Session
	hub := handler.Hub

	time.Sleep(100 * time.Millisecond)
	hub.BroadcastFullState()

	select {
	case <-session.GetSetupCompleteChan():
	case <-session.IsAbortedChan():
		logging.Debug(logging.TagWeb, "Game aborted during setup: %s", session.ID)
		if _, stillExists := s.GetSession(session.ID); stillExists {
			s.RemoveSession(session.ID)
		}
		return
	case <-s.Ctx.Done():
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.Ctx.Done():
			return
		case <-session.IsAbortedChan():
			logging.Debug(logging.TagWeb, "Game aborted during gameplay: %s", session.ID)
			if _, stillExists := s.GetSession(session.ID); stillExists {
				s.RemoveSession(session.ID)
			}
			return
		case <-session.GetMoveNotifyChan():
			isHeadless := session.IsHeadless()
			if !isHeadless {
				logging.Debug(logging.TagGame, "Move executed in game %s", session.ID)
			}

			if isHeadless {
				state := session.GetGameState()
				if state.IsGameOver {
					time.Sleep(100 * time.Millisecond)
					s.handleGameOver(session, hub)
					return
				}
				session.AckMoveProcessed()
				continue
			}

			combat := session.GetLastCombat()
			if combat != nil && combat.Occurred {
				hub.BroadcastCombat(combat)
				session.WaitForAnimationComplete(3 * time.Second)
				session.ClearLastCombat()
				hub.BroadcastFullState()
			} else {
				hub.BroadcastFullState()
			}

			session.AckMoveProcessed()

			state := session.GetGameState()
			if state.IsGameOver {
				time.Sleep(500 * time.Millisecond)
				s.handleGameOver(session, hub)
				return
			}
		case <-ticker.C:
			if !session.IsRunning() && session.GetGameState().IsGameOver {
				s.handleGameOver(session, hub)
				return
			}
		}
	}
}
