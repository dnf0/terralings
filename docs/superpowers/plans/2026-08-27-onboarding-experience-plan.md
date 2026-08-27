# Onboarding Experience & Environment Diagnostics Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a complete, friendly onboarding experience and environment diagnostics toolset (`terralings tour`, `terralings doctor`, first-run welcome hooks, and `docs/onboarding-guide.md`) for Terralings.

**Architecture:** 
- `internal/tour`: 5-step interactive terminal tour with Lip Gloss styling, key navigation, step jumping, non-interactive mode, and JSON export.
- `internal/doctor`: Pre-flight environment diagnostics verifying OpenTofu/Terraform binary availability, exercise directory scaffold integrity, provider plugin cache directory permissions, state store health, and git repository ignore rules.
- `cmd/terralings`: Cobra subcommands `tour` and `doctor` with first-run guidance hooks on clean state store.
- `docs/onboarding-guide.md`: Illustrated beginner-to-advanced onboarding manual.

**Tech Stack:** Go 1.22+, Cobra, Lip Gloss, OpenTofu / Terraform CLI.

---

### Task 1: CLI Tour Engine (`internal/tour`) & Unit Tests

**Files:**
- Create: `internal/tour/tour.go`
- Create: `test/tour_test.go`

- [ ] **Step 1: Write unit tests in `test/tour_test.go`**

```go
package test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/tour"
)

func TestTour_DefaultStepsCountAndContent(t *testing.T) {
	steps := tour.DefaultSteps()
	if len(steps) != 5 {
		t.Fatalf("expected 5 default steps, got %d", len(steps))
	}

	expectedTitles := []string{
		"Welcome & Core Philosophy",
		"Anatomy of an Exercise",
		"Continuous Watch & Verification",
		"Interactive TUI, Hints & Analytics",
		"Editor Integration & LSP",
	}

	for i, expected := range expectedTitles {
		if steps[i].Index != i+1 {
			t.Errorf("step %d index mismatch: got %d, want %d", i, steps[i].Index, i+1)
		}
		if steps[i].Title != expected {
			t.Errorf("step %d title mismatch: got %q, want %q", i, steps[i].Title, expected)
		}
		if len(steps[i].Body) == 0 {
			t.Errorf("step %d body should not be empty", i)
		}
	}
}

func TestTour_NonInteractiveAllSteps(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("")

	tr := tour.NewTour(&out, in)
	tr.NonInteractive = true

	err := tr.Run(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error running non-interactive tour: %v", err)
	}

	output := out.String()
	for _, title := range []string{
		"Welcome & Core Philosophy",
		"Anatomy of an Exercise",
		"Continuous Watch & Verification",
		"Interactive TUI, Hints & Analytics",
		"Editor Integration & LSP",
	} {
		if !strings.Contains(output, title) {
			t.Errorf("output missing expected step title %q", title)
		}
	}
}

func TestTour_SpecificStepRender(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("")

	tr := tour.NewTour(&out, in)
	tr.NonInteractive = true

	err := tr.Run(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Continuous Watch & Verification") {
		t.Errorf("expected step 3 output, got: %s", output)
	}
	if strings.Contains(output, "Welcome & Core Philosophy") {
		t.Errorf("should not contain step 1 output when step 3 requested")
	}
}

func TestTour_ExportJSON(t *testing.T) {
	tr := tour.NewTour(nil, nil)
	jsonData, err := tr.ExportJSON()
	if err != nil {
		t.Fatalf("failed to export tour JSON: %v", err)
	}

	var parsed []tour.Step
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if len(parsed) != 5 {
		t.Fatalf("expected 5 steps in JSON, got %d", len(parsed))
	}
}

func TestTour_InteractiveNavigation(t *testing.T) {
	var out bytes.Buffer
	// Navigate: next (n) -> next (Enter) -> prev (p) -> jump (5) -> quit (q)
	in := strings.NewReader("n\n\np\n5\nq\n")

	tr := tour.NewTour(&out, in)
	tr.NonInteractive = false

	err := tr.Run(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error in interactive tour: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Editor Integration & LSP") {
		t.Errorf("expected jump to step 5 in output, got: %s", output)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/tour_test.go`
Expected: FAIL with undefined `tour` package.

- [ ] **Step 3: Implement `internal/tour/tour.go`**

```go
package tour

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		if startStep > 1 && startStep <= len(t.Steps) {
			return t.RenderStep(startStep)
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v -race ./test/tour_test.go`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/tour/tour.go test/tour_test.go
git commit -m "feat(tour): implement 5-step interactive onboarding tour engine"
```

---

### Task 2: Environment Diagnostics Doctor (`internal/doctor`) & Unit Tests

**Files:**
- Create: `internal/doctor/doctor.go`
- Create: `test/doctor_test.go`

- [ ] **Step 1: Write unit tests in `test/doctor_test.go`**

```go
package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/doctor"
)

