package tour

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Step struct {
	Index          int      `json:"index"`
	Title          string   `json:"title"`
	Subtitle       string   `json:"subtitle"`
	Body           []string `json:"body"`
	CommandExample string   `json:"command_example,omitempty"`
	KeyTakeaways   []string `json:"key_takeaways,omitempty"`
}

func DefaultSteps() []Step {
	return []Step{
		{
			Index:    1,
			Title:    "Welcome & Core Philosophy",
			Subtitle: "Master Terraform & OpenTofu through interactive hands-on practice",
			Body: []string{
				"Terralings is designed to teach you Infrastructure-as-Code from first principles.",
				"All exercises run in isolated, sandboxed environments without requiring real cloud credentials.",
				"We follow the Ziglings / Rustlings v6 model: pure deterministic validation with zero magic comment friction.",
			},
			CommandExample: "terralings watch",
			KeyTakeaways: []string{
				"100% local, safe evaluation with OpenTofu / Terraform.",
				"Real compiler errors & plan outputs guide your progress.",
			},
		},
		{
			Index:    2,
			Title:    "Anatomy of an Exercise",
			Subtitle: "How exercises are structured and solved",
			Body: []string{
				"Exercises live in the `exercises/` folder across 13 progressive chapters.",
				"Each file starts with `# TODO:` instructions explaining the required infrastructure declaration.",
				"Every exercise begins in a failing state (syntax error, missing block, or failed test assertion).",
			},
			CommandExample: "code exercises/01_primitives/primitives01.tf",
			KeyTakeaways: []string{
				"Fix the deliberate bug to make the exercise pass.",
				"Reference solutions are available in `solutions/` for comparison.",
			},
		},
		{
			Index:    3,
			Title:    "Continuous Watch & Verification",
			Subtitle: "The rapid edit-save-verify feedback loop",
			Body: []string{
				"Running `terralings watch` starts continuous file monitoring with instant re-evaluation on save.",
				"When an exercise passes, the watcher pauses with interactive controls:",
				"  [Enter / n] Next exercise  |  [p] Previous  |  [r] Rerun  |  [q] Quit",
			},
			CommandExample: "terralings watch",
			KeyTakeaways: []string{
				"Run single exercises on demand with `terralings run <name>`.",
				"Verify all solutions at any time with `terralings verify`.",
			},
		},
		{
			Index:    4,
			Title:    "Interactive TUI, Hints & Analytics",
			Subtitle: "Powerful terminal dashboard, multi-level hints, and progress tracking",
			Body: []string{
				"Launch `terralings tui` (or `watch -i`) for a full-screen Bubble Tea split-pane dashboard.",
				"Stuck on an exercise? Get progressive hints with `terralings hint <name>` or press 'h' in the TUI.",
				"View your learning stats, attempts, and chapter completion with `terralings stats`.",
			},
			CommandExample: "terralings tui\nterralings hint primitives01\nterralings stats",
			KeyTakeaways: []string{
				"Fuzzy search curriculum topics anytime with `terralings search <term>`.",
				"Reset any exercise back to its clean starting template with `terralings reset <name>`.",
			},
		},
		{
			Index:    5,
			Title:    "Editor Integration & LSP",
			Subtitle: "Real-time compiler diagnostics and hover docs in your favorite editor",
			Body: []string{
				"Terralings includes a built-in Language Server Protocol daemon: `terralings lsp`.",
				"Configure your editor to receive inline OpenTofu/Terraform error diagnostics and rich markdown hint cards.",
				"Works seamlessly with VS Code, Neovim (`nvim-lspconfig`), and Helix (`languages.toml`).",
			},
			CommandExample: "terralings lsp",
			KeyTakeaways: []string{
				"Zero fake warning squiggles — only true HCL syntax and test errors.",
				"Instant code actions to reveal progressive hints right in your editor.",
			},
		},
	}
}

type Tour struct {
	Steps          []Step
	Writer         io.Writer
	Reader         io.Reader
	NonInteractive bool
	JSONMode       bool
}

func NewTour(w io.Writer, r io.Reader) *Tour {
	if w == nil {
		w = os.Stdout
	}
	if r == nil {
		r = os.Stdin
	}
	return &Tour{
		Steps:  DefaultSteps(),
		Writer: w,
		Reader: r,
	}
}

