# Specification: Terralings Onboarding Experience & Environment Diagnostics

**Date:** 2026-08-27  
**Status:** Approved  
**Author:** Pair Programming Agent & Daniel Fisher  
**Scope:** `internal/tour`, `internal/doctor`, CLI commands (`tour`, `doctor`), First-run recommendation hooks, and `docs/onboarding-guide.md`.

---

## 1. Executive Summary & Goals

The goal of this specification is to create a seamless, guided onboarding experience for **Terralings** learners, matching the gold-standard architecture established in Kubelings.

Key objectives:
1. **Interactive CLI Tour Engine (`terralings tour`)**: A rich, multi-step interactive terminal walkthrough teaching learners how Terralings works, how exercises are structured, how to use watch/TUI modes, progressive hints, and editor LSP integration.
2. **Environment Pre-flight Diagnostics (`terralings doctor`)**: A comprehensive health-check tool verifying local OpenTofu/Terraform installations, version compatibility, directory scaffold integrity, plugin cache permissions, and state persistence.
3. **First-Run Guidance**: Friendly, non-intrusive terminal callouts on `terralings init` and `terralings` root executions when a new user has 0 attempts recorded.
4. **Headless & Editor Integration**: `--json` and `--non-interactive` flags on `tour` and `doctor` to empower IDE extensions (VS Code, Neovim, Helix) and automated workflows.
5. **Comprehensive Onboarding Documentation (`docs/onboarding-guide.md`)**: A complete beginner-to-advanced guide covering installation, doctor checks, tour steps, and editor workflows.

---

## 2. Architecture & Component Design

```
                     ┌─────────────────────────────────────────────────┐
                     │              cmd/terralings/main.go             │
                     │  (subcommands: tour, doctor, init, watch, tui)  │
                     └───────────────┬─────────────────┬───────────────┘
                                     │                 │
                  ┌──────────────────▼──┐           ┌──▼──────────────────┐
                  │    internal/tour    │           │   internal/doctor   │
                  │  - 5-Step Engine    │           │  - Binary Checks    │
                  │  - Interactive TTY  │           │  - Exercises Check  │
                  │  - NDJSON / JSON    │           │  - Cache Dir Check  │
                  │  - Lip Gloss UI     │           │  - State File Check │
                  └─────────────────────┘           └─────────────────────┘
```

### 2.1 The CLI Tour Engine (`internal/tour`)

The tour engine provides a guided 5-step walkthrough rendered with **Lip Gloss** styling:

#### Steps Breakdown:
1. **Step 1: Welcome & Core Philosophy**
   - Introduces Infrastructure-as-Code learning with OpenTofu & Terraform.
   - Explains sandboxed isolated workspace evaluation without cloud credentials.
   - Explains the pure deterministic validation model (Zero comment busywork).
2. **Step 2: Anatomy of an Exercise**
   - Explains `exercises/` folder organization across 13 progressive chapters.
   - Explains `# TODO:` task markers and deliberate educational bugs/stubs.
   - Shows how exercises and reference solutions (`solutions/`) correspond.
3. **Step 3: Continuous Watch & Verification**
   - Teaches `terralings watch` workflow (`Edit -> Save -> Green -> Next`).
   - Explains the interactive pass prompt (`[Enter / n] Next | [p] Prev | [r] Rerun`).
   - Explains single-exercise runs via `terralings run <exercise>`.
4. **Step 4: Interactive TUI, Hints & Analytics**
   - Teaches the full-screen terminal dashboard: `terralings tui` / `watch -i`.
   - Explains progressive hint extraction (`terralings hint <exercise> [-i <idx>]`).
   - Explains learning statistics and progress tracking (`terralings stats`).
   - Explains template reset (`terralings reset <exercise>`).
5. **Step 5: Editor Integration & Language Server (LSP)**
   - Teaches `terralings lsp` daemon usage for inline compiler diagnostics and markdown hover docs.
   - Provides quick setup instructions for VS Code, Neovim, and Helix.

#### Data Structures (`internal/tour/tour.go`):
```go
package tour

import (
	"context"
	"io"
)

type Step struct {
	Index          int      `json:"index"`
	Title          string   `json:"title"`
	Subtitle       string   `json:"subtitle"`
	Body           []string `json:"body"`
	CommandExample string   `json:"command_example,omitempty"`
	KeyTakeaways   []string `json:"key_takeaways,omitempty"`
}

type Tour struct {
	Steps          []Step
	CurrentIndex   int
	Writer         io.Writer
	Reader         io.Reader
	NonInteractive bool
	JSONMode       bool
}

func NewTour(w io.Writer, r io.Reader) *Tour
func (t *Tour) Run(ctx context.Context, startStep int) error
func (t *Tour) RenderStep(stepIndex int) error
func (t *Tour) ExportJSON() ([]byte, error)
```

