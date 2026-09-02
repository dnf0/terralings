# Full Browser WebAssembly (Pyodide) Learning Platform Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a zero-install, 100% in-browser interactive learning platform for Terralings featuring all 74 exercises across 13 chapters, powered by Pyodide WebAssembly, in-memory HCL AST validation, client-side `localStorage` state persistence, Monaco Editor, and responsive split-pane syllabus UI.

**Architecture:** A Python bundle generator compiles all 74 exercises, starter templates, reference solutions, hints, and metadata into a static JSON asset. In the browser, a dedicated Web Worker running Pyodide v0.26+ WebAssembly executes a pure-Python HCL AST parser and semantic rule validator to verify learner code in <15ms. The UI layer coordinates Monaco Editor, debounced state persistence, and ANSI diagnostics.

**Tech Stack:** Python 3.12, Pyodide v0.26+ WebAssembly, Monaco Editor (AMD Loader), Vanilla JavaScript (ES6 Modules/Classes), CSS Grid/Flexbox, MkDocs Material.

## Global Constraints
- Target 100% test coverage across all 74 exercises and 13 chapters in the bundle.
- In-memory HCL validation execution time < 25ms per exercise in WebAssembly.
- State persistence strictly client-side via `localStorage` under `terralings_learning_state_v1` with zero server telemetry.
- Support full dark/light theme switching matching MkDocs Material scheme.
- Include responsive 100vw × 100vh fullscreen mode toggle (`F11` and UI button).

---

### Task 1: In-Memory HCL AST Tokenizer, Parser & Rule Validator (`hcl_validator.py`)

**Files:**
- Create: `docs/assets/playground/hcl_validator.py`
- Test: `tests/test_hcl_validator.py`

**Interfaces:**
- Consumes: User HCL code string, exercise name/ID.
- Produces: `validate_exercise(code: str, exercise_id: str, rules: dict) -> dict` returning `{"passed": bool, "error": str | None, "output": str, "duration_ms": float, "line": int | None}`.

- [ ] **Step 1: Write the failing unit tests for HCL tokenizer, parser, and marker detection**

```python
# tests/test_hcl_validator.py
import pytest
from docs.assets.playground.hcl_validator import parse_hcl, check_markers, validate_exercise

def test_check_markers_detects_not_done_and_todo():
    code_with_not_done = '# I AM NOT DONE\nterraform {\n  required_version = ">= 1.6.0"\n}'
    assert check_markers(code_with_not_done) is True

    code_with_todo = 'terraform {\n  # TODO: specify version\n  required_version = "___"\n}'
    assert check_markers(code_with_todo) is True

    clean_code = 'terraform {\n  required_version = ">= 1.6.0"\n}'
    assert check_markers(clean_code) is False

def test_parse_hcl_extracts_blocks_and_attributes():
    hcl = '''
    terraform {
      required_version = ">= 1.6.0"
      required_providers {
        local = {
          source  = "hashicorp/local"
          version = "~> 2.0"
        }
      }
    }
    variable "environment" {
      type    = string
      default = "dev"
    }
    '''
    ast = parse_hcl(hcl)
    assert "terraform" in ast
    assert ast["terraform"][0]["required_version"] == ">= 1.6.0"
    assert "variable" in ast
    assert ast["variable"][0]["_label"] == "environment"
    assert ast["variable"][0]["default"] == "dev"

def test_validate_exercise_primitives01_success_and_failure():
    starter_fail = 'terraform {\n  required_version = ">= 1.6.0"\n  # TODO: complete\n}'
    res_fail = validate_exercise(starter_fail, "primitives01", {})
    assert res_fail["passed"] is False

    valid_sol = '''
    terraform {
      required_version = ">= 1.6.0"
      required_providers {
        local = {
          source  = "hashicorp/local"
          version = "~> 2.0"
        }
      }
    }
    '''
    res_pass = validate_exercise(valid_sol, "primitives01", {})
    assert res_pass["passed"] is True
```

- [ ] **Step 2: Run test to verify it fails**

Run: `pytest tests/test_hcl_validator.py -v`
Expected: FAIL with `ModuleNotFoundError: No module named 'docs.assets.playground.hcl_validator'`

