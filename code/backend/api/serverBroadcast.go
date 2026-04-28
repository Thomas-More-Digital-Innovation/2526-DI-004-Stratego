package api

import (
	"digital-innovation/stratego/game"
	"digital-innovation/stratego/models"
	"log"
)

// broadcastFullState sends complete game state and board to all clients
func (s *GameServer) broadcastFullState(hub *WSHub, gameType string) {
	state := hub.session.GetGameState()

	// Broadcast game state
	hub.BroadcastGameState()

	// Broadcast board state
	switch {
	case state.IsSetupPhase:
		s.broadcastSetupBoard(hub, gameType)
	case gameType == models.AiVsAi:
		s.broadcastBoardStateRevealed(hub)
	default:
		s.broadcastBoardStatePerClient(hub)
	}
}

// broadcastBoardState sends board state to all clients
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func (s *GameServer) broadcastBoardState(hub *WSHub, viewerID int) {
	board := hub.session.GetBoard()
	field := board.GetField()

	boardDTO := make([][]PieceDTO, 10)
	for y := 0; y < 10; y++ {
		boardDTO[y] = make([]PieceDTO, 10)
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			if piece != nil && piece.IsAlive() {
				dto := PieceToDTO(piece, viewerID)
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

	hub.BroadcastMessage(MsgTypeBoardState, boardMsg)
}

// broadcastSetupBoard sends the setup board state (pieces not yet placed on board)
func (s *GameServer) broadcastSetupBoard(hub *WSHub, gameType string) {
	hub.BroadcastSetupBoard()
}

// broadcastBoardStatePerClient sends personalized board state to each connected client
func (s *GameServer) broadcastBoardStatePerClient(hub *WSHub) {
	hub.mutex.RLock()
	clients := make([]*WSClient, 0, len(hub.clients))
	for client := range hub.clients {
		clients = append(clients, client)
	}
	hub.mutex.RUnlock()

	// Send personalized board state to each client
	for _, client := range clients {
		hub.sendBoardState(client)
	}
}

// broadcastCombat sends combat information to all clients
func (s *GameServer) broadcastCombat(hub *WSHub, combat *game.CombatResult, gameType string) {
	if combat == nil || !combat.Occurred {
		return
	}

	attacker := combat.AttackerPiece
	defender := combat.DefenderPiece

	// For AI vs AI or spectators, reveal both pieces
	// For player games, reveal based on ownership
	attackerDTO := PieceToDTO(attacker, attacker.GetOwner().GetID())
	attackerDTO.Position = PositionToDTO(combat.AttackerPosition)
	attackerDTO.Revealed = true

	defenderDTO := PieceToDTO(defender, defender.GetOwner().GetID())
	defenderDTO.Position = PositionToDTO(combat.DefenderPosition)
	defenderDTO.Revealed = true

	combatMsg := CombatMessage{
		Attacker:     attackerDTO,
		Defender:     defenderDTO,
		AttackerWon:  attacker.IsAlive(),
		DefenderWon:  defender.IsAlive(),
		AttackerDied: !attacker.IsAlive(),
		DefenderDied: !defender.IsAlive(),
	}

	hub.BroadcastMessage(MsgTypeCombat, combatMsg)
	log.Printf("Combat message sent: %+v", combatMsg)
}

// broadcastBoardStateRevealed sends board state with all pieces revealed (for AI vs AI spectating)
func (s *GameServer) broadcastBoardStateRevealed(hub *WSHub) {
	board := hub.session.GetBoard()
	field := board.GetField()

	boardDTO := make([][]PieceDTO, 10)
	for y := 0; y < 10; y++ {
		boardDTO[y] = make([]PieceDTO, 10)
		for x := 0; x < 10; x++ {
			piece := field[y][x]
			if piece != nil && piece.IsAlive() {
				// Force reveal all pieces for spectators
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

	hub.BroadcastMessage(MsgTypeBoardState, boardMsg)
}
