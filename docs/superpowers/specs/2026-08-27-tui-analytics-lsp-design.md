# Terralings Milestone 2: Interactive TUI, Analytics Engine & Editor LSP Integration

**Date**: 2026-08-27  
**Status**: Approved  
**Target Release**: `v0.2.0`  

---

## 1. Executive Summary

This document specifies the architecture and design for three major developer experience enhancements to Terralings:
1. **Progress Persistence & Learning Analytics Engine (`internal/state` & `terralings stats`)**: Stores local progress, attempt frequencies, time tracking, and hint metrics in `.terralings/state.json` and renders an analytics dashboard.
2. **Interactive Full-Screen TUI (`internal/tui` & `terralings tui` / `terralings watch -i`)**: A Bubble Tea terminal user interface featuring a split-pane layout with chapter/exercise navigation, live validator output, dynamic progress bars, and progressive hint drawers.
3. **Diagnostic Parsing & Language Server Protocol Integration (`internal/diagnostics`, `internal/lsp`, `terralings watch --json`, `terralings lsp`)**: Normalizes OpenTofu/Terraform compiler diagnostics and serves real-time inline diagnostics, hovers, and code actions to editors (VS Code, Neovim, Helix).

---

## 2. Architecture & Component Boundaries

```
                         ┌───────────────────────────────┐
                         │      cmd/terralings/main.go   │
                         └──────────────┬────────────────┘
                                        │
         ┌──────────────────────────────┼─────────────────────────────┐
         │                              │                             │
         ▼                              ▼                             ▼
┌──────────────────┐           ┌──────────────────┐          ┌──────────────────┐
│  internal/state  │           │   internal/tui   │          │   internal/lsp   │
│  State & Metrics │           │  Bubble Tea TUI  │          │  JSON-RPC Server │
└────────┬─────────┘           └────────┬─────────┘          └────────┬─────────┘
         │                              │                             │
         │                              ▼                             ▼
         │                     ┌──────────────────┐          ┌──────────────────┐
         │                     │ internal/watcher │          │ internal/diagnos │
         │                     └────────┬─────────┘          └────────┬─────────┘
         │                              │                             │
         └──────────────────────────────┼─────────────────────────────┘
                                        │
                                        ▼
                               ┌──────────────────┐
                               │ internal/runner  │
                               └────────┬─────────┘
                                        │
                                        ▼
                               ┌──────────────────┐
                               │  tofu/terraform  │
                               └──────────────────┘
```

---

## 3. Detailed Specifications

### 3.1 Progress Persistence & Analytics (`internal/state`)

#### State File Location & Concurrency
- Primary Path: `.terralings/state.json` inside the active workspace directory.
- Global Flag: `--state <path>` or `TERRALINGS_STATE_PATH` environment variable.
- Thread Safety: All read and write operations are synchronized using `sync.RWMutex`.
- Automatic Gitignore: When initializing a new state file in a directory that has a `.git/` directory or `.gitignore`, Terralings ensures `.terralings/` is ignored.

#### JSON Data Schema
```json
{
  "version": "1.0",
  "created_at": "2026-08-27T08:00:00Z",
  "last_active_at": "2026-08-27T09:20:00Z",
  "total_time_spent_seconds": 1845,
  "exercises": {
    "primitives01": {
      "name": "primitives01",
      "chapter": "01_primitives",
      "status": "passed",
      "attempts": 2,
      "hints_viewed": 1,
      "first_attempt_at": "2026-08-27T08:05:00Z",
      "completed_at": "2026-08-27T08:08:12Z",
      "time_spent_seconds": 192
    }
  }
}
```

#### Go Interfaces & Types
```go
package state

type ExerciseState struct {
    Name             string     `json:"name"`
    Chapter          string     `json:"chapter"`
    Status           string     `json:"status"` // "not_started" | "in_progress" | "passed"
    Attempts         int        `json:"attempts"`
    HintsViewed      int        `json:"hints_viewed"`
    FirstAttemptAt   time.Time  `json:"first_attempt_at"`
    CompletedAt      *time.Time `json:"completed_at,omitempty"`
    TimeSpentSeconds int        `json:"time_spent_seconds"`
}

type Store struct {
    mu       sync.RWMutex
    filePath string
    data     StateData
}

type AnalyticsSummary struct {
    TotalExercises    int
    CompletedCount    int
    InProgressCount   int
    TotalAttempts     int
    TotalHintsViewed  int
    TotalTimeSpent    time.Duration
    ChapterSummaries  []ChapterSummary
}
```

#### CLI Command: `terralings stats`
Renders an overview of user performance and chapter metrics using Lip Gloss styling:
- Overall completion percentage bar.
- Average attempts per exercise and hint reliance index.
- Chapter-by-chapter status breakdown.

---

### 3.2 Interactive Bubble Tea Terminal Dashboard (`internal/tui`)

