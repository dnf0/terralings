package doctor

import (
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

// CheckStatus represents the status of an individual diagnostic check.
type CheckStatus string

const (
	// StatusPass indicates the check succeeded completely.
	StatusPass CheckStatus = "PASS"
	// StatusWarn indicates a non-fatal warning condition.
	StatusWarn CheckStatus = "WARN"
	// StatusFail indicates a critical blocker.
	StatusFail CheckStatus = "FAIL"
)

// CheckResult records the outcome of a single diagnostic check.
type CheckResult struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Status     CheckStatus `json:"status"`
	Message    string      `json:"message"`
	Resolution string      `json:"resolution,omitempty"`
}

// DiagnosticReport is the complete evaluation outcome of all environment checks.
type DiagnosticReport struct {
	EngineDetected string        `json:"engine_detected"`
	EngineVersion  string        `json:"engine_version"`
	WorkspaceDir   string        `json:"workspace_dir"`
	Checks         []CheckResult `json:"checks"`
	Passed         bool          `json:"passed"`
	FailureCount   int           `json:"failure_count"`
	WarningCount   int           `json:"warning_count"`
}

// RunDiagnostics executes all pre-flight environment checks against the given workspace.
func RunDiagnostics(workspaceDir string, binOverride string, stateOverride string) DiagnosticReport {
	if workspaceDir == "" {
		workspaceDir, _ = os.Getwd()
	}

	var checks []CheckResult
	var failureCount, warningCount int
	var detectedEngine, detectedVersion string

	// 1. IaC Engine Binary check
	binPath, err := detector.DetectBinary(binOverride)
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
		baseName := strings.ToLower(filepath.Base(binPath))
		if strings.Contains(baseName, "tofu") {
			detectedEngine = "opentofu"
		} else {
			detectedEngine = "terraform"
		}

		ver, verErr := detector.GetBinaryVersion(binPath)
		if verErr != nil {
			cmd := exec.Command(binPath, "version")
			if out, cmdErr := cmd.Output(); cmdErr == nil {
				firstLine := strings.Split(strings.TrimSpace(string(out)), "\n")[0]
				detectedVersion = firstLine
			} else {
				detectedVersion = "unknown"
			}
		} else {
			detectedVersion = ver
		}

		if strings.Contains(strings.ToLower(detectedVersion), "opentofu") {
			detectedEngine = "opentofu"
		}

		checks = append(checks, CheckResult{
			ID:      "binary",
			Name:    "IaC Engine Binary",
			Status:  StatusPass,
			Message: fmt.Sprintf("Found %s at %s (%s)", detectedEngine, binPath, detectedVersion),
		})
	}

	// 2. Curriculum Scaffold check
	exercisesPath := filepath.Join(workspaceDir, "exercises")
	if info, err := os.Stat(exercisesPath); err != nil || !info.IsDir() {
		warningCount++
		checks = append(checks, CheckResult{
			ID:         "exercises",
			Name:       "Curriculum Scaffold",
			Status:     StatusWarn,
			Message:    fmt.Sprintf("Directory %q not found in current workspace.", exercisesPath),
			Resolution: "Run `terralings init` to scaffold exercises into this directory.",
		})
	} else {
		var exerciseFiles int
		_ = filepath.Walk(exercisesPath, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr == nil && !info.IsDir() && (strings.HasSuffix(path, ".tf") || strings.HasSuffix(path, ".hcl")) {
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

	// 3. Provider Plugin Cache check
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
		testFile := filepath.Join(cacheDir, ".doctor_write_test")
		if writeErr := os.WriteFile(testFile, []byte("ok"), 0644); writeErr != nil {
			warningCount++
			checks = append(checks, CheckResult{
				ID:         "cache",
				Name:       "Provider Plugin Cache",
				Status:     StatusWarn,
				Message:    fmt.Sprintf("Plugin cache directory %s is not writable: %v", cacheDir, writeErr),
				Resolution: "Check write permissions for ~/.terralings/plugin-cache.",
			})
		} else {
			_ = os.Remove(testFile)
			checks = append(checks, CheckResult{
				ID:      "cache",
				Name:    "Provider Plugin Cache",
				Status:  StatusPass,
				Message: fmt.Sprintf("Plugin cache directory ready at %s", cacheDir),
			})
		}
	}

	// 4. Git Ignore check (checked before state store creation to avoid masking missing .gitignore)
	gitDir := filepath.Join(workspaceDir, ".git")
	if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
		gitignorePath := filepath.Join(workspaceDir, ".gitignore")
		content, readErr := os.ReadFile(gitignorePath)
		if readErr == nil && strings.Contains(string(content), ".terralings") {
			checks = append(checks, CheckResult{
				ID:      "git",
				Name:    "Git Ignore Integration",
				Status:  StatusPass,
				Message: ".terralings directory is properly git-ignored.",
			})
		} else {
			warningCount++
			checks = append(checks, CheckResult{
				ID:         "git",
				Name:       "Git Ignore Integration",
				Status:     StatusWarn,
				Message:    ".terralings directory is not listed in .gitignore.",
				Resolution: "Add `.terralings/` to `.gitignore` to avoid committing local progress.",
			})
		}
	}

	// 5. Progress Persistence State check
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

// FormatReport creates a styled terminal representation of the diagnostic report.
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
