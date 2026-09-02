# Full Browser WebAssembly (Pyodide) Learning Platform Design Specification

- **Date:** 2026-09-02
- **Status:** Approved
- **Target Release:** `v0.12.0`
- **Topic:** Full 74-Exercise In-Browser Learning Platform with State Persistence

---

## 1. Objective & Scope

Transform the Terralings documentation site into an interactive, zero-install in-browser learning platform featuring:
1. **Full Curriculum Coverage**: All 74 exercises across all 13 chapters available directly in-browser.
2. **Client-Side State Persistence**: Complete progress, per-exercise working code, completion badges, and hint states stored locally via `localStorage` (`TerralingsStorage`) with zero backend server dependencies.
3. **Interactive Split-Pane Workspace**: Modern IDE layout with a collapsible chapter syllabus sidebar (320px) on the left and Monaco Editor + Terminal Diagnostics on the right.
4. **In-Memory HCL AST & Rule Validator**: High-performance pure-Python HCL tokenizer, AST parser, and rule-verification engine running in Pyodide WebAssembly (<15ms evaluation).
5. **Data Portability**: Full JSON export and import of learner progress (`terralings-progress.json`) for seamless backup and migration across browsers/devices.
6. **Deterministic Testing**: End-to-end bundle generation script and pytest verification ensuring parity with repository curriculum sources.

---

## 2. System Architecture & Component Interactions

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│ Browser Application (docs/playground.md / site/playground/)                            │
│                                                                                        │
│  ┌───────────────────────────┬──────────────────────────────────────────────────────┐  │
│  │ 📚 Curriculum Sidebar     │ 🏗️ Code & Diagnostics Workspace                      │  │
│  │  • 13 Chapters            │  • Chapter & Exercise Title Header                   │  │
│  │  • 74 Exercises           │  • Monaco Editor (HCL/Terraform + Auto-save)         │  │
│  │  • Status Badges & Search │  • Terminal Diagnostics (Pass/Fail, Timing, Errors)  │  │
│  │  • Overall Progress Bar   │  • Action Bar: [▶ Run] [💡 Hint] [↺ Reset] [🔍 Diff] │  │
│  └─────────────┬─────────────┴──────────────────────────┬───────────────────────────┘  │
│                │                                        │                              │
│                ▼                                        ▼                              │
│  ┌───────────────────────────┐            ┌─────────────────────────────────────────┐  │
│  │ LocalStorage State Engine │            │ Web Worker (Pyodide Wasm Runtime)       │  │
│  │  • TerralingsStorage      │            │  • In-memory Python 3.12 Wasm           │  │
│  │  • Per-exercise code      │◄───────────┤  • Pure-Python HCL AST & Rule Validator │  │
│  │  • Export / Import JSON   │            │  • Evaluates manifests in <15ms         │  │
│  └───────────────────────────┘            └─────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

### Component Roles:
1. **`scripts/build_playground_bundle.py`**:
   Extracts all 74 exercises, reference solutions, hints, chapter metadata, and in-memory test rules from `exercises/`, `solutions/`, and `internal/manifest/manifest.go`, writing `docs/assets/playground/playground-bundle.json`.
2. **`docs/assets/playground/hcl_validator.py`**:
   Pure-Python HCL lexer, AST parser, marker checker (`I AM NOT DONE`, `TODO`, `___`), and rule validation engine.
3. **`docs/assets/playground/playground-worker.js`**:
   Web Worker running Pyodide v0.26+ WebAssembly runtime. Mounts `hcl_validator.py` and validates user submissions asynchronously without blocking the UI thread.
4. **`docs/assets/playground/playground.js` (`TerralingsUI` & `TerralingsStorage`)**:
   Controls Monaco Editor, syllabus tree navigation, text search filtering, debounced auto-saving, progress calculation, fullscreen mode, and diff views.
5. **`docs/assets/playground/playground.css`**:
   Responsive split-view styles with CSS Grid/Flexbox, status indicators, fullscreen mode overrides, and dark/light theme integration with MkDocs Material.

---

## 3. State Management Specification (`TerralingsStorage`)

State is persisted under the `localStorage` key `terralings_learning_state_v1`.

### Data Schema:
```typescript
interface TerralingsSavedExercise {
  status: "not_started" | "in_progress" | "completed";
  userCode: string;           // Active working code
  hintsRevealed: number;       // Number of hints viewed (0..N)
  lastEvaluatedAt?: string;   // ISO 8601 timestamp
  passedAt?: string;          // ISO 8601 timestamp
}

interface TerralingsLearningState {
  version: 1;
  lastActiveExerciseId: string;
  exercises: Record<string, TerralingsSavedExercise>;
  stats: {
    completedCount: number;
    totalCount: number;
    completionPercentage: number;
  };
}
```