- [ ] **Step 3: Implement `docs/assets/playground/hcl_validator.py`**

Implement pure-Python HCL lexer, block parser, marker detector, and semantic validator.

- [ ] **Step 4: Run test to verify it passes**

Run: `pytest tests/test_hcl_validator.py -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add docs/assets/playground/hcl_validator.py tests/test_hcl_validator.py
PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "feat(playground): add in-memory pure-python hcl ast and rule validator"
```

---

### Task 2: Build-Time Bundle Generator (`build_playground_bundle.py`)

**Files:**
- Create: `scripts/build_playground_bundle.py`
- Output: `docs/assets/playground/playground-bundle.json`
- Test: `tests/test_playground_bundle.py`

**Interfaces:**
- Consumes: `exercises/`, `solutions/`, `internal/manifest/manifest.go`
- Produces: `docs/assets/playground/playground-bundle.json` containing `chapters`, `exercises`, `hcl_validator_code`, and `stats`.

- [ ] **Step 1: Write unit tests for bundle generation**

```python
# tests/test_playground_bundle.py
import json
import os
from scripts.build_playground_bundle import generate_bundle

def test_generate_bundle_includes_all_74_exercises():
    bundle = generate_bundle()
    assert len(bundle["chapters"]) == 13
    assert bundle["stats"]["totalExercises"] == 74
    assert len(bundle["exercises"]) == 74
    
    # Verify sample exercise
    ex = bundle["exercises"]["primitives01"]
    assert ex["name"] == "primitives01"
    assert "required_version" in ex["solutionCode"]
    assert len(ex["hints"]) >= 2
    assert ex["starterCode"] != ""
```

- [ ] **Step 2: Run test to verify it fails**

Run: `pytest tests/test_playground_bundle.py -v`
Expected: FAIL with `ModuleNotFoundError: No module named 'scripts.build_playground_bundle'`

- [ ] **Step 3: Implement `scripts/build_playground_bundle.py`**

Parse Go manifest or file tree to extract all 13 chapters, 74 exercises, starter files, solutions, hints, and serialize to `docs/assets/playground/playground-bundle.json`.

- [ ] **Step 4: Execute generator and verify output**

Run: `python3 scripts/build_playground_bundle.py && pytest tests/test_playground_bundle.py -v`
Expected: PASS and bundle generated (~350KB).

- [ ] **Step 5: Commit**

```bash
git add scripts/build_playground_bundle.py docs/assets/playground/playground-bundle.json tests/test_playground_bundle.py
PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "feat(playground): add build-time curriculum bundle generator"
```

---

### Task 3: In-Memory Validation Engine Parity for All 74 Solutions

**Files:**
- Modify: `tests/test_playground_bundle.py`
- Modify: `docs/assets/playground/hcl_validator.py` (if needed for edge cases)

**Interfaces:**
- Validates all 74 reference solutions against `hcl_validator.validate_exercise`.

- [ ] **Step 1: Add parameterized test for all 74 reference solutions**

```python
# in tests/test_playground_bundle.py
from docs.assets.playground.hcl_validator import validate_exercise

def test_all_74_reference_solutions_pass():
    bundle = generate_bundle()
    failures = []
    for ex_id, ex in bundle["exercises"].items():
        res = validate_exercise(ex["solutionCode"], ex_id, ex.get("rules", {}))
        if not res["passed"]:
            failures.append(f"{ex_id}: {res.get('error')}")
    assert len(failures) == 0, f"Failed solutions: {failures}"
```

- [ ] **Step 2: Run test and fix any edge cases**

Run: `pytest tests/test_playground_bundle.py -k test_all_74_reference_solutions_pass -v`
Expected: PASS (All 74 pass).

- [ ] **Step 3: Commit**

```bash
git add tests/test_playground_bundle.py docs/assets/playground/hcl_validator.py
PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "test(playground): verify 100% pass rate across all 74 reference solutions"
```

---

### Task 4: Pyodide WebAssembly Background Worker (`playground-worker.js`)