func TestDoctor_RunDiagnostics(t *testing.T) {
	tmpDir := t.TempDir()

	report := doctor.RunDiagnostics(tmpDir, "", "")
	if len(report.Checks) < 4 {
		t.Fatalf("expected at least 4 diagnostic checks, got %d", len(report.Checks))
	}

	// Should contain binary, exercises, cache, state checks
	var foundBinary, foundExercises, foundCache, foundState bool
	for _, check := range report.Checks {
		switch check.ID {
		case "binary":
			foundBinary = true
		case "exercises":
			foundExercises = true
		case "cache":
			foundCache = true
		case "state":
			foundState = true
		}
	}

	if !foundBinary || !foundExercises || !foundCache || !foundState {
		t.Errorf("missing core checks: binary=%v, exercises=%v, cache=%v, state=%v",
			foundBinary, foundExercises, foundCache, foundState)
	}
}

func TestDoctor_FormatReport(t *testing.T) {
	rep := doctor.DiagnosticReport{
		EngineDetected: "opentofu",
		EngineVersion:  "1.8.0",
		WorkspaceDir:   "/tmp/test",
		Checks: []doctor.CheckResult{
			{
				ID:      "binary",
				Name:    "IaC Engine Binary",
				Status:  doctor.StatusPass,
				Message: "Found OpenTofu v1.8.0",
			},
			{
				ID:         "exercises",
				Name:       "Curriculum Scaffold",
				Status:     doctor.StatusWarn,
				Message:    "Missing exercises directory",
				Resolution: "Run `terralings init` to scaffold",
			},
		},
		Passed:       true,
		FailureCount: 0,
		WarningCount: 1,
	}

	formatted := doctor.FormatReport(rep)
	if !strings.Contains(formatted, "Terralings Doctor Report") {
		t.Errorf("expected header in formatted report, got: %s", formatted)
	}
	if !strings.Contains(formatted, "IaC Engine Binary") {
		t.Errorf("expected binary check in formatted report")
	}
	if !strings.Contains(formatted, "Run `terralings init`") {
		t.Errorf("expected resolution message in formatted report")
	}
}

func TestDoctor_JSONSerialization(t *testing.T) {
	tmpDir := t.TempDir()
	rep := doctor.RunDiagnostics(tmpDir, "", "")

	jsonData, err := json.Marshal(rep)
	if err != nil {
		t.Fatalf("failed to serialize report to JSON: %v", err)
	}

	var decoded doctor.DiagnosticReport
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("failed to deserialize report JSON: %v", err)
	}

	if decoded.WorkspaceDir != tmpDir {
		t.Errorf("workspace mismatch: got %q, want %q", decoded.WorkspaceDir, tmpDir)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/doctor_test.go`
Expected: FAIL with undefined `doctor` package.

- [ ] **Step 3: Implement `internal/doctor/doctor.go`**

```go
package doctor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dnf0/terralings/internal/detector"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/state"
)

type CheckStatus string

const (
	StatusPass CheckStatus = "PASS"
	StatusWarn CheckStatus = "WARN"
	StatusFail CheckStatus = "FAIL"
)

type CheckResult struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Status     CheckStatus `json:"status"`
	Message    string      `json:"message"`
	Resolution string      `json:"resolution,omitempty"`
}

type DiagnosticReport struct {
	EngineDetected string        `json:"engine_detected"`
	EngineVersion  string        `json:"engine_version"`
	WorkspaceDir   string        `json:"workspace_dir"`
	Checks         []CheckResult `json:"checks"`
	Passed         bool          `json:"passed"`
	FailureCount   int           `json:"failure_count"`
	WarningCount   int           `json:"warning_count"`
}

