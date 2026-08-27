package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/search"
	"github.com/dnf0/terralings/internal/state"
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

// FormatSuccess renders a celebratory success message.
func FormatSuccess(msg string) string {
	return successStyle.Render(msg) + "\n"
}

// FormatInteractivePrompt renders the interactive watcher key commands.
func FormatInteractivePrompt() string {
	promptStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#BD93F9"))
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D7D7"))
	dim := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4"))

	return fmt.Sprintf("\n%s %s %s %s %s %s %s %s %s %s %s\n",
		keyStyle.Render("[Enter / n]"), promptStyle.Render("Next exercise"),
		dim.Render("|"),
		keyStyle.Render("[p]"), promptStyle.Render("Previous"),
		dim.Render("|"),
		keyStyle.Render("[r]"), promptStyle.Render("Rerun"),
		dim.Render("|"),
		keyStyle.Render("[q]"), promptStyle.Render("Quit"),
	)
}

// FormatResult renders a formatted string describing the outcome of an exercise run.
func FormatResult(res runner.RunResult) string {
	var b strings.Builder
	if res.Passed {
		b.WriteString(successStyle.Render(fmt.Sprintf("✓ Exercise %s passed!\n", res.Exercise.Name)))
	} else {
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

// FormatAnalytics renders comprehensive learning and progress analytics.
func FormatAnalytics(summary state.AnalyticsSummary) string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D7D7"))
	b.WriteString(titleStyle.Render("📊 TERRALINGS LEARNING ANALYTICS") + "\n\n")

	// Overall Progress bar
	barWidth := 20
	var filled int
	var percent int
	if summary.TotalExercises > 0 {
		percent = (summary.CompletedCount * 100) / summary.TotalExercises
		filled = (summary.CompletedCount * barWidth) / summary.TotalExercises
	}
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled
	bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
	b.WriteString(fmt.Sprintf("Overall Progress: %s %d%% (%d/%d completed)\n", bar, percent, summary.CompletedCount, summary.TotalExercises))

	// Total Time Invested
	timeStr := formatDuration(summary.TotalTimeSpent)
	b.WriteString(fmt.Sprintf("Time Invested:    %s\n", timeStr))

	// Attempts & average attempts
	avgAttempts := 0.0
	if summary.TotalExercises > 0 {
		avgAttempts = float64(summary.TotalAttempts) / float64(summary.TotalExercises)
	}
	b.WriteString(fmt.Sprintf("Total Attempts:   %d (avg %.1f per exercise)\n", summary.TotalAttempts, avgAttempts))

	// Hints consulted
	b.WriteString(fmt.Sprintf("Hints Consulted:  %d\n\n", summary.TotalHintsViewed))

	// Chapter Breakdown
	b.WriteString(chapterHeaderStyle.Render("Chapter Breakdown:") + "\n")
	for _, cs := range summary.ChapterSummaries {
		chPercent := 0
		chAvgAttempts := 0.0
		if cs.Total > 0 {
			chPercent = (cs.Completed * 100) / cs.Total
			chAvgAttempts = float64(cs.TotalAttempts) / float64(cs.Total)
		}
		b.WriteString(fmt.Sprintf("  • %-20s %3d%% (%d/%d) | Avg Attempts: %.1f | Hints: %d\n",
			cs.ChapterID, chPercent, cs.Completed, cs.Total, chAvgAttempts, cs.TotalHints))
	}

	return b.String()
}

func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "0m"
	}
	hours := int(d.Hours())
	mins := int(d.Minutes()) % 60
	if hours > 0 {
		if mins > 0 {
			return fmt.Sprintf("%dh %dm", hours, mins)
		}
		return fmt.Sprintf("%dh", hours)
	}
	if mins > 0 {
		return fmt.Sprintf("%dm", mins)
	}
	return fmt.Sprintf("%ds", int(d.Seconds()))
}

// FormatFirstRunWelcome renders a styled welcome and getting-started guidance banner.
func FormatFirstRunWelcome() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00D7D7")).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#50FA7B"))

	cmdStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D7D7"))

	dimTxt := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4"))

	var b strings.Builder
	b.WriteString(titleStyle.Render("🚀 Welcome to Terralings!") + "\n\n")
	b.WriteString("New to Terralings or OpenTofu / Terraform? Here is how to get started:\n\n")
	b.WriteString(fmt.Sprintf("  • %-20s %s\n", cmdStyle.Render("terralings tour"), dimTxt.Render("Interactive 5-step guided walkthrough of concepts & tools")))
	b.WriteString(fmt.Sprintf("  • %-20s %s\n", cmdStyle.Render("terralings doctor"), dimTxt.Render("Verify your local IaC engine, cache, and workspace setup")))
	b.WriteString(fmt.Sprintf("  • %-20s %s\n", cmdStyle.Render("terralings watch"), dimTxt.Render("Start continuous interactive exercise watcher")))

	return boxStyle.Render(b.String())
}
