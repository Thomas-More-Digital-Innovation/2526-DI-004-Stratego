// Package tui coordinates the Bubble Tea fullscreen loop and Huh overlay setups.
package tui

import (
	"digital-innovation/gostrategy/internal/ai/playground/tui/style"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

// View renders the AI Playground TUI interface based on state and window constraints.
func (m *Model) View() string {
	if m.Width < 40 || m.Height < 20 {
		return "Terminal is too small."
	}

	bgString := m.renderNormalView(m.Width, m.Height)

	if m.ActiveComponent != nil {
		dialogContent := m.ActiveComponent.View(54)
		bgLines := strings.Split(stripANSI(bgString), "\n")

		dialogWidth := min(62, m.Width-4)
		dialogHeight := min(18, m.Height-2)

		dialogBox := style.Dialog.Width(dialogWidth).Height(dialogHeight).Render(dialogContent)
		dialogLines := strings.Split(dialogBox, "\n")

		dialogW := lipgloss.Width(dialogBox)
		dialogH := len(dialogLines)

		startRow := (len(bgLines) - dialogH) / 2
		startCol := (m.Width - dialogW) / 2

		var finalLines []string
		for i, bgLine := range bgLines {
			bgRunes := []rune(bgLine)
			if len(bgRunes) < m.Width {
				bgRunes = append(bgRunes, []rune(strings.Repeat(" ", m.Width-len(bgRunes)))...)
			} else if len(bgRunes) > m.Width {
				bgRunes = bgRunes[:m.Width]
			}

			if i >= startRow && i < startRow+dialogH {
				dialogLineIdx := i - startRow
				leftPart := bgRunes[:startCol]
				rightPart := bgRunes[startCol+dialogW:]

				mutedLeft := style.Muted.Render(string(leftPart))
				mutedRight := style.Muted.Render(string(rightPart))
				dialogLine := dialogLines[dialogLineIdx]

				finalLines = append(finalLines, mutedLeft+dialogLine+mutedRight)
			} else {
				finalLines = append(finalLines, style.Muted.Render(string(bgRunes)))
			}
		}
		bgString = strings.Join(finalLines, "\n")
	}

	return bgString
}

func (m *Model) renderNormalView(width, height int) string {
	headerBoxHeight := 3
	footerBoxHeight := 3
	bodyBoxHeight := max(height-headerBoxHeight-footerBoxHeight, 2)

	headerContent := m.renderHeader()
	headerBox := style.HeaderBox.
		Width(width - 2).
		Height(headerBoxHeight - 2).
		Render(headerContent)

	bodyContent := m.renderBody(width-4, bodyBoxHeight-2)
	bodyBox := style.BodyBox.
		Width(width - 2).
		Height(bodyBoxHeight - 2).
		Render(bodyContent)

	footerContent := m.renderFooter()
	footerBox := style.FooterBox.
		Width(width - 2).
		Height(footerBoxHeight - 2).
		Render(footerContent)

	return lipgloss.JoinVertical(lipgloss.Left,
		headerBox,
		bodyBox,
		footerBox,
	)
}

func (m *Model) renderHeader() string {
	var tabs []string
	for _, t := range m.Tabs {
		if t == m.ActiveTab {
			tabs = append(tabs, style.TabActive.Render(t))
		} else {
			tabs = append(tabs, style.TabInactive.Render(t))
		}
	}
	tabsRow := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)
	title := style.Title.Render(" GoStrategy AI Playground ")
	return lipgloss.JoinHorizontal(lipgloss.Left, title, "── ", tabsRow)
}