**Files:**
- Create: `docs/assets/playground/playground-worker.js`

**Interfaces:**
- Message In: `{ type: "INIT", bundle: Object }`, `{ type: "RUN_EXERCISE", exerciseId: string, code: string }`
- Message Out: `{ type: "STATUS", stage: string, message: string }`, `{ type: "RUN_RESULT", exerciseId: string, passed: bool, error: string, output: string, durationMs: number, line: number }`

- [ ] **Step 1: Implement `docs/assets/playground/playground-worker.js`**

Implement Web Worker importing `pyodide.js`, mounting `/lib/terralings/hcl_validator.py`, and handling asynchronous evaluation requests.

- [ ] **Step 2: Commit**

```bash
git add docs/assets/playground/playground-worker.js
PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "feat(playground): add pyodide webassembly worker runtime"
```

---

### Task 5: Client-Side State Persistence Engine (`TerralingsStorage`)

**Files:**
- Create: `docs/assets/playground/playground.js` (State portion)

**Interfaces:**
- `TerralingsStorage`:
  - `loadState()`
  - `saveExerciseCode(id, code)`
  - `markCompleted(id)`
  - `resetExercise(id)`
  - `exportProgress()`
  - `importProgress(jsonStr)`
  - `resetAllProgress()`

- [ ] **Step 1: Implement `TerralingsStorage` class in `docs/assets/playground/playground.js`**

Implement robust `localStorage` state management, debounced auto-saving, JSON export/import, and migration handlers.

- [ ] **Step 2: Commit**

```bash
git add docs/assets/playground/playground.js
PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "feat(playground): implement localStorage state persistence engine"
```

---

### Task 6: Interactive Split-Pane Workspace & Monaco UI Controller

**Files:**
- Modify: `docs/assets/playground/playground.js` (`TerralingsUI`)
- Create: `docs/assets/playground/playground.css`

**Interfaces:**
- Syllabus sidebar with 13 collapsible chapter accordions, real-time search/filter, progress bars.
- Monaco Editor with HCL language configuration, theme syncing, diff view, progressive hint tiers, ANSI terminal output pane, and ⛶ Fullscreen mode (`F11`).

- [ ] **Step 1: Implement `TerralingsUI` controller in `docs/assets/playground/playground.js`**
- [ ] **Step 2: Implement responsive CSS styling in `docs/assets/playground/playground.css`**
- [ ] **Step 3: Commit**

```bash
git add docs/assets/playground/playground.js docs/assets/playground/playground.css
PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "feat(playground): add interactive split-pane syllabus ui and monaco workspace"
```

---

### Task 7: MkDocs Documentation & Navigation Integration

**Files:**
- Create: `docs/playground.md`
- Modify: `mkdocs.yml`
- Modify: `docs/index.md`

**Interfaces:**
- Adds "Interactive Playground" to MkDocs nav, mounts HTML container for split-pane IDE, and embeds hero CTA in `docs/index.md`.

- [ ] **Step 1: Create `docs/playground.md`**
- [ ] **Step 2: Update `mkdocs.yml` navigation and extra javascript/css assets**
- [ ] **Step 3: Update `docs/index.md` with playground banner**
- [ ] **Step 4: Commit**

```bash
git add docs/playground.md mkdocs.yml docs/index.md
PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "docs: add webassembly playground to documentation site"
```

---

### Task 8: End-to-End Build & Verification

**Files:**
- All touched files

- [ ] **Step 1: Run complete test suite**

Run: `pytest -q`
Expected: 100% tests pass.

- [ ] **Step 2: Re-generate playground bundle to verify fresh synchronization**

Run: `python3 scripts/build_playground_bundle.py`
Expected: Clean generation with 74 exercises.

- [ ] **Step 3: Build MkDocs documentation in strict mode**

Run: `uv run mkdocs build --strict`
Expected: Build succeeds with 0 errors and 0 broken links.

- [ ] **Step 4: Final verification commit**

```bash
git add -A
PRE_COMMIT_ALLOW_NO_CONFIG=1 git commit --no-gpg-sign -m "chore: verify and synchronize browser learning platform assets"
```
