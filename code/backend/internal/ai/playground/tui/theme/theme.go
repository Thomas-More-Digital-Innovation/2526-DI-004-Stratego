// Package theme defines the visual color palette and token definitions for the TUI.
package theme

import "github.com/charmbracelet/lipgloss"

// Theme bundles the color configuration for consistent component styling.
type Theme struct {
	Primary   lipgloss.TerminalColor
	Secondary lipgloss.TerminalColor
	Muted     lipgloss.TerminalColor
	Bg        lipgloss.TerminalColor
	Success   lipgloss.TerminalColor
	Error     lipgloss.TerminalColor
}

// Global bundles the system theme colors for UI components.
var Global = Theme{
	Primary:   lipgloss.Color("#38BDF8"),
	Secondary: lipgloss.Color("#1F1F1F"),
	Muted:     lipgloss.Color("#6B7280"),
	Bg:        lipgloss.Color("#111827"),
	Success:   lipgloss.Color("#10B981"),
	Error:     lipgloss.Color("#EF4444"),
}
