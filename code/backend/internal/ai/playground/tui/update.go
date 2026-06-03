// Package tui coordinates the Bubble Tea fullscreen loop and Huh overlay setups.
package tui

import (
	"digital-innovation/gostrategy/internal/ai/playground"
	"digital-innovation/gostrategy/internal/game"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg struct{}

func tickCmd(ms int) tea.Cmd {
	return tea.Tick(time.Duration(ms)*time.Millisecond, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Update processes terminal input and updates the state machine model losslessly.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case playground.ProgressMsg:
		m.CurrentMatch = msg.Current
		m.TotalMatches = msg.Total
		return m, playground.ListenToSimChan(m.ProgChan)

	case playground.ResultsMsg:
		m.Running = false
		if msg.Error != nil {
			m.ResultsErr = msg.Error
		} else {
			m.Results = msg.Results
			m.ResultsErr = nil
			m.SelectedGame = 0
			m.ActiveTab = tabHistory // Auto-navigate to History on complete
		}
		return m, nil

	case tickMsg:
		if m.ActiveTab == tabReplay && m.ReplayAutoplay {
			gameExport := m.Results.Games[m.SelectedGame]
			totalMoves := len(gameExport.Moves)
			if m.ReplayMoveIndex < totalMoves {
				m.ReplayMoveIndex++
				return m, tickCmd(m.ReplaySpeed)
			}
			m.ReplayAutoplay = false
		}
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	// Delegate input to active modal if one is open
	if m.ActiveComponent != nil {
		activeCmd, done := m.ActiveComponent.Update(msg)
		if done {
			m.ActiveComponent = nil
		}
		return m, activeCmd
	}

	// Normal key handling for layout tabs
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "tab":
			m.navigateTabs(1)
			return m, nil
		case "shift+tab":
			m.navigateTabs(-1)
			return m, nil

		default:
			// Tab-specific shortcuts
			switch m.ActiveTab {
			case tabSetup:
				return m.handleSetupKeys(keyMsg)
			case tabHistory:
				return m.handleHistoryKeys(keyMsg)
			case tabReplay:
				return m.handleReplayKeys(keyMsg)
			}
		}
	}

	return m, nil
}

func (m *Model) navigateTabs(dir int) {
	idx := 0
	for i, t := range m.Tabs {
		if t == m.ActiveTab {
			idx = i
			break
		}
	}
	newIdx := (idx + dir + len(m.Tabs)) % len(m.Tabs)
	m.ActiveTab = m.Tabs[newIdx]
	m.ReplayAutoplay = false
}

func (m *Model) handleSetupKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.Running {
		return m, nil
	}
	if msg.String() == "enter" {
		return m.startSimulation()
	}
	if msg.String() == "c" || msg.String() == "C" {
		m.ActiveComponent = m.BuildSetupForm()
		return m, m.ActiveComponent.Init()
	}
	return m, nil
}

func (m *Model) handleHistoryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if len(m.Results.Games) == 0 {
		return m, nil
	}

	switch msg.String() {
	case "up":
		m.SelectedGame--
		if m.SelectedGame < 0 {
			m.SelectedGame = len(m.Results.Games) - 1
		}
	case "down":
		m.SelectedGame++
		if m.SelectedGame >= len(m.Results.Games) {
			m.SelectedGame = 0
		}
	case "enter":
		m.ReplayMoveIndex = 0
		m.ReplayAutoplay = false
		m.ActiveTab = tabReplay
	case "e", "E":
		_ = playground.ExportSimulation("playground_dataset.json", m.Results)
	}
	return m, nil
}

func (m *Model) handleReplayKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if len(m.Results.Games) == 0 {
		return m, nil
	}
	gameExport := m.Results.Games[m.SelectedGame]
	totalMoves := len(gameExport.Moves)

	switch msg.String() {
	case "right", "n":
		m.ReplayAutoplay = false
		if m.ReplayMoveIndex < totalMoves {
			m.ReplayMoveIndex++
		}
	case "left", "p":
		m.ReplayAutoplay = false
		if m.ReplayMoveIndex > 0 {
			m.ReplayMoveIndex--
		}
	case "space":
		m.ReplayAutoplay = !m.ReplayAutoplay
		if m.ReplayAutoplay && m.ReplayMoveIndex < totalMoves {
			return m, tickCmd(m.ReplaySpeed)
		}
	case "esc":
		m.ReplayAutoplay = false
		m.ActiveTab = tabHistory
	}
	return m, nil
}

func (m *Model) startSimulation() (tea.Model, tea.Cmd) {
	m.Running = true
	m.CurrentMatch = 0
	m.TotalMatches = m.MatchesCount
	m.ProgChan = make(chan tea.Msg, 50)

	go func() {
		playerAlice := game.NewPlayer(0, "Alice AI - "+m.AliceAI, "red")
		playerBob := game.NewPlayer(1, "Bob AI - "+m.BobAI, "blue")

		sim := playground.SimulationExport{
			MatchesCount: m.MatchesCount,
			AliceAI:      m.AliceAI,
			BobAI:        m.BobAI,
			Games:        make([]playground.GameTelemetryExport, m.MatchesCount),
		}

		var setupStr string
		switch m.SetupSelection {
		case "Honey Pot":
			setupStr = "220M22B5BB971877364422342258B56663533B4B"
		case "Random":
			setupStr = "random"
		case "Custom":
			setupStr = m.CustomSetup
		}

		for i := 0; i < m.MatchesCount; i++ {
			m.ProgChan <- playground.ProgressMsg{Current: i + 1, Total: m.MatchesCount}
			gameExport, err := playground.RunSingleMatch(i, m.AliceAI, m.BobAI, m.FatoAggression, setupStr, &playerAlice, &playerBob)
			if err != nil {
				m.ProgChan <- playground.ResultsMsg{Error: err}
				return
			}
			sim.Games[i] = gameExport
			sim.TotalRounds += gameExport.TotalTurns
			switch gameExport.WinnerID {
			case 0:
				sim.AliceWins++
			case 1:
				sim.BobWins++
			default:
				sim.Draws++
			}
		}
		sim.AvgRounds = float64(sim.TotalRounds) / float64(m.MatchesCount)
		m.ProgChan <- playground.ResultsMsg{Results: sim}
	}()

	return m, playground.ListenToSimChan(m.ProgChan)
}
