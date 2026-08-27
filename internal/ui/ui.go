package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/search"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D7D7")).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF87"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700"))

	errorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF5F87")).
			Padding(1).
			MarginTop(1)

	hintBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5FAFFF")).
			Padding(1).
			MarginTop(1)

	chapterHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#BD93F9"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))
)

// FormatBanner renders the application header banner.
func FormatBanner() string {
	return headerStyle.Render("⚡ TERRALINGS: Master Terraform & OpenTofu from Scratch ⚡\n")
}

// FormatResult renders a formatted string describing the outcome of an exercise run.
func FormatResult(res runner.RunResult) string {
	var b strings.Builder
	if res.Passed {
		b.WriteString(successStyle.Render(fmt.Sprintf("✓ Exercise %s passed!\n", res.Exercise.Name)))
	} else {
		if res.HasNotDoneMarker {
			b.WriteString(warningStyle.Render(fmt.Sprintf("⌛ %s still contains '%s' marker. Keep going!\n", res.Exercise.Name, runner.NotDoneMarker)))
		}
		if res.Error != "" {
			b.WriteString(errorBoxStyle.Render(fmt.Sprintf("Error in %s:\n%s", res.Exercise.Name, res.Error)))
			b.WriteString("\n")
		} else if res.Output != "" {
			b.WriteString(fmt.Sprintf("\n%s\n", res.Output))
		}
	}
	return b.String()
}

// FormatHint renders a hint box for an exercise at the given index.
func FormatHint(ex *models.Exercise, hintIdx int) string {
	if ex == nil || len(ex.Hints) == 0 {
		return warningStyle.Render("No hints available for this exercise.")
	}
	idx := hintIdx
	if idx < 0 {
		idx = 0
	} else if idx >= len(ex.Hints) {
		idx = len(ex.Hints) - 1
	}
	return hintBoxStyle.Render(fmt.Sprintf("💡 Hint (%d/%d) for %s:\n%s", idx+1, len(ex.Hints), ex.Name, ex.Hints[idx]))
}

// FormatProgress renders a progress bar and completion statistics.
func FormatProgress(completed, total int) string {
	if total <= 0 {
		return "Progress: [--------------------] 0/0 (0%)"
	}
	if completed < 0 {
		completed = 0
	}
	if completed > total {
		completed = total
	}

	percent := (completed * 100) / total
	barWidth := 20
	filled := (completed * barWidth) / total
	empty := barWidth - filled

	bar := "[" + strings.Repeat("=", filled) + strings.Repeat("-", empty) + "]"
	return fmt.Sprintf("Progress: %s %d/%d (%d%%)", bar, completed, total, percent)
}

// FormatChapterList renders the curriculum table with chapter information and exercise completion status.
func FormatChapterList(m *models.Manifest, statuses map[string]models.ExerciseStatus) string {
	if m == nil {
		return warningStyle.Render("No curriculum manifest loaded.")
	}

	var b strings.Builder
	for _, ch := range m.Chapters {
		b.WriteString(chapterHeaderStyle.Render(fmt.Sprintf("Chapter %02d: %s", ch.Number, ch.Title)))
		if ch.Description != "" {
			b.WriteString(dimStyle.Render(fmt.Sprintf(" - %s", ch.Description)))
		}
		b.WriteString("\n")

		for _, ex := range ch.Exercises {
			status := models.StatusNotStarted
			if statuses != nil {
				if s, ok := statuses[ex.Name]; ok {
					status = s
				}
			}

			var statusIcon string
			switch status {
			case models.StatusCompleted:
				statusIcon = successStyle.Render("✓")
			case models.StatusInProgress:
				statusIcon = warningStyle.Render("•")
			case models.StatusFailed:
				statusIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87")).Render("✕")
			default:
				statusIcon = dimStyle.Render("·")
			}

			b.WriteString(fmt.Sprintf("  %s %-16s : %s (%s)\n", statusIcon, ex.Name, ex.Title, ex.Path))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// FormatSearchResults renders exercise search results in a clean styled list.
func FormatSearchResults(query string, results []search.SearchResult) string {
	if len(results) == 0 {
		return warningStyle.Render(fmt.Sprintf("No exercises found matching '%s'.\n", query))
	}

	var b strings.Builder
	b.WriteString(headerStyle.Render(fmt.Sprintf("🔍 Search Results for '%s' (%d matches):\n", query, len(results))))

	for _, r := range results {
		b.WriteString(fmt.Sprintf("  • %-16s %s\n", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00D7D7")).Render(r.Exercise.Name), r.Exercise.Title))
		b.WriteString(fmt.Sprintf("    %s | matched in: %s\n", dimStyle.Render(r.Exercise.Path), lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Render(r.MatchedIn)))
	}
	return b.String()
}
