package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/state"
)

var (
	headerTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00D7D7"))

	headerProgressStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#50FA7B"))

	headerDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))

	activeBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00D7D7")).
			Padding(0, 1)

	inactiveBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#44475A")).
				Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#00D7D7"))

	passedItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B"))

	inProgressItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFB86C"))

	failedItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555"))

	dimItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))

	chapterTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#BD93F9"))

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D7D7"))

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#44475A"))

	hintBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5FAFFF")).
			Padding(0, 1).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF5555")).
			Padding(0, 1).
			Foreground(lipgloss.Color("#FF5555"))

	footerHelpKey = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D7D7"))

	footerHelpDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))
)

// View renders the complete Bubble Tea user interface.
func (m Model) View() string {
	w := m.width
	if w < 60 {
		w = 60
	}
	h := m.height
	if h < 20 {
		h = 20
	}

	header := m.renderHeader(w)
	footer := m.renderFooter(w)

	sidebarW := w / 3
	if sidebarW < 28 {
		sidebarW = 28
	}
	if sidebarW > 42 {
		sidebarW = 42
	}

	contentW := w - sidebarW - 4
	if contentW < 24 {
		contentW = 24
	}

	bodyH := h - 6
	if bodyH < 10 {
		bodyH = 10
	}

	sidebar := m.renderSidebar(sidebarW, bodyH)
	mainContent := m.renderMain(contentW, bodyH)

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainContent)
	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m Model) renderHeader(width int) string {
	title := headerTitleStyle.Render("⚡ TERRALINGS v0.2.0")

	// Calculate overall progress
	total := len(m.exercises)
	completed := 0
	if m.store != nil {
		allStates := m.store.GetAllExerciseStates()
		for _, st := range allStates {
			if st.Status == state.StatusPassed {
				completed++
			}
		}
	}

	percent := 0
	if total > 0 {
		percent = (completed * 100) / total
	}

	barWidth := 16
	filled := 0
	if total > 0 {
		filled = (completed * barWidth) / total
	}
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled
	bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"

	progText := fmt.Sprintf("%s %d%% (%d/%d)", bar, percent, completed, total)
	progRendered := headerProgressStyle.Render(progText)

	// Distribute space
	gap := width - lipgloss.Width(title) - lipgloss.Width(progRendered) - 2
	if gap < 1 {
		gap = 1
	}

	return title + strings.Repeat(" ", gap) + progRendered
}

func (m Model) renderSidebar(width, height int) string {
	var lines []string
	cursorLine := 0

	if m.searching {
		lines = append(lines, chapterTitleStyle.Render(fmt.Sprintf("Search: %q (%d)", m.searchInput.Value(), len(m.filteredIndices))))
		for i, idx := range m.filteredIndices {
			ex := m.exercises[idx]
			st := state.StatusNotStarted
			if m.store != nil {
				if s := m.store.GetExerciseState(ex.Name); s != nil {
					st = s.Status
				}
			}

			icon := dimItemStyle.Render("·")
			if st == state.StatusPassed {
				icon = passedItemStyle.Render("✓")
			} else if st == state.StatusInProgress {
				icon = inProgressItemStyle.Render("•")
			}

			line := fmt.Sprintf(" %s %s", icon, ex.Name)
			if i == m.cursor {
				cursorLine = len(lines)
				line = selectedItemStyle.Render(fmt.Sprintf(" ▶ %s", ex.Name))
			}
			lines = append(lines, line)
		}
	} else {
		// Group by chapter
		currentCh := -1
		exItemIdx := 0

		for _, ex := range m.exercises {
			// Find chapter info
			chNum := 0
			chTitle := ex.ChapterName
			for _, ch := range m.manifest.Chapters {
				if ch.Name == ex.ChapterName {
					chNum = ch.Number
					chTitle = ch.Title
					break
				}
			}

			if chNum != currentCh {
				currentCh = chNum
				lines = append(lines, chapterTitleStyle.Render(fmt.Sprintf("▼ Ch %02d: %s", chNum, chTitle)))
			}

			st := state.StatusNotStarted
			if m.store != nil {
				if s := m.store.GetExerciseState(ex.Name); s != nil {
					st = s.Status
				}
			}

			icon := dimItemStyle.Render("·")
			if st == state.StatusPassed {
				icon = passedItemStyle.Render("✓")
			} else if st == state.StatusInProgress {
				icon = inProgressItemStyle.Render("•")
			}

			line := fmt.Sprintf("   %s %s", icon, ex.Name)
			if exItemIdx == m.cursor {
				cursorLine = len(lines)
				line = selectedItemStyle.Render(fmt.Sprintf(" ▶ %s", ex.Name))
			}
			lines = append(lines, line)
			exItemIdx++
		}
	}

	// Window lines to fit inside height
	maxLines := height - 2
	if maxLines < 1 {
		maxLines = 1
	}

	start := 0
	if cursorLine > maxLines/2 {
		start = cursorLine - maxLines/2
	}
	if start+maxLines > len(lines) {
		start = len(lines) - maxLines
	}
	if start < 0 {
		start = 0
	}

	end := start + maxLines
	if end > len(lines) {
		end = len(lines)
	}

	visibleLines := lines[start:end]
	content := strings.Join(visibleLines, "\n")

	boxStyle := inactiveBoxStyle
	if m.activePane == PaneSidebar {
		boxStyle = activeBoxStyle
	}

	return boxStyle.Width(width).Height(height).Render(content)
}

