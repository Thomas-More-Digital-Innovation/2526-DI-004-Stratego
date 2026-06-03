// Package tui coordinates the Bubble Tea fullscreen loop and Huh overlay setups.
package tui

import (
	"io"
	"log"
	"os"

	"digital-innovation/gostrategy/internal/ai/playground"
	"digital-innovation/gostrategy/internal/ai/playground/tui/components"
	"digital-innovation/gostrategy/internal/game/models"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	tabSetup   = "Setup"
	tabHistory = "History"
	tabReplay  = "Replay"
)

// Model holds the state machine parameters for the Bubble Tea application loop.
type Model struct {
	// Dimensions
	Width  int
	Height int

	// Active tab
	ActiveTab string
	Tabs      []string

	// Active form modal overlay
	ActiveComponent components.Component

	// Simulation parameters
	AliceAI        string
	BobAI          string
	FatoAggression float64
	MatchesCount   int
	SetupSelection string // "Honey Pot", "Random", "Custom"
	CustomSetup    string // 40-character board setup layout string

	// Simulation execution state
	Running      bool
	ProgChan     chan tea.Msg
	CurrentMatch int
	TotalMatches int
	Results      playground.SimulationExport
	ResultsErr   error

	// History navigation
	SelectedGame int

	// Replay state
	ReplayMoveIndex int
	ReplayAutoplay  bool
	ReplaySpeed     int // ms delay for autoplay
}

// NewModel initializes and configures a new TUI application state model.
func NewModel() *Model {
	// Completely disable all logs from stdout/stderr to avoid terminal corruption
	log.SetOutput(io.Discard)
	_ = os.Setenv("GOSTRATEGY_DEBUG", "false")

	return &Model{
		ActiveTab:      tabSetup,
		Tabs:           []string{tabSetup, tabHistory, tabReplay},
		AliceAI:        models.Fafo,
		BobAI:          models.Fato,
		FatoAggression: 0.5,
		MatchesCount:   100,
		SetupSelection: "Honey Pot",
		CustomSetup:    "220M22B5BB971877364422342258B56663533B4B",
		ReplaySpeed:    300,
	}
}

// Init initializes the Bubble Tea application state and enters alternative screen.
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("GoStrategy AI Playground"),
		tea.EnterAltScreen,
	)
}