#### Interactive Navigation Controls:
- `Enter` or `n` / `Right`: Advance to next step.
- `p` or `Left`: Go back to previous step.
- `1` - `5`: Direct jump to specified step.
- `r`: Re-render current step.
- `q` or `Ctrl+C`: Exit tour cleanly.

---

### 2.2 Environment Diagnostics Doctor (`internal/doctor`)

The diagnostic engine performs pre-flight verification of the user's environment to detect and resolve configuration issues before the user encounters runtime errors.

#### Diagnostic Checks:
1. **Engine Binary Detection**:
   - Checks whether `tofu` or `terraform` is discoverable via PATH or `--bin` / `TERRALINGS_BIN`.
   - Validates execution permissions and extracts version string (`tofu version` / `terraform version`).
   - Status: `PASS` if found, `FAIL` if neither exists (with install instructions for Homebrew, curl, apt).
2. **Curriculum Scaffold Integrity**:
   - Verifies whether `exercises/` exists in current workspace.
   - Validates presence of 13 chapter folders and 56 exercise configurations.
   - Status: `PASS` if complete, `WARN` if missing (recommending `terralings init`).
3. **Provider Plugin Cache**:
   - Verifies provider cache directory (`runner.PluginCacheDir()`).
   - Checks directory creation and write permissions.
   - Status: `PASS` if writable, `WARN` if fallback local cache is used.
4. **Progress State Persistence**:
   - Verifies `.terralings/state.json` store path.
   - Validates JSON parseability and write permissions.
   - Status: `PASS` if healthy, `WARN` if corrupt (noting automatic `.bak` recovery).
5. **Git Repository & Ignore Rules**:
   - Checks if workspace is a git repository (`.git`).
   - Checks if `.terralings/` is ignored in `.gitignore`.
   - Status: `PASS` if properly ignored, `WARN` if missing from `.gitignore`.

#### Data Structures (`internal/doctor/doctor.go`):
```go
package doctor

type CheckStatus string

const (
	StatusPass CheckStatus = "PASS"
	StatusWarn CheckStatus = "WARN"
	StatusFail CheckStatus = "FAIL"
)

type CheckResult struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Status      CheckStatus `json:"status"`
	Message     string      `json:"message"`
	Resolution  string      `json:"resolution,omitempty"`
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

func RunDiagnostics(workspaceDir string, binOverride string, stateOverride string) DiagnosticReport
func FormatReport(report DiagnosticReport) string
```

---

### 2.3 First-Run Guidance Hooks

When a learner executes `terralings` or `terralings init`:
- The CLI inspects the active `state.Store`.
- If `analytics.TotalAttempts == 0 && analytics.CompletedCount == 0`:
  - A friendly Lip Gloss callout banner is displayed:
  ```text
  👋 Welcome to Terralings! New here?
     Run `terralings tour` for a 2-minute interactive guided walkthrough.
     Run `terralings doctor` to verify your OpenTofu/Terraform setup.
  ```

---

### 2.4 Comprehensive Onboarding Guide (`docs/onboarding-guide.md`)

Authored at `docs/onboarding-guide.md`, containing:
1. **Quickstart Checklist**: From zero to first solved exercise.
2. **Environment Diagnostics**: Troubleshooting binary detection with `terralings doctor`.
3. **CLI Guided Tour**: Step-by-step summary of the 5 tour modules.
4. **Learning Modalities**: When to use standard watch mode (`terralings watch`) vs full TUI dashboard (`terralings tui`).
5. **Editor Setup Walkthrough**: Full configurations for VS Code, Neovim (`nvim-lspconfig`), and Helix (`languages.toml`).
6. **Curriculum Roadmap**: Overview of all 13 chapters and key competencies.

---

## 3. Testing Strategy

1. **`test/tour_test.go`**:
   - Verify all 5 tour steps are present with complete metadata.
   - Verify interactive navigation (`n`, `p`, direct numbers `1-5`, `q`).
   - Verify `--non-interactive` mode prints all steps cleanly without blocking on stdin.
   - Verify `--step <n>` renders the exact requested step.
   - Verify `--json` outputs valid JSON containing all step objects.
2. **`test/doctor_test.go`**:
   - Verify detection of OpenTofu and Terraform engines.
   - Verify handling when binary is missing (returns `StatusFail` with installation instructions).
   - Verify exercises directory check with complete vs partial scaffold.
   - Verify state store permissions and JSON diagnostics output.
3. **`test/cli_onboarding_test.go`**:
   - End-to-end CLI integration testing of `terralings tour` and `terralings doctor`.
   - Test first-run welcome hook output on clean state.

---

## 4. Verification & Release Checklist

- `make check`: `go mod verify`, `go vet ./...`, `gofmt`.
- `make all`: Full build and test suite pass cleanly with `-race`.
- `roborev status`: Check automated review feedback.
