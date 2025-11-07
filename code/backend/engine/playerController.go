package engine

import "sync"

type ControllerType int

const (
	HumanController ControllerType = iota
	AIController
)

// PlayerController is the interface that all player controllers must implement
// It allows for both AI and Human players to be used interchangeably
type PlayerController interface {
	GetPlayer() *Player
	GetControllerType() ControllerType
	MakeMove(board *Board) Move
}

// HumanPlayerController represents a human player waiting for input
type HumanPlayerController struct {
	player      *Player
	pendingMove *Move
	mutex       sync.RWMutex
}

func NewHumanPlayerController(player *Player) *HumanPlayerController {
	return &HumanPlayerController{
		player: player,
	}
}

func (h *HumanPlayerController) GetPlayer() *Player {
	return h.player
}

func (h *HumanPlayerController) GetControllerType() ControllerType {
	return HumanController
}

// MakeMove for human returns an empty move - the game should wait for SetPendingMove
func (h *HumanPlayerController) MakeMove(board *Board) Move {
	// Return empty move - game loop should check for this and wait
	return Move{}
}

// SetPendingMove is called by external input (e.g., HTTP handler) to provide the human's move
func (h *HumanPlayerController) SetPendingMove(move Move) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.pendingMove = &move
}

// GetPendingMove retrieves and clears the pending move
func (h *HumanPlayerController) GetPendingMove() *Move {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	move := h.pendingMove
	h.pendingMove = nil
	return move
}

// HasPendingMove checks if a move is waiting
func (h *HumanPlayerController) HasPendingMove() bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.pendingMove != nil
}