func RunDiagnostics(workspaceDir string, binOverride string, stateOverride string) DiagnosticReport {
	if workspaceDir == "" {
		workspaceDir, _ = os.Getwd()
	}

	var checks []CheckResult
	var failureCount, warningCount int
	var detectedEngine, detectedVersion string

	// 1. Binary check
	binInfo, err := detector.DetectBinary(binOverride)
	if err != nil {
		failureCount++
		checks = append(checks, CheckResult{
			ID:         "binary",
			Name:       "IaC Engine Binary",
			Status:     StatusFail,
			Message:    "Neither OpenTofu nor Terraform was detected in your PATH.",
			Resolution: "Install OpenTofu (`brew install opentofu` or `curl -fsSL https://get.opentofu.org/install-opentofu.sh | sh`) or specify `--bin <path>`.",
		})
	} else {
		detectedEngine = string(binInfo.Engine)
		cmd := exec.Command(binInfo.Path, "version")
		if out, err := cmd.Output(); err == nil {
			firstLine := strings.Split(string(out), "\n")[0]
			detectedVersion = strings.TrimSpace(firstLine)
		}
		checks = append(checks, CheckResult{
			ID:      "binary",
			Name:    "IaC Engine Binary",
			Status:  StatusPass,
			Message: fmt.Sprintf("Found %s at %s (%s)", binInfo.Engine, binInfo.Path, detectedVersion),
		})
	}

	// 2. Exercises directory check
	exercisesPath := filepath.Join(workspaceDir, "exercises")
	if info, err := os.Stat(exercisesPath); err != nil || !info.IsDir() {
		warningCount++
		checks = append(checks, CheckResult{
			ID:         "exercises",
			Name:       "Curriculum Scaffold",
			Status:     StatusWarn,
			Message:    fmt.Sprintf("Directory %q not found in current workspace.", exercisesPath),
			Resolution: "Run `terralings init` to extract exercises into this directory.",
		})
	} else {
		// Count exercises
		var exerciseFiles int
		_ = filepath.Walk(exercisesPath, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && (strings.HasSuffix(path, ".tf") || strings.HasSuffix(path, ".hcl")) {
				exerciseFiles++
			}
			return nil
		})
		checks = append(checks, CheckResult{
			ID:      "exercises",
			Name:    "Curriculum Scaffold",
			Status:  StatusPass,
			Message: fmt.Sprintf("Exercises directory present (%d configuration files found)", exerciseFiles),
		})
	}

	// 3. Provider plugin cache check
	cacheDir := runner.PluginCacheDir()
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		warningCount++
		checks = append(checks, CheckResult{
			ID:         "cache",
			Name:       "Provider Plugin Cache",
			Status:     StatusWarn,
			Message:    fmt.Sprintf("Cannot create plugin cache directory %s: %v", cacheDir, err),
			Resolution: "Ensure write permissions for ~/.terralings or set TERRALINGS_PLUGIN_CACHE_DIR.",
		})
	} else {
		checks = append(checks, CheckResult{
			ID:      "cache",
			Name:    "Provider Plugin Cache",
			Status:  StatusPass,
			Message: fmt.Sprintf("Plugin cache directory ready at %s", cacheDir),
		})
	}

	// 4. State store check
	statePath := stateOverride
	if statePath == "" {
		statePath = filepath.Join(workspaceDir, ".terralings", "state.json")
	}
	store, err := state.NewStore(statePath)
	if err != nil {
		warningCount++
		checks = append(checks, CheckResult{
			ID:         "state",
			Name:       "Progress Persistence Store",
			Status:     StatusWarn,
			Message:    fmt.Sprintf("State store error: %v", err),
			Resolution: "Check file permissions in .terralings/",
		})
	} else {
		analytics := store.GetAnalytics(nil)
		checks = append(checks, CheckResult{
			ID:      "state",
			Name:    "Progress Persistence Store",
			Status:  StatusPass,
			Message: fmt.Sprintf("State store healthy at %s (%d completed, %d attempts)", statePath, analytics.CompletedCount, analytics.TotalAttempts),
		})
	}

	// 5. Git ignore check
	gitDir := filepath.Join(workspaceDir, ".git")
	if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
		gitignorePath := filepath.Join(workspaceDir, ".gitignore")
		content, _ := os.ReadFile(gitignorePath)
		if strings.Contains(string(content), ".terralings") {
			checks = append(checks, CheckResult{
				ID:      "git",
				Name:    "Git Ignore Integration",
				Status:  StatusPass,
				Message: ".terralings directory is properly git-ignored.",
			})
		} else {
			checks = append(checks, CheckResult{
				ID:         "git",
				Name:       "Git Ignore Integration",
				Status:     StatusWarn,
				Message:    ".terralings directory is not listed in .gitignore.",
				Resolution: "Add `.terralings/` to `.gitignore` to avoid committing local progress.",
			})
		}
	}

	return DiagnosticReport{
		EngineDetected: detectedEngine,
		EngineVersion:  detectedVersion,
		WorkspaceDir:   workspaceDir,
		Checks:         checks,
		Passed:         failureCount == 0,
		FailureCount:   failureCount,
		WarningCount:   warningCount,
	}
}