func (t *Tour) ExportJSON() ([]byte, error) {
	return json.MarshalIndent(t.Steps, "", "  ")
}

func (t *Tour) RenderStep(stepIndex int) error {
	if stepIndex < 1 || stepIndex > len(t.Steps) {
		return fmt.Errorf("invalid step index %d (must be between 1 and %d)", stepIndex, len(t.Steps))
	}
	s := t.Steps[stepIndex-1]

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#50FA7B"))
	stepBadge := lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#6272A4")).
		Foreground(lipgloss.Color("#F8F8F2")).
		Padding(0, 1)
	subStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#F1FA8C"))
	bodyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F8F8F2"))
	codeBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00D7D7")).
		Padding(0, 1).
		Foreground(lipgloss.Color("#50FA7B"))
	takeawayStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8BE9FD"))

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%s %s\n", stepBadge.Render(fmt.Sprintf("STEP %d OF %d", s.Index, len(t.Steps))), titleStyle.Render(s.Title)))
	b.WriteString(subStyle.Render(s.Subtitle) + "\n\n")

	for _, line := range s.Body {
		b.WriteString(bodyStyle.Render("  "+line) + "\n")
	}

	if s.CommandExample != "" {
		b.WriteString("\n  Example Command:\n")
		b.WriteString(codeBoxStyle.Render(s.CommandExample) + "\n")
	}

	if len(s.KeyTakeaways) > 0 {
		b.WriteString("\n  Key Takeaways:\n")
		for _, item := range s.KeyTakeaways {
			b.WriteString(takeawayStyle.Render("  ✓ "+item) + "\n")
		}
	}
	b.WriteString("\n")

	_, err := io.WriteString(t.Writer, b.String())
	return err
}

func (t *Tour) Run(ctx context.Context, startStep int) error {
	if t.JSONMode {
		data, err := t.ExportJSON()
		if err != nil {
			return err
		}
		_, err = t.Writer.Write(data)
		return err
	}

	if t.NonInteractive {
		if startStep >= 1 && startStep <= len(t.Steps) {
			return t.RenderStep(startStep)
		}
		if startStep > len(t.Steps) || startStep < 0 {
			return fmt.Errorf("invalid step index %d (must be between 1 and %d)", startStep, len(t.Steps))
		}
		for i := 1; i <= len(t.Steps); i++ {
			if err := t.RenderStep(i); err != nil {
				return err
			}
			if i < len(t.Steps) {
				io.WriteString(t.Writer, strings.Repeat("─", 60)+"\n")
			}
		}
		return nil
	}

	// Interactive Loop
	current := startStep
	if current < 1 || current > len(t.Steps) {
		current = 1
	}

	scanner := bufio.NewScanner(t.Reader)

	promptStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#BD93F9"))
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D7D7"))
	dim := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4"))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := t.RenderStep(current); err != nil {
			return err
		}

		prompt := fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s %s\n> ",
			keyStyle.Render("[Enter / n]"), promptStyle.Render("Next"),
			dim.Render("|"),
			keyStyle.Render("[p]"), promptStyle.Render("Prev"),
			dim.Render("|"),
			keyStyle.Render("[1-5]"), promptStyle.Render("Jump"),
			dim.Render("|"),
			keyStyle.Render("[q]"), promptStyle.Render("Quit"),
		)
		io.WriteString(t.Writer, prompt)

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(strings.ToLower(scanner.Text()))
		switch input {
		case "", "n", "next", "right":
			if current < len(t.Steps) {
				current++
			} else {
				io.WriteString(t.Writer, "\n🎉 You've reached the end of the tour! Run `terralings watch` to begin learning.\n\n")
				return nil
			}
		case "p", "prev", "previous", "left":
			if current > 1 {
				current--
			}
		case "q", "quit", "exit":
			io.WriteString(t.Writer, "\nExited tour. Happy learning!\n\n")
			return nil
		case "r", "rerun":
			// re-render current
		default:
			if num, err := strconv.Atoi(input); err == nil && num >= 1 && num <= len(t.Steps) {
				current = num
			}
		}
	}

	return scanner.Err()
}