### Key Operations:
- **`loadExercise(id)`**: Retrieves saved code; if none exists, populates with starter template from bundle.
- **`saveExerciseCode(id, code)`**: Debounced auto-save (300ms) on Monaco Editor content change.
- **`markCompleted(id)`**: Updates status to `completed`, records `passedAt`, recalculates overall progress stats, and triggers badge updates in the sidebar.
- **`resetExercise(id)`**: Reverts `userCode` to the clean starter template from bundle and resets `hintsRevealed`.
- **`exportProgress()`**: Generates and downloads `terralings-progress-<date>.json`.
- **`importProgress(jsonString)`**: Validates schema, merges or overwrites state, and refreshes the UI.
- **`resetAllProgress()`**: Clears storage key and re-initializes starter state after user confirmation.

---

## 4. User Interface & Layout Specification

### 4.1 Curriculum Sidebar (Left Pane - 320px)
- **Header**: Global progress bar (`X / 74 Exercises Completed • Y%`), Export JSON icon, Import JSON icon, and Reset Progress icon.
- **Search & Filter Bar**: Instant client-side search box filtering exercises by title, ID, chapter, or keyword (e.g. `dynamic`, `moved`, `validation`), plus status toggles (`All`, `Incomplete`, `Completed`).
- **Chapter Accordions**:
  - 13 collapsible chapters with chapter index, title, and completed count badge (e.g. `01. Primitives (6/6 ✓)`).
  - Exercise item rows displaying status icon (`○` Not Started, `⏳` In Progress, `✓` Completed) and exercise name.
  - Active item highlight with automatic scroll-into-view.

### 4.2 Code & Diagnostics Workspace (Right Pane)
- **Top Bar**: Chapter title breadcrumb, exercise title, difficulty badge, and `← Prev` / `Next →` navigation.
- **Action Toolbar**:
  - `▶ Run Solution (Ctrl+Enter)`: Dispatches code to Pyodide worker.
  - `💡 Reveal Hint (H)`: Progressively unhides hint tiers in a collapsible panel.
  - `↺ Reset Code`: Restores exercise starter manifest.
  - `🔍 Compare Solution`: Switches Monaco to side-by-side Diff Editor.
  - `⛶ Fullscreen (F11)`: Expands workspace to 100vw × 100vh full-bleed edge-to-edge mode.
- **Monaco Editor Container**:
  - Full-featured code editor with syntax highlighting for Terraform/HCL.
  - Automatic theme switching matching MkDocs (`vs-dark` in slate mode, `vs` in default mode).
- **Terminal Output Pane**:
  - ANSI-styled output box rendering green passes, red assertion errors with line highlights, and execution duration in milliseconds.
  - On pass, displays celebration message and auto-advance button.

---

## 5. In-Memory HCL Validation Engine

The in-memory validation engine (`hcl_validator.py`) runs in Pyodide WebAssembly and performs:
1. **Marker Check**: Fails if `I AM NOT DONE`, `TODO`, `___`, or `/* ??? */` are present in uncommented code.
2. **HCL Structure Parsing**: Parses blocks, labels, attributes, and expressions.
3. **Semantic Verification**:
   - Compares parsed block structure, required attributes, variable definitions, and locals against the reference solution and rule set.
   - Evaluates string expressions, collection indexing, and function transforms where applicable.
4. **Line-Specific Error Reporting**: If a syntax or semantic error occurs, highlights the offending line with an error message in the terminal.

---

## 6. Verification & Test Plan

1. **Bundle Generator Test (`tests/test_playground_bundle.py`)**:
   - Verify `scripts/build_playground_bundle.py` bundles all 74 exercises across 13 chapters.
   - Assert all 74 exercises have non-empty starter code, reference solution, hints, and valid metadata.
   - Assert generated JSON structure adheres to bundle schema.
2. **In-Memory Validation Engine Parity**:
   - Run `hcl_validator.py` against all 74 reference solutions to guarantee 100% pass rate.
   - Verify incomplete starter templates fail validation with appropriate diagnostics.
3. **State Engine Integrity**:
   - Test localStorage initialization, auto-save debounce, progress calculations, and export/import roundtrip integrity.
4. **Site Build & Asset Integrity**:
   - Run `mkdocs build --strict` to verify all assets and markdown pages compile without warnings or broken references.