func FormatReport(report DiagnosticReport) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#50FA7B")).
		Padding(0, 1)

	passStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#50FA7B"))
	warnStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFB86C"))
	failStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5555"))
	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4"))
	resStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8BE9FD")).
		Italic(true)

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(headerStyle.Render("🩺 Terralings Doctor Diagnostic Report"))
	b.WriteString("\n" + strings.Repeat("─", 60) + "\n\n")

	for _, check := range report.Checks {
		var icon string
		switch check.Status {
		case StatusPass:
			icon = passStyle.Render("✓")
		case StatusWarn:
			icon = warnStyle.Render("!")
		case StatusFail:
			icon = failStyle.Render("✗")
		}

		b.WriteString(fmt.Sprintf(" %s %s\n", icon, lipgloss.NewStyle().Bold(true).Render(check.Name)))
		b.WriteString(fmt.Sprintf("   %s\n", check.Message))
		if check.Resolution != "" {
			b.WriteString(fmt.Sprintf("   %s %s\n", dimStyle.Render("→ Fix:"), resStyle.Render(check.Resolution)))
		}
		b.WriteString("\n")
	}

	b.WriteString(strings.Repeat("─", 60) + "\n")
	if report.FailureCount == 0 && report.WarningCount == 0 {
		b.WriteString(passStyle.Render(" All diagnostics passed! Your environment is 100% ready for Terralings.\n\n"))
	} else if report.FailureCount == 0 {
		b.WriteString(warnStyle.Render(fmt.Sprintf(" Diagnostics passed with %d warning(s). You can run Terralings safely.\n\n", report.WarningCount)))
	} else {
		b.WriteString(failStyle.Render(fmt.Sprintf(" Found %d critical issue(s). Please resolve them before continuing.\n\n", report.FailureCount)))
	}

	return b.String()
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v -race ./test/doctor_test.go`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/doctor/doctor.go test/doctor_test.go
git commit -m "feat(doctor): implement pre-flight environment diagnostics engine"
```

---

### Task 3: CLI Subcommand Registration & First-Run Guidance Hooks

**Files:**
- Modify: `cmd/terralings/main.go`
- Create: `test/cli_onboarding_test.go`

- [ ] **Step 1: Write integration tests in `test/cli_onboarding_test.go`**

```go
package test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/doctor"
	"github.com/dnf0/terralings/internal/tour"
)

func TestCLI_TourJSONOutput(t *testing.T) {
	var buf bytes.Buffer
	tr := tour.NewTour(&buf, nil)
	tr.JSONMode = true

	if err := tr.Run(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var steps []tour.Step
	if err := json.Unmarshal(buf.Bytes(), &steps); err != nil {
		t.Fatalf("failed to decode tour JSON: %v", err)
	}

	if len(steps) != 5 {
		t.Errorf("expected 5 steps, got %d", len(steps))
	}
}

func TestCLI_DoctorExecution(t *testing.T) {
	tmpDir := t.TempDir()
	rep := doctor.RunDiagnostics(tmpDir, "", "")

	formatted := doctor.FormatReport(rep)
	if !strings.Contains(formatted, "Terralings Doctor Diagnostic Report") {
		t.Errorf("expected doctor banner in formatted report")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/cli_onboarding_test.go`
Expected: FAIL or PASS depending on underlying packages.

- [ ] **Step 3: Update `cmd/terralings/main.go`**

Register `tour` and `doctor` commands with flags:
- `tour`: `--step <int>`, `--non-interactive`, `--json`.
- `doctor`: `--json`.
- Add first-run welcome prompt to `mainCmd` / `initCmd`.

- [ ] **Step 4: Run all CLI tests to verify they pass**

Run: `go test -v -race ./test/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/terralings/main.go test/cli_onboarding_test.go
git commit -m "feat(cli): add tour and doctor subcommands with first-run guidance hooks"
```

---

### Task 4: Comprehensive Onboarding Documentation & Release Verification

**Files:**
- Create: `docs/onboarding-guide.md`
- Modify: `README.md`
- Test: All test suites (`make all`)

- [ ] **Step 1: Create `docs/onboarding-guide.md`**
Author comprehensive illustrated guide covering:
1. System requirements & pre-flight diagnostics (`terralings doctor`).
2. Scaffold initialization (`terralings init`).
3. CLI Guided Tour (`terralings tour`).
4. Learning modalities (`terralings watch` vs `terralings tui`).
5. Editor configurations (VS Code, Neovim, Helix).
6. Troubleshooting & FAQs.

- [ ] **Step 2: Update `README.md`**
Add `terralings tour` and `terralings doctor` to CLI command table and quickstart section.

- [ ] **Step 3: Run full verification suite**
- `make check`
- `make all`

- [ ] **Step 4: Commit**

```bash
git add docs/onboarding-guide.md README.md
git commit -m "docs: add comprehensive onboarding guide and update CLI reference"
```

---