func (m Model) renderMain(width, height int) string {
	ex := m.SelectedExercise()
	var b strings.Builder

	// Header
	b.WriteString(labelStyle.Render("EXERCISE: ") + headerTitleStyle.Render(ex.Name) + "\n")
	if ex.Title != "" {
		b.WriteString(headerDimStyle.Render("Title: ") + ex.Title + "\n")
	}
	b.WriteString(headerDimStyle.Render("File:  ") + ex.Path + "\n")
	b.WriteString(headerDimStyle.Render("Mode:  ") + string(ex.Mode) + "\n\n")

	// Divider
	divW := width - 4
	if divW < 10 {
		divW = 10
	}
	b.WriteString(dividerStyle.Render(strings.Repeat("─", divW)) + "\n")

	// Compiler output
	if m.running {
		b.WriteString(fmt.Sprintf("\n %s Evaluating %s...\n", m.spinner.View(), ex.Name))
	} else if m.lastResult != nil {
		res := m.lastResult
		if res.Passed {
			b.WriteString(passedItemStyle.Render(fmt.Sprintf("\n✓ Exercise %s passed! All checks succeeded.\n", res.Exercise.Name)))
		} else {
			if res.HasNotDoneMarker {
				b.WriteString(inProgressItemStyle.Render(fmt.Sprintf("\n⌛ %s still contains '%s' marker. Keep going!\n", res.Exercise.Name, runner.NotDoneMarker)))
			}
			if res.Error != "" {
				b.WriteString(errorStyle.Render(fmt.Sprintf("Error in %s:\n%s", res.Exercise.Name, res.Error)))
				b.WriteString("\n")
			} else if res.Output != "" {
				b.WriteString(fmt.Sprintf("\n%s\n", res.Output))
			}
		}
	} else {
		b.WriteString(headerDimStyle.Render("\nPress [Enter] to run validation, or edit the exercise file in your editor.\n"))
	}

	// Hint Drawer
	if m.showHints && len(ex.Hints) > 0 {
		idx := m.hintIndex
		if idx < 0 {
			idx = 0
		} else if idx >= len(ex.Hints) {
			idx = len(ex.Hints) - 1
		}
		b.WriteString("\n" + dividerStyle.Render(strings.Repeat("─", divW)) + "\n")
		b.WriteString(hintBoxStyle.Render(fmt.Sprintf("💡 Hint (%d/%d) for %s:\n%s", idx+1, len(ex.Hints), ex.Name, ex.Hints[idx])))
	}

	boxStyle := inactiveBoxStyle
	if m.activePane == PaneViewport {
		boxStyle = activeBoxStyle
	}

	return boxStyle.Width(width).Height(height).Render(b.String())
}

func (m Model) renderFooter(width int) string {
	if m.searching {
		return fmt.Sprintf("Search: %s  %s %s",
			m.searchInput.View(),
			footerHelpKey.Render("[Enter]"), footerHelpDesc.Render("select"),
		)
	}

	if m.statusMessage != "" {
		return inProgressItemStyle.Render(m.statusMessage)
	}

	return fmt.Sprintf("%s %s  %s %s  %s %s  %s %s  %s %s  %s %s  %s %s",
		footerHelpKey.Render("[↑/↓/j/k]"), footerHelpDesc.Render("Navigate"),
		footerHelpKey.Render("[Tab]"), footerHelpDesc.Render("Switch Pane"),
		footerHelpKey.Render("[Enter]"), footerHelpDesc.Render("Run"),
		footerHelpKey.Render("[h]"), footerHelpDesc.Render("Hint"),
		footerHelpKey.Render("[r]"), footerHelpDesc.Render("Reset"),
		footerHelpKey.Render("[/]"), footerHelpDesc.Render("Search"),
		footerHelpKey.Render("[q]"), footerHelpDesc.Render("Quit"),
	)
}
