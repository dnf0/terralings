# Terralings Milestone 2 Implementation Plan: TUI, Analytics Engine & Editor LSP Integration

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build and deliver Milestone 2 for Terralings: progress state persistence and learning analytics, an interactive Bubble Tea terminal dashboard, and an LSP diagnostic language server with NDJSON event streaming.

**Architecture:** Modular architecture separating state management (`internal/state`), terminal UI rendering (`internal/tui`), compiler diagnostic normalization (`internal/diagnostics`), and JSON-RPC LSP protocol handling (`internal/lsp`), unified through Cobra commands in `cmd/terralings/main.go`.

**Tech Stack:** Go 1.23, `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/bubbles`, `github.com/charmbracelet/lipgloss`, `github.com/fsnotify/fsnotify`, `github.com/spf13/cobra`.

---

### Task 1: Progress State & Learning Analytics Engine

**Files:**
- Create: `internal/state/state.go`
- Create: `internal/state/analytics.go`
- Create: `test/state_test.go`

- [ ] **Step 1: Write state and analytics tests in `test/state_test.go`**

```go
package test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/state"
)

func TestState_NewStoreAndAutoInitialize(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, ".terralings", "state.json")

	store, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to initialize state store: %v", err)
	}

	if store.GetVersion() != "1.0" {
		t.Errorf("Expected version 1.0, got %s", store.GetVersion())
	}

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Errorf("Expected state file %s to be created on disk", statePath)
	}
}

func TestState_RecordAttemptAndHints(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, ".terralings", "state.json")

	store, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	// 1. Record failed attempt
	if err := store.RecordAttempt("primitives01", "01_primitives", false); err != nil {
		t.Fatalf("RecordAttempt failed: %v", err)
	}

	exState := store.GetExerciseState("primitives01")
	if exState.Attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", exState.Attempts)
	}
	if exState.Status != state.StatusInProgress {
		t.Errorf("Expected in_progress, got %s", exState.Status)
	}

	// 2. Record hint view
	if err := store.RecordHint("primitives01", "01_primitives", 1); err != nil {
		t.Fatalf("RecordHint failed: %v", err)
	}
	exState = store.GetExerciseState("primitives01")
	if exState.HintsViewed != 1 {
		t.Errorf("Expected 1 hint viewed, got %d", exState.HintsViewed)
	}

	// 3. Record passed attempt
	if err := store.RecordAttempt("primitives01", "01_primitives", true); err != nil {
		t.Fatalf("RecordAttempt failed: %v", err)
	}
	exState = store.GetExerciseState("primitives01")
	if exState.Attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", exState.Attempts)
	}
	if exState.Status != state.StatusPassed {
		t.Errorf("Expected passed, got %s", exState.Status)
	}
	if exState.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}
}

func TestState_ThreadSafety(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, ".terralings", "state.json")

	store, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_ = store.RecordAttempt("primitives01", "01_primitives", idx%2 == 0)
			_ = store.RecordHint("primitives01", "01_primitives", 1)
			_ = store.GetAnalytics(manifest.GetManifest())
		}(i)
	}
	wg.Wait()
}

func TestState_AnalyticsCalculations(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, ".terralings", "state.json")

	store, _ := state.NewStore(statePath)
	_ = store.RecordAttempt("primitives01", "01_primitives", true)
	_ = store.RecordAttempt("primitives02", "01_primitives", false)
	_ = store.RecordHint("primitives02", "01_primitives", 2)

	m := manifest.GetManifest()
	analytics := store.GetAnalytics(m)

	if analytics.CompletedCount != 1 {
		t.Errorf("Expected 1 completed, got %d", analytics.CompletedCount)
	}
	if analytics.TotalAttempts != 2 {
		t.Errorf("Expected 2 total attempts, got %d", analytics.TotalAttempts)
	}
	if analytics.TotalHintsViewed != 2 {
		t.Errorf("Expected 2 total hints viewed, got %d", analytics.TotalHintsViewed)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/state_test.go`  
Expected: FAIL (package `github.com/dnf0/terralings/internal/state` not found).

- [ ] **Step 3: Implement `internal/state/state.go` and `internal/state/analytics.go`**

Implement thread-safe JSON persistence with atomic file writes, automatic directory creation, and analytics summary calculations.

- [ ] **Step 4: Run tests and verify they pass**

Run: `go test -v -race ./test/state_test.go`  
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/state/ test/state_test.go
git commit -m "feat(state): implement thread-safe progress persistence and analytics store"
```

---

### Task 2: CLI Stats Command & Watcher State Integration

**Files:**
- Modify: `internal/ui/ui.go`
- Modify: `internal/watcher/watcher.go`
- Modify: `cmd/terralings/main.go`
- Modify: `test/cli_test.go`

- [ ] **Step 1: Write CLI stats tests in `test/cli_test.go`**

Add `TestCLI_StatsCommand` to test output formatting when state is empty and when exercises have recorded progress.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v -run TestCLI_StatsCommand ./test/cli_test.go`  
Expected: FAIL (`stats` command not recognized).

- [ ] **Step 3: Implement UI formatter and Cobra stats command**