func (m *Model) renderBody(width, _ int) string {
	var s strings.Builder

	switch m.ActiveTab {
	case tabSetup:
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8")).Bold(true).Render("Simulation Settings:") + "\n\n")
		fmt.Fprintf(&s, "  Alice Strategy (P1): %s\n", m.AliceAI)
		fmt.Fprintf(&s, "  Bob Strategy (P2):   %s\n", m.BobAI)
		fmt.Fprintf(&s, "  FATO Aggression:     %.2f\n", m.FatoAggression)
		fmt.Fprintf(&s, "  Total Matches:       %d\n", m.MatchesCount)
		fmt.Fprintf(&s, "  Board Setup Layout:  %s\n\n", m.SetupSelection)

		if m.Running {
			s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B")).Bold(true).Render("Simulating matches...") + "\n")
			percent := float64(m.CurrentMatch) / float64(m.TotalMatches)
			widthBar := min(width-12, 40)
			filled := int(percent * float64(widthBar))
			if filled < 0 {
				filled = 0
			}
			empty := widthBar - filled
			if empty < 0 {
				empty = 0
			}
			bar := lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")).Render(strings.Repeat("█", filled)) +
				lipgloss.NewStyle().Foreground(lipgloss.Color("#374151")).Render(strings.Repeat("░", empty))
			fmt.Fprintf(&s, "  [%s] %.1f%% (%d/%d matches)\n", bar, percent*100.0, m.CurrentMatch, m.TotalMatches)
		} else {
			s.WriteString(style.Muted.Render("Press [C] to open configuration modal and edit parameters.") + "\n")
		}

	case tabHistory:
		if len(m.Results.Games) == 0 {
			s.WriteString(style.Muted.Render("No simulation run records. Head to the Setup tab and launch one!") + "\n")
			break
		}
		total := float64(m.Results.MatchesCount)
		fmt.Fprintf(&s, "Tourney Win Ratio: Alice %d (%.1f%%) | Bob %d (%.1f%%) | Draws %d (%.1f%%) | Avg Turns: %.1f\n\n",
			m.Results.AliceWins, float64(m.Results.AliceWins)/total*100.0,
			m.Results.BobWins, float64(m.Results.BobWins)/total*100.0,
			m.Results.Draws, float64(m.Results.Draws)/total*100.0,
			m.Results.AvgRounds,
		)

		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#38BDF8")).Bold(true).Render("Simulation History (Press Enter to Replay):") + "\n")
		start := m.SelectedGame - 2
		if start < 0 {
			start = 0
		}
		end := start + 6
		if end > len(m.Results.Games) {
			end = len(m.Results.Games)
			start = end - 6
			if start < 0 {
				start = 0
			}
		}

		for i := start; i < end; i++ {
			g := m.Results.Games[i]
			var winnerStr string
			switch g.WinnerID {
			case 0:
				winnerStr = "Alice"
			case 1:
				winnerStr = "Bob"
			default:
				winnerStr = "Draw"
			}
			row := fmt.Sprintf("  %s Match #%d: Winner: %s (Turns: %d, Cause: %s)",
				ternary(i == m.SelectedGame, "●", " "),
				g.GameIndex+1, winnerStr, g.TotalTurns, g.WinCause,
			)
			if i == m.SelectedGame {
				s.WriteString(style.RowActive.Render(row) + "\n")
			} else {
				s.WriteString(style.RowInactive.Render(row) + "\n")
			}
		}

	case tabReplay:
		if len(m.Results.Games) == 0 {
			s.WriteString(style.Muted.Render("No active match loaded to replay.") + "\n")
			break
		}
		m.renderReplayBody(&s)
	}

	return s.String()
}

func (m *Model) renderFooter() string {
	switch m.ActiveTab {
	case tabSetup:
		if m.Running {
			return style.Footer.Render("Simulating matches in progress... Please wait.")
		}
		return style.Footer.Render("[Tab/Shift-Tab] Switch Tabs | [C] Configure Simulation | [Enter] Start Simulation | [Ctrl+C] Quit")
	case tabHistory:
		return style.Footer.Render("[Tab/Shift-Tab] Switch Tabs | [Up/Down] Select Match | [Enter] Replay | [E] Export JSON | [Ctrl+C] Quit")
	case tabReplay:
		return style.Footer.Render("[Left/Right] Step Move | [Space] Play/Pause | [Esc] Return to History | [Ctrl+C] Quit")
	default:
		return ""
	}
}

