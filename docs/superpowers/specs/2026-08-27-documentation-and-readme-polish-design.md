# Terralings Documentation Suite & README Polish Design Specification

## 1. Overview & Goals

This specification defines the comprehensive overhaul of Terralings documentation, repository assets, and developer guides to match the standard established by [Kubelings](https://github.com/dnf0/kubelings).

### Core Goals
1. **Visual Appeal**: Provide a high-fidelity SVG terminal demo graphic (`assets/demo.svg`) embedded in the README hero header and documentation site.
2. **Static Documentation Site (Material for MkDocs)**: Provide a searchable, themed documentation site powered by `mkdocs-material` with Slate/Light modes, instant navigation, code copying, and tabbed content.
3. **Comprehensive Documentation Pages (`docs/`)**:
   - `docs/index.md`: Architecture diagram, core philosophy, and feature matrix.
   - `docs/getting-started.md`: Installation via curl, Go, and source, plus engine detection.
   - `docs/onboarding-guide.md`: Step-by-step interactive learner's tutorial.
   - `docs/syllabus.md`: Full 13-chapter, 56-exercise breakdown with learning goals and verification modes.
   - `docs/cli-reference.md`: Reference documentation for all 14 CLI commands, flags, environment variables, and shell completions.
   - `docs/vscode-extension.md`: Official companion extension setup, features, and troubleshooting.
   - `docs/contributing.md`: Authoring guide for exercises, solutions, tests, and PR workflows.
4. **Top-Level `README.md` & `CONTRIBUTING.md` Overhaul**: Modernize with badges, hero visual, architecture diagram, command walkthroughs, curriculum table, and `*lings` ecosystem cross-links.
5. **Makefile Automation**: Add `docs-install`, `docs-serve`, and `docs-build` targets.

---

## 2. Component Architecture & Assets

### 2.1 Terminal Demo Graphic (`assets/demo.svg`)
- **Dimensions**: 840x480 SVG.
- **Window Controls**: macOS rounded buttons (close, minimize, maximize) with title `"terralings watch — 80x24"`.
- **Styling**: Dracula / Slate dark background (`#1e1e2e` / `#181825`), vibrant accents (`#89b4fa`, `#a6e3a1`, `#f38ba8`, `#f9e2af`, `#cba6f7`).
- **Content Simulated**:
  - Terralings ASCII banner with OpenTofu/Terraform branding.
  - Step progress: `Evaluating exercises/01_primitives/primitives01.tf...`
  - Realistic diagnostic error callout with line numbers and fix instructions.
  - Success celebratory indicator: `✓ primitives01 passed! Next: primitives02`.
  - Interactive hotkey navigation prompt: `[Enter / n] Next  |  [p] Prev  |  [h] Hint  |  [r] Rerun  |  [q] Quit`.
- **Location**: `assets/demo.svg` (and mirrored in `docs/assets/demo.svg`).

---

### 2.2 Material for MkDocs Configuration (`mkdocs.yml`)
- **Theme**: `material` with dark slate toggle (primary: `indigo`, accent: `cyan`).
- **Features**:
  - `navigation.instant`, `navigation.tracking`, `navigation.sections`, `navigation.top`
  - `search.suggest`, `search.highlight`
  - `content.code.copy`
- **Markdown Extensions**:
  - `admonition`, `pymdownx.details`, `pymdownx.superfences`, `pymdownx.highlight`, `pymdownx.inlinehilite`, `pymdownx.snippets`, `pymdownx.tabbed`, `tables`, `attr_list`.
- **Navigation Map**:
  - Overview: `index.md`
  - Getting Started: `getting-started.md`
  - Onboarding Guide: `onboarding-guide.md`
  - Curriculum Syllabus: `syllabus.md`
  - CLI Reference: `cli-reference.md`
  - VS Code Extension: `vscode-extension.md`
  - Contributing: `contributing.md`

---

### 2.3 Documentation Content Specifications (`docs/*.md`)

#### `docs/index.md`
- Hero banner with `assets/demo.svg`.
- Project mission: Mastering Terraform & OpenTofu through interactive, hands-on terminal exercises.
- Architecture ASCII diagram illustrating the core subsystem flow (CLI $\rightarrow$ Watcher $\rightarrow$ Manifest $\rightarrow$ Runner $\rightarrow$ Diagnostics / LSP $\rightarrow$ Engine).
- Key feature callouts (Deterministic Validation, Live Watcher, Full-Screen TUI, Environment Doctor, CLI Tour, VS Code Companion, Shell Autocompletion).

#### `docs/getting-started.md`
- System prerequisites (OpenTofu >= 1.6 or Terraform >= 1.5, Go >= 1.22 for source builds).
- Installation methods:
  - 1-Line Curl installer: `curl -sSL https://raw.githubusercontent.com/dnf0/terralings/main/install.sh | bash`
  - Go install: `go install github.com/dnf0/terralings/cmd/terralings@latest`
  - Source build: `git clone ... && make build`
- Initial verification: `terralings doctor` and `terralings tour`.
- Starting watch mode: `terralings watch`.

#### `docs/syllabus.md`
- Detailed breakdown of all 13 chapters and 56 exercises:
  1. `01_primitives` (6 exercises, validate mode): Blocks, strings, numbers, bools, lists, maps.
  2. `02_resources` (5 exercises, plan mode): Local file, random pet, directory resources, provider configs.
  3. `03_variables` (5 exercises, validate mode): Input variables, validations, type constraints, outputs.
  4. `04_expressions` (5 exercises, validate mode): For expressions, splat syntax, conditionals, string templates.
  5. `05_metaarguments` (4 exercises, plan mode): `count`, `for_each`, `depends_on`, `lifecycle` blocks.
  6. `06_modules` (4 exercises, plan mode): Local modules, child outputs, variable pass-through, composition.
  7. `07_state` (4 exercises, validate mode): Local state, backend config, data sources, outputs.
  8. `08_provisioners` (3 exercises, validate mode): `local-exec`, environment injection, execution triggers.
  9. `09_workspaces` (3 exercises, validate mode): Multi-environment workspace conditionals, naming prefixes.
  10. `10_security` (4 exercises, validate mode): Sensitive variables, ephemeral values, precondition/postcondition checks.
  11. `11_opentofu` (4 exercises, validate mode): State encryption keys, pbkdf2 ciphers, engine compatibility.
  12. `12_testing` (4 exercises, test mode): Terraform 1.6+ test blocks (`run`), command assertions, unit mocks.
  13. `13_capstone` (5 exercises, plan/validate mode): End-to-end multi-tier microservice architecture.

#### `docs/cli-reference.md`
- Comprehensive command manual covering all 14 subcommands:
  - `terralings watch` (interactive file watcher with hotkeys)
  - `terralings tui` / `dashboard` (Bubble Tea full-screen explorer)
  - `terralings tour` (5-step onboarding walkthrough)
  - `terralings doctor` (environmental probe)
  - `terralings run <exercise>` (single exercise validation)
  - `terralings hint <exercise>` (progressive hints)
  - `terralings reset <exercise>` (template restoration)
  - `terralings search <query>` (curriculum search)
  - `terralings list` (syllabus listing)
  - `terralings verify` (progress report)
  - `terralings test` (solution self-test)
  - `terralings lsp` (Language Server Protocol daemon)
  - `terralings completion <shell>` (autocompletions for bash, zsh, fish, powershell)
  - `terralings version` (build metadata)
- Configuration environment variables: `TERRALINGS_BINARY`, `NO_COLOR`, `TERRALINGS_STATE_PATH`.

#### `docs/vscode-extension.md`
- Extension overview and feature showcase (LSP daemon, Activity Bar tree view, Walkthrough, status bar, runners).
- VSIX installation instructions (`code --install-extension dist/terralings-vscode.vsix`).
- Available settings table (`terralings.binaryPath`, `terralings.enableLsp`, `terralings.autoOpenWalkthrough`).
- Full command palette reference.

#### `docs/contributing.md`
- Development setup and prerequisites.
- Step-by-step authoring guide for new exercises and reference solutions.
- Test execution (`go test -v -race ./...`, `npm test` in `extensions/vscode`).
- Coding and style standards (`gofmt`, `terraform fmt`, Conventional Commits).

---

### 2.4 Top-Level `README.md` & `CONTRIBUTING.md`
- Badges: GitHub Actions CI, Apache-2.0 License, Go 1.22+, OpenTofu / Terraform compatible, VS Code Extension.
- Hero Demo: `assets/demo.svg` terminal graphic.
- Pedagogical Philosophy and Architecture ASCII diagram.
- Quickstart guide (curl script, Go install, Git clone).
- Comprehensive command walkthroughs.
- VS Code companion section.
- Full 13-Chapter Curriculum matrix.
- `*lings` ecosystem section with cross-links to [Kubelings](https://github.com/dnf0/kubelings), [Spanglings](https://github.com/dnf0/spanglings), and [Raylings](https://github.com/dnf0/raylings).

---

### 2.5 Makefile Automation
- `docs-install`: Installs `mkdocs-material` using `uv` (or `pip`).
- `docs-serve`: Spawns `mkdocs serve` with live reload at `http://127.0.0.1:8000`.
- `docs-build`: Compiles static documentation site into `site/`.
