// Package tui coordinates the Bubble Tea fullscreen loop and Huh overlay setups.
package tui

import (
	"digital-innovation/gostrategy/internal/ai/playground/tui/style"
	"digital-innovation/gostrategy/internal/game/models"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) renderReplayBody(s *strings.Builder) {
	gameExport := m.Results.Games[m.SelectedGame]
	board := m.getReplayBoardState()
	totalMoves := len(gameExport.Moves)

	fmt.Fprintf(s, "Replay: Match #%d | Total Turns: %d | Move: %d/%d | Autoplay: %t\n",
		gameExport.GameIndex+1, gameExport.TotalTurns, m.ReplayMoveIndex, totalMoves, m.ReplayAutoplay,
	)

	if m.ReplayMoveIndex > 0 && m.ReplayMoveIndex-1 < totalMoves {
		lastMove := gameExport.Moves[m.ReplayMoveIndex-1]
		player := ternary(lastMove.PlayerID == 0, "Alice", "Bob")
		pieceName := "Piece"
		if lastMove.Attacker != nil {
			pieceName = lastMove.Attacker.Type
		}
		moveStr := fmt.Sprintf("Last Move: %s (%s) moved (%d,%d) -> (%d,%d)",
			player, pieceName, lastMove.FromX, lastMove.FromY, lastMove.ToX, lastMove.ToY,
		)
		if lastMove.Result != models.ResultMove {
			combatStr := fmt.Sprintf(" COMBAT RESULT: %s! (%s vs %s)",
				lastMove.Result, lastMove.Attacker.Type,
				ternary(lastMove.Defender != nil, lastMove.Defender.Type, "?"),
			)
			s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B")).Bold(true).Render(moveStr+combatStr) + "\n\n")
		} else {
			s.WriteString(style.RowInactive.Render(moveStr) + "\n\n")
		}
	} else {
		s.WriteString(style.Muted.Render("Initial placement layout setup.") + "\n\n")
	}

	// Render board
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Bold(true)
	blueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6")).Bold(true)
	lakeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#06B6D4")).Background(lipgloss.Color("#164E63"))
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#4B5563"))

	s.WriteString(borderStyle.Render("  +---------------------+") + "\n")
	for y := 0; y < 10; y++ {
		s.WriteString(borderStyle.Render("  | "))
		for x := 0; x < 10; x++ {
			cell := board[y][x]
			isLake := (y == 4 || y == 5) && (x == 2 || x == 3 || x == 6 || x == 7)

			switch {
			case cell.OwnerID == 0:
				s.WriteString(redStyle.Render(getRankChar(cell.Type)) + " ")
			case cell.OwnerID == 1:
				s.WriteString(blueStyle.Render(getRankChar(cell.Type)) + " ")
			case isLake:
				s.WriteString(lakeStyle.Render("~") + " ")
			default:
				s.WriteString(style.Muted.Render(".") + " ")
			}
		}
		s.WriteString(borderStyle.Render("|") + "\n")
	}
	s.WriteString(borderStyle.Render("  +---------------------+") + "\n")
}

func (m *Model) getReplayBoardState() [][]models.PieceData {
	gameExport := m.Results.Games[m.SelectedGame]
	board := make([][]models.PieceData, 10)
	for y := 0; y < 10; y++ {
		board[y] = make([]models.PieceData, 10)
		copy(board[y], gameExport.InitialBoard[y])
	}
	for i := 0; i < m.ReplayMoveIndex; i++ {
		if i >= len(gameExport.Moves) {
			break
		}
		move := gameExport.Moves[i]
		attacker := board[move.FromY][move.FromX]
		board[move.FromY][move.FromX] = models.PieceData{OwnerID: -1}

		switch move.Result {
		case models.ResultMove, models.ResultWin, models.ResultCapture:
			if move.Attacker != nil {
				attacker = *move.Attacker
			}
			board[move.ToY][move.ToX] = attacker
		case models.ResultLoss:
			if move.Defender != nil {
				board[move.ToY][move.ToX] = *move.Defender
			}
		case models.ResultTie:
			board[move.ToY][move.ToX] = models.PieceData{OwnerID: -1}
		}
	}
	return board
}

func getRankChar(typeName string) string {
	switch typeName {
	case "Flag":
		return "0"
	case "Bomb":
		return "B"
	case "Spy":
		return "1"
	case "Scout":
		return "2"
	case "Miner":
		return "3"
	case "Sergeant":
		return "4"
	case "Lieutenant":
		return "5"
	case "Captain":
		return "6"
	case "Major":
		return "7"
	case "Colonel":
		return "8"
	case "General":
		return "9"
	case "Marshal":
		return "M"
	default:
		return "?"
	}
}

func ternary(cond bool, t, f string) string {
	if cond {
		return t
	}
	return f
}
