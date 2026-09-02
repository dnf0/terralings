# Documentation Architecture Diagrams to Native Mermaid Implementation Plan

- **Spec:** `docs/superpowers/specs/2026-09-02-docs-mermaid-diagrams-design.md`
- **Branch:** `feat/docs-mermaid-diagrams`

---

## Task 1: Enable Mermaid Custom Fence in MkDocs & Update Homepage Diagram

### Goal
Configure `mkdocs.yml` with `pymdownx.superfences` custom fence for Mermaid and convert the WebAssembly architecture diagram in `docs/index.md` to a responsive Mermaid flowchart.

### Steps
1. In `mkdocs.yml`, configure `pymdownx.superfences` with the `mermaid` custom fence.
2. In `docs/index.md`, replace the ASCII architecture box with an interactive `flowchart LR` Mermaid diagram.
3. Verify site build with `mkdocs build --strict`.
4. Commit changes: `feat(docs): enable mkdocs mermaid rendering and update homepage architecture diagram`.

---

## Task 2: Convert Core HCL & Dynamic Logic Guides (Chapters 01 - 08)

### Goal
Replace ASCII text plots in Section 1 of chapters 01 through 08 with detailed Mermaid flowcharts.

### Steps
1. Update `docs/guides/01-primitives.md` (HCL Engine Pipeline: lexer, schema, DAG, topological execution).
2. Update `docs/guides/02-variables.md` (Variable Precedence Ladder & Custom Validation Condition).
3. Update `docs/guides/03-outputs-locals.md` (Locals DAG Transformation & Sensitive Masking Pipeline).
4. Update `docs/guides/04-functions.md` (Built-in Functions Execution & Safe `try`/`can` Fallbacks).
5. Update `docs/guides/05-meta-arguments.md` (Meta-Argument Expansion Engine & Lifecycle Hooks).
6. Update `docs/guides/06-dynamic-blocks.md` (Dynamic Block AST Generation Engine).
7. Update `docs/guides/07-data-sources.md` (Data Source Lifecycle & Assertion Gates).
8. Update `docs/guides/08-modules.md` (Module Hierarchy & Provider Alias Mapping).
9. Verify with `mkdocs build --strict`.
10. Commit changes: `feat(docs): convert chapters 01-08 architecture diagrams to native mermaid`.

---

## Task 3: Convert State, Testing, Governance & Cloud Guides (Chapters 09 - 15)

### Goal
Replace ASCII text plots in Section 1 of chapters 09 through 15 with rich Mermaid flowcharts and cloud topologies.

### Steps
1. Update `docs/guides/09-state-management.md` (Declarative State Surgery with `moved`/`import` blocks).
2. Update `docs/guides/10-testing.md` (Native Testing Pipeline `.tftest.hcl` with `run` & `mock_provider`).
3. Update `docs/guides/11-production-patterns.md` (Environment Configuration Factory & Tagging Matrix).
4. Update `docs/guides/12-opentofu.md` (OpenTofu State Encryption Engine at Rest).
5. Update `docs/guides/13-governance.md` (Root Module Encapsulation & ADR-0005 Policy Ownership).
6. Update `docs/guides/14-aws-architecture.md` (Multi-AZ Resilient AWS Enterprise Blueprint).
7. Update `docs/guides/15-gcp-architecture.md` (High-Availability Google Cloud Blueprint).
8. Verify with `mkdocs build --strict`.
9. Commit changes: `feat(docs): convert chapters 09-15 architecture diagrams to native mermaid`.

---

## Task 4: Full Verification, Knowledge Graph & PR Integration

### Goal
Run all test suites, verify strict docs build, update knowledge graph, and prepare PR.

### Steps
1. Run `uv run pytest -q`.
2. Run `mkdocs build --strict`.
3. Run `go test -v -run 'Test(Manifest|List|Curriculum)' ./test` & `make check build-all`.
4. Run `uvx --from graphifyy graphify update .`.
5. Check `roborev status` and `roborev show HEAD`.
6. Present completion options via `finishing-a-development-branch`.