#### UI Framework
Built using Charm's:
- `github.com/charmbracelet/bubbletea` (Elm architecture loop)
- `github.com/charmbracelet/bubbles/viewport` (Scrollable output viewer)
- `github.com/charmbracelet/bubbles/list` / tree (Curriculum navigation)
- `github.com/charmbracelet/lipgloss` (Adaptive terminal layout styling)

#### Interface Layout
```text
┌─ Terralings v0.2.0 ─────────────────────────────────── [██████████████████░░] 78% (44/56) ─┐
│ CURRICULUM                     │ EXERCISE: dynamic03 (Chapter 06: Dynamic Blocks)          │
│                                │                                                           │
│ ▼ Chapter 01: Primitives (6/6) │ Goal: Configure nested dynamic blocks for security rules  │
│   ✓ primitives01               │ File: exercises/06_dynamic_blocks/dynamic03.tf            │
│   ✓ primitives02               │                                                           │
│   ...                          │ ─── Compiler Output ────────────────────────────────────  │
│ ▼ Chapter 06: Dynamic (2/4)    │ ❌ Error: Missing required attribute "content"             │
│   ✓ dynamic01                  │                                                           │
│   ✓ dynamic02                  │   on dynamic03.tf line 18, in dynamic "ingress":          │
│   ▶ dynamic03 (current)        │   18:   dynamic "ingress" {                               │
│   ⬜ dynamic04                 │                                                           │
│                                │ ─── Hints (Press 'h' to cycle) ─────────────────────────  │
│                                │ 💡 Hint 1/2: Dynamic blocks require a nested content { }  │
├────────────────────────────────┴───────────────────────────────────────────────────────────┤
│ [↑/↓] Navigate  [Tab] Switch Pane  [Enter] Run  [h] Hint  [r] Reset  [/] Search  [q] Quit     │
└────────────────────────────────────────────────────────────────────────────────────────────┘
```

#### Event & State Handling
- `FileChangedMsg`: Dispatched by background `fsnotify` watcher on `.tf`/`.hcl` modification.
- `RunCompletedMsg`: Returns execution result from `runner.Run()` to update viewport without UI locking.
- `KeyMsg`: Dispatches interactive keyboard hotkeys (`j`/`k`, `Tab`, `h`, `r`, `/`, `q`).

---

### 3.3 Diagnostic Parser & LSP Server (`internal/diagnostics`, `internal/lsp`)

#### Diagnostic Normalization (`internal/diagnostics`)
```go
type Diagnostic struct {
    File      string `json:"file"`
    Line      int    `json:"line"`
    Column    int    `json:"column"`
    EndLine   int    `json:"end_line,omitempty"`
    EndColumn int    `json:"end_column,omitempty"`
    Severity  string `json:"severity"` // "error" | "warning" | "information"
    Summary   string `json:"summary"`
    Detail    string `json:"detail"`
}

func ParseDiagnostics(rawOutput string, ex models.Exercise) []Diagnostic
```

#### Streaming Watch Mode: `terralings watch --json`
Emits NDJSON stream of evaluation lifecycle events:
- `exercise_start`: Notifies when evaluation begins.
- `exercise_result`: Full result object with structured diagnostic array.

#### Language Server Protocol (`terralings lsp`)
Implements standard JSON-RPC 2.0 stdio server:
- `initialize`: Returns server capabilities (textDocumentSync, hoverProvider, codeActionProvider).
- `textDocument/didOpen` & `textDocument/didSave`: Compiles the active exercise via `runner.Runner` and dispatches `textDocument/publishDiagnostics` with precise line numbers for `# I AM NOT DONE` markers and compiler errors.
- `textDocument/hover`: Displays markdown card with exercise goal, description, and progressive hints.
- `textDocument/codeAction`: Provides quick fixes (e.g. remove marker, view hint).

---

## 4. Testing & Verification Plan

1. **State Persistence Tests (`test/state_test.go`)**:
   - Verify creation, read/write concurrency safety under 50 simultaneous goroutines.
   - Verify state calculations, persistence across restart, and corrupt JSON recovery.
2. **TUI Tests (`test/tui_test.go`)**:
   - Test model initialization, keybinding dispatch, viewport scrolling, and search modal filtering.
3. **Diagnostic Parser Tests (`test/diagnostics_test.go`)**:
   - Verify parsing of OpenTofu/Terraform standard errors, JSON errors, and marker detection.
4. **LSP Server Tests (`test/lsp_test.go`)**:
   - Test JSON-RPC handshake, `textDocument/didSave` diagnostic publication, and hover documentation.
5. **CLI End-to-End Tests (`test/cli_test.go`)**:
   - Test `terralings stats`, `terralings watch --json`, `terralings lsp --help`.
6. **Full Matrix Verification**:
   - `make all` & `make test-race` on macOS and Linux runners with OpenTofu and Terraform.
