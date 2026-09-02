# Bidirectional Docs $\leftrightarrow$ Playground Linking Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Choose an execution mode:
> 1. `superpowers:subagent-driven-development` (recommended for multi-agent reviews, backed by `SKILL.state` / `.agent-state/state.json`)
> 2. `agent-rules:stateful-execution` (SKILL.state) (recommended for deterministic single-agent linear execution)
> 3. `superpowers:executing-plans` (batch execution with manual checkpoints)
> Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement comprehensive bidirectional deep-linking between all 15 documentation guides / syllabus and the interactive WebAssembly playground.

**Architecture:** Enrich the client-side JavaScript UI controller in `playground.js` to provide contextual guide links in sidebar chapter headers, the editor toolbar, hints drawer, and victory completion modal; synchronize `docs/index.md` and `docs/syllabus.md` with deep links to the Wasm playground.

**Tech Stack:** JavaScript (ES6), HTML5/CSS3, Monaco Editor, Pyodide Wasm, Markdown/MkDocs Material, Python 3.12 (pytest).

## Global Constraints
- Zero-breakage policy for existing localStorage state format (`terralings_learning_state_v1`).
- All links from playground to guides must open in a new tab (`target="_blank" rel="noopener noreferrer"`).
- All MkDocs links must build cleanly under `mkdocs build --strict`.
- Pytest test suite `uv run pytest -q` must remain 100% green.

---

### Task 1: WebAssembly Playground UI Controller & Link Enhancements

**Files:**
- Modify: `docs/assets/playground/playground.js:590-640, 720-755, 870-940, 1020-1080`
- Modify: `docs/assets/playground/playground.css:200-280`
- Modify: `docs/playground/index.html:140-170`

**Interfaces:**
- Consumes: `CHAPTER_GUIDES` map (chapters 1-15).
- Produces: Contextual chapter guide pills in sidebar accordions, editor header, hints drawer, and victory modal.

- [ ] **Step 1: Write / Update test for playground bundle and links**
- [ ] **Step 2: Add guide link pill to sidebar chapter headers in `playground.js`**
- [ ] **Step 3: Enhance hints drawer and victory modal with guide buttons**
- [ ] **Step 4: Add CSS styles for `.playground-chapter-guide-pill` and victory action buttons**
- [ ] **Step 5: Verify in browser / node / test suite**
- [ ] **Step 6: Commit changes**

---

### Task 2: Documentation Hub & Guide Link Standardization

**Files:**
- Modify: `docs/index.md:1-105`
- Modify: `docs/syllabus.md:1-120`
- Check: `docs/guides/01-primitives.md` through `15-gcp-architecture.md`

**Interfaces:**
- Produces: Complete bidirectional links from homepage, syllabus, and all 15 reference guides to `playground/index.html?exercise=<id>` and `playground/index.html?chapter=<N>`.

- [ ] **Step 1: Update `docs/index.md` with 15 chapters, 68 exercises, and Cloud Blueprints section**
- [ ] **Step 2: Audit all 15 `docs/guides/*.md` for hero launch button and Section 6 practice matrix**
- [ ] **Step 3: Run `mkdocs build --strict` to verify no broken links**
- [ ] **Step 4: Commit changes**

---

### Task 3: Full End-to-End Verification & Quality Gate

**Files:**
- Test: `tests/test_playground_bundle.py`
- Verify: `uv run pytest -q`
- Verify: `mkdocs build --strict`
- Verify: `make check build-all`

- [ ] **Step 1: Run pytest suite**
- [ ] **Step 2: Run MkDocs strict build**
- [ ] **Step 3: Run Go checks & build**
- [ ] **Step 4: Commit & push**
