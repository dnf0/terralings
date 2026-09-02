# Design Specification: Bidirectional Docs $\leftrightarrow$ Playground Linking

- **Author**: Antigravity & User
- **Date**: 2026-09-02
- **Status**: Approved
- **Scope**: Documentation (`docs/`), WebAssembly Playground (`docs/assets/playground/`, `docs/playground/`)

---

## 1. Context & Motivation

Terralings provides a dual-model learning experience:
1. **Interactive Client-Side WebAssembly Playground**: Real-time Monaco editor, Pyodide Wasm execution, hints, and test runner for 15 chapters and 68 exercises.
2. **Comprehensive Architecture Reference Guides**: 15 in-depth architectural guides featuring theory, field specifications, DAG visualizations, failure triage trees, and anti-pattern catalogs.

To match the Kubelings interactive standard, every page in the documentation must deep-link directly into the playground, and the playground must provide direct, contextual entry points back to the corresponding architectural guides across all key UI touchpoints.

---

## 2. Architecture & Linking Touchpoints

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Documentation Site                               │
│                                                                             │
│  ┌───────────────────────┐             ┌─────────────────────────────────┐  │
│  │     docs/index.md     │             │     docs/guides/XX-*.md         │  │
│  │   (Home / Overview)   │             │   (15 Chapter Reference Guides) │  │
│  └───────────┬───────────┘             └────────────────┬────────────────┘  │
│              │ [Launch Playground →]                    │ [Hero Button &    │
│              │                                          │  Matrix Rows]     │
│              ▼                                          ▼                   │
│   ┌───────────────────────────────────────────────────────────────┐         │
│   │               Interactive Wasm Playground                     │         │
│   │                 (docs/playground/index.html)                  │         │
│   │                                                               │         │
│   │  • Top Navigation: [📖 Guide: Ch XX ↗]                        │         │
│   │  • Syllabus Sidebar: [📖 Guide ↗] on every chapter heading    │         │
│   │  • Editor Toolbar: [Ch XX Guide] quick-link                   │         │
│   │  • Hints Drawer: [Open Chapter XX Guide (Title) ↗]            │         │
│   │  • Victory Modal: [📖 Read Chapter / Next Chapter Guide]      │         │
│   └───────────────────────────────────────────────────────────────┘         │
│                                      │                                      │
│                                      ▼ (Opens in New Tab)                   │
│                        Back to docs/guides/XX-*.md                          │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Detailed UI & Content Specifications

### A. WebAssembly Playground Controller (`docs/assets/playground/playground.js`)
1. **Sidebar Chapter Headings**:
   - In `renderSidebarChapters()`, render an external link pill (`📖 Guide ↗`) for every chapter header pointing to `../guides/${guideInfo.slug}/`.
   - Prevent click event bubbling so clicking the guide link does not toggle the accordion dropdown.
2. **Top Bar Header & In-Editor Toolbar**:
   - Dynamically update `#btn-chapter-guide` and `.playground-btn-guide` on active exercise switch to reflect the active chapter slug and title.
3. **Progressive Hints Drawer**:
   - In `renderHints()`, include a dedicated reference callout linking to the specific chapter guide and topic title.
4. **Victory / Success State Modal**:
   - When an exercise passes, include secondary action buttons:
     - `📖 Review Chapter Guide`
     - When passing the last exercise in a chapter: `📖 Explore Next Chapter Guide →`.

### B. Documentation Hub (`docs/index.md` & `docs/syllabus.md`)
1. **`docs/index.md`**:
   - Update badges and descriptions to reflect **15 Chapters | 68 Exercises**.
   - Add **Cloud Architecture Blueprints** group (Chapter 14: AWS Architecture, Chapter 15: GCP Architecture).
   - Ensure all CTA cards and quick links route directly to `playground/index.html`.
2. **`docs/syllabus.md`**:
   - Verify every chapter section heading links to `guides/<slug>.md`.
   - Verify every exercise table row links directly to `playground/index.html?exercise=<id>`.

### C. 15 Architectural Reference Guides (`docs/guides/*.md`)
- Standardize all 15 guides:
  - Top card: `:material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html?chapter=<N>){ .md-button .md-button--primary }`
  - Section 6 Interactive Practice Matrix: `[**⚡ Solve in Playground →**](../playground/index.html?exercise=<id>){ .md-button .md-button--primary }` for every exercise.

---

## 4. Verification & Quality Gates

1. **Deterministic Pytest Suite**: `uv run pytest -q` passes 100%.
2. **MkDocs Strict Build**: `mkdocs build --strict` exits 0 with no broken internal links.
3. **Playground Query Routing**: Navigation via `?exercise=aws01` or `?chapter=14` properly sets active exercise, expands sidebar accordion, and updates guide links.
