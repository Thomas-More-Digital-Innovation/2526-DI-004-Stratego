// Package tui coordinates the Bubble Tea fullscreen loop and Huh overlay setups.
package tui

import (
	"fmt"

	"digital-innovation/gostrategy/internal/ai/playground/tui/components"
	"digital-innovation/gostrategy/internal/game/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// BuildSetupForm creates a multi-step interactive form using huh for configuring the simulation.
func (m *Model) BuildSetupForm() *components.Form {
	var groups []*huh.Group

	// Group 1: Player selection and parameters
	groups = append(groups, huh.NewGroup(
		huh.NewSelect[string]().
			Title("Alice AI (Player 1)").
			Options(
				huh.NewOption("FAFO (Random)", models.Fafo),
				huh.NewOption("FATO (Aggressive)", models.Fato),
			).
			Value(&m.AliceAI),
		huh.NewSelect[string]().
			Title("Bob AI (Player 2)").
			Options(
				huh.NewOption("FAFO (Random)", models.Fafo),
				huh.NewOption("FATO (Aggressive)", models.Fato),
			).
			Value(&m.BobAI),
		huh.NewSelect[float64]().
			Title("FATO Aggression Level").
			Options(
				huh.NewOption("Low (0.2)", 0.2),
				huh.NewOption("Medium (0.5)", 0.5),
				huh.NewOption("High (0.8)", 0.8),
				huh.NewOption("Maximum (1.0)", 1.0),
			).
			Value(&m.FatoAggression),
		huh.NewSelect[int]().
			Title("Simulation Round Matches").
			Options(
				huh.NewOption("10 matches", 10),
				huh.NewOption("100 matches", 100),
				huh.NewOption("500 matches", 500),
				huh.NewOption("1000 matches", 1000),
			).
			Value(&m.MatchesCount),
	))

	// Group 2: Setup layout options
	groups = append(groups, huh.NewGroup(
		huh.NewSelect[string]().
			Title("Board Placement Layout").
			Options(
				huh.NewOption("Honey Pot Preset", "Honey Pot"),
				huh.NewOption("Random Setup Placement", "Random"),
				huh.NewOption("Custom 40-character Layout", "Custom"),
			).
			Value(&m.SetupSelection),
	))

	// Group 3: Custom Layout input (always append; if they skip it, default is used)
	groups = append(groups, huh.NewGroup(
		huh.NewInput().
			Title("Custom Layout Setup").
			Description("40 characters representing piece ranks (e.g. 0=Flag, B=Bomb, 1-9, M=Marshal)").
			Value(&m.CustomSetup).
			Validate(func(s string) error {
				if len(s) != 40 {
					return fmt.Errorf("layout must be exactly 40 characters (currently %d)", len(s))
				}
				return nil
			}),
	))

	form := huh.NewForm(groups...).
		WithTheme(huh.ThemeCharm()).
		WithWidth(60)

	return &components.Form{
		Form: form,
		OnSubmit: func() tea.Cmd {
			_, cmd := m.startSimulation()
			return cmd
		},
	}
}