- Add `FormatAnalytics(summary state.AnalyticsSummary) string` to `internal/ui/ui.go`.
- Add `stats` subcommand to `cmd/terralings/main.go` with `--state` persistent flag.
- Hook `state.Store` updates into `internal/watcher/watcher.go` and `terralings run`.

- [ ] **Step 4: Run test to verify it passes**

Run: `make build && go test -v -run TestCLI_StatsCommand ./test/cli_test.go`  
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/ui/ui.go internal/watcher/watcher.go cmd/terralings/main.go test/cli_test.go
git commit -m "feat(cli): add stats command and wire state recording to watcher and runner"
```

---

### Task 3: Compiler Diagnostic Parser Engine

**Files:**
- Create: `internal/diagnostics/diagnostics.go`
- Create: `test/diagnostics_test.go`

- [ ] **Step 1: Write diagnostic parser tests in `test/diagnostics_test.go`**

Test parsing of OpenTofu/Terraform JSON diagnostics, standard text compiler errors, and `# I AM NOT DONE` marker diagnostics.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/diagnostics_test.go`  
Expected: FAIL (`internal/diagnostics` not found).

- [ ] **Step 3: Implement `internal/diagnostics/diagnostics.go`**

Implement:
```go
package diagnostics

type Severity string
const (
    SeverityError   Severity = "error"
    SeverityWarning Severity = "warning"
    SeverityInfo    Severity = "info"
)

type Diagnostic struct {
    File      string   `json:"file"`
    Line      int      `json:"line"`
    Column    int      `json:"column"`
    EndLine   int      `json:"end_line,omitempty"`
    EndColumn int      `json:"end_column,omitempty"`
    Severity  Severity `json:"severity"`
    Summary   string   `json:"summary"`
    Detail    string   `json:"detail"`
}

func Parse(rawStderr string, ex models.Exercise) []Diagnostic
```

- [ ] **Step 4: Run tests and verify they pass**

Run: `go test -v ./test/diagnostics_test.go`  
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/diagnostics/ test/diagnostics_test.go
git commit -m "feat(diagnostics): add compiler diagnostic normalizer and error parser"
```

---

### Task 4: Language Server Protocol (LSP) Server & NDJSON Streaming

**Files:**
- Create: `internal/lsp/protocol.go`
- Create: `internal/lsp/server.go`
- Modify: `cmd/terralings/main.go`
- Modify: `internal/watcher/watcher.go`
- Create: `test/lsp_test.go`

- [ ] **Step 1: Write LSP protocol and server tests in `test/lsp_test.go`**

Test:
- Server initialization handshake.
- `textDocument/didOpen` & `textDocument/didSave` diagnostic publishing.
- `textDocument/hover` response with exercise description and progressive hint.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/lsp_test.go`  
Expected: FAIL (`internal/lsp` not found).

- [ ] **Step 3: Implement `internal/lsp` and add `lsp` command + `--json` flag**

- Implement JSON-RPC 2.0 stdio server in `internal/lsp/server.go`.
- Add `terralings lsp` subcommand.
- Add `--json` flag to `terralings watch` emitting NDJSON stream events.

- [ ] **Step 4: Run tests and verify they pass**

Run: `go test -v -race ./test/lsp_test.go`  
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/lsp/ cmd/terralings/main.go internal/watcher/watcher.go test/lsp_test.go
git commit -m "feat(lsp): implement Language Server Protocol daemon and NDJSON watch stream"
```

---

### Task 5: Interactive Full-Screen Bubble Tea Terminal Dashboard

**Files:**
- Create: `internal/tui/model.go`
- Create: `internal/tui/view.go`
- Create: `internal/tui/update.go`
- Create: `internal/tui/keymap.go`
- Modify: `cmd/terralings/main.go`
- Create: `test/tui_test.go`

- [ ] **Step 1: Write TUI unit tests in `test/tui_test.go`**

Test:
- Model initialization with manifest and state store.
- Keybinding event transitions (exercise navigation, hint toggling, search filtering).
- View rendering composition (sidebar, main pane, hint drawer, footer).

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/tui_test.go`  
Expected: FAIL (`internal/tui` not found).

- [ ] **Step 3: Implement `internal/tui` Bubble Tea dashboard**

- Use `bubbletea`, `bubbles/viewport`, `bubbles/spinner`, and `lipgloss`.
- Implement split-pane reactive layout, background file watcher channel, and hotkeys.
- Register `terralings tui` and `terralings watch --interactive` / `-i`.

- [ ] **Step 4: Run tests and verify they pass**

Run: `go test -v -race ./test/tui_test.go`  
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/tui/ cmd/terralings/main.go test/tui_test.go
git commit -m "feat(tui): implement interactive Bubble Tea full-screen terminal dashboard"
```

---

### Task 6: Documentation, Integration & Verification

**Files:**
- Modify: `README.md`
- Run: `make all` and `make test-race`

- [ ] **Step 1: Update README.md with TUI, Stats, and LSP documentation**
- [ ] **Step 2: Run complete verification suite (`make all`)**
- [ ] **Step 3: Commit, push branch, open PR, verify CI matrix, and merge into `main`**

```bash
git add README.md
git commit -m "docs: document TUI dashboard, analytics stats, and LSP editor integration"
```
