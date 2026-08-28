# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.1] - 2026-08-28

### Fixed
- **Real-Time Sidebar & Status Bar Progress Synchronization**: Upgraded the VS Code extension's state tracking architecture to a robust 4-layer synchronization system:
  1. Direct Node `fs.watchFile` polling on `.terralings/state.json` bypassing VS Code dot-folder and gitignore file watcher exclusions.
  2. 1-second background polling heartbeat ensuring terminal and external CLI runs (`terralings watch`, `terralings tui`, `terralings run`) immediately update progress.
  3. Integrated VS Code lifecycle event triggers on document saves, active editor tab switches, and window focus transitions.
  4. Synchronous UI refresh hooks across all interactive commands (`terralings.runCurrent`, `terralings.reset`, `terralings.initExercises`, `terralings.hint`, `terralings.doctor`).
- **Comprehensive State File Discovery (`findStateJsonPath`)**: Intelligently locates `.terralings/state.json` across active editor document hierarchies, open workspace roots, parent folders, and standard candidate directories.

## [0.4.0] - 2026-08-28

### Added
- **Intelligent Multi-Depth Path Resolution (`pathUtils.ts`)**: VS Code extension seamlessly handles root directory workspaces, `exercises/` subfolders, individual chapter directories, and performs recursive parent directory traversal up to 6 levels.
- **1-Click Workspace Auto-Initialization UX**: When an exercise is opened in an uninitialized directory, VS Code prompts the user with a 1-click **"Initialize Exercises"** action that executes `terralings init` and opens the exercise immediately.
- **Command Palette & Sidebar Shortcuts**: Added `terralings.initExercises` ("Terralings: Initialize Exercises in Workspace") and `terralings.openNextExercise` ("Terralings: Open Next Exercise") with inline navigation actions.
- **JSON Output & State Awareness for `terralings list`**: Added `--json` output flag to `terralings list` for machine consumption, tool integration, and real-time completion state reflection.
- **Enhanced Core Engine Marker & Blank Detection**: Extended `CheckMarker` evaluation to check both comment markers (`# I AM NOT DONE`, `// I AM NOT DONE`, `<!-- I AM NOT DONE -->`) and unfilled blanks (`___`, `/* ??? */`, `<!-- ANSWER -->`).

### Changed
- **Workspace Fallback Safety**: Replaced unsafe root directory (`/`) fallbacks with user home directory resolution (`os.homedir()`), preventing read-only file system errors (`os error 30` / `EROFS`).

## [0.3.0] - 2026-08-28

### Added
- **Official VS Code Companion Extension (`terralings-vscode`)**: Rich IDE extension with embedded LSP diagnostics, interactive Exercise Explorer TreeView with real-time status badges, multi-step welcome walkthroughs, and status bar integration.
- **Shell Autocompletion Scripts**: Multi-shell completion generators for Bash, Zsh, Fish, and PowerShell (`terralings completion <shell>`).
- **Comprehensive Material for MkDocs Documentation Suite**: Multi-page documentation website (`mkdocs.yml`, `docs/`) with custom theme styling, instant search, navigation tabs, and syntax highlighting.
- **Dedicated Documentation Pages**: Getting Started (`docs/getting-started.md`), Onboarding Guide (`docs/onboarding-guide.md`), Curriculum Syllabus (`docs/syllabus.md`), CLI Reference (`docs/cli-reference.md`), VS Code Companion Guide (`docs/vscode-extension.md`), and Contributing Guide (`docs/contributing.md`).
- **Rich SVG Visual Assets**: High-fidelity terminal demo vector graphic (`assets/demo.svg` & `docs/assets/demo.svg`).
- **Root Documentation & Guide Overhaul**: Modernized top-level `README.md` with ecosystem links, architecture flowcharts, and 56-exercise curriculum matrix, updated `CONTRIBUTING.md`, and refreshed `CHANGELOG.md`.

### Changed
- **Curriculum-Wide TODO Overhaul**: Enriched all 56 exercises across 13 chapters with structured `# TODO (What)` and `# TODO (Why)` instructions explaining IaC mechanics and architectural design rationale.
- **Governance Validation Hardening**: Ensured deterministic evaluation for ADR-0005 policy encapsulation.
- **Unified Centralized Versioning**: Automated CI and test gates preventing version drift between CLI and VS Code extension.

## [0.2.0] - 2026-08-27

### Added
- **Interactive Full-Screen TUI Dashboard (`terralings tui` / `watch -i`)**: Split-pane Bubble Tea terminal dashboard with responsive exercise sidebar, live compiler output viewport, search modal (`/`), and expandable hint drawer (`h`).
- **Learning Analytics & Progress Persistence (`terralings stats`)**: Thread-safe persistent progress store (`internal/state`) tracking attempts, time invested, and per-chapter completion metrics with atomic `.tmp` + `fsync` + rename updates and automatic `.gitignore`.
- **Language Server Protocol (LSP) Daemon (`terralings lsp`)**: Built-in JSON-RPC 2.0 stdio LSP server providing live compiler diagnostics (`publishDiagnostics`), rich markdown hover documentation (`textDocument/hover`), and hint code actions (`textDocument/codeAction`) for Neovim, Helix, and VS Code.
- **NDJSON Event Streaming (`terralings watch --json`)**: Machine-readable newline-delimited JSON stream for tooling, CI pipelines, and IDE extensions.
- **Curriculum & Scaffolding Commands**: Embedded curriculum scaffolding (`terralings init`, `terralings reset`), fuzzy curriculum search (`terralings search`), and shell completion generation (`terralings completion`).

### Changed
- **Modernized to Ziglings / Rustlings v6 Hybrid Learning Model**: Removed legacy `# I AM NOT DONE` magic comments across all 56 exercises in favor of realistic `# TODO:` instructions and pure deterministic validation.
- **Interactive Watcher Controls**: When an exercise passes, `terralings watch` displays an interactive prompt (`[Enter / n] Next | [p] Prev | [r] Rerun | [q] Quit`) preventing unintended auto-skipping.

## [0.1.1] - 2026-08-27

### Changed
- Bumped E2E test verification to dynamically assert valid semantic versioning.

## [0.1.0] - 2026-08-26

### Added
- **Interactive Watch Mode**: Continuous filesystem watcher (`watcher.RunWatch`) using `fsnotify` with event debouncing, live UI feedback, and automatic exercise advancement.
- **Dual Engine Auto-Detection**: Binary discovery system (`detector.DetectBinary`) supporting OpenTofu (`tofu`) and Terraform (`terraform`) with `--bin` and `TERRALINGS_BIN` overrides.
- **Execution & Isolation Runner**: Sandboxed workspace runner (`runner.Runner`) supporting `validate`, `plan`, and `test` execution modes with shared provider plugin caching for fast sub-100ms evaluation.
- **Comprehensive 13-Chapter Curriculum**: 56 progressive exercises and reference solutions covering:
  - Chapter 01: HCL Foundations & Core Primitives (`primitives01`–`primitives06`)
  - Chapter 02: Input Variables, Types & Validations (`variables01`–`variables05`)
  - Chapter 03: Outputs, Locals & Expressions (`outputs01`, `locals01`, `expr01`, `expr02`)
  - Chapter 04: Built-in Functions & Collections (`func01`–`func05`)
  - Chapter 05: Meta-Arguments & Resource Scaling (`meta01`–`meta05`)
  - Chapter 06: Dynamic Blocks & Advanced HCL (`dynamic01`–`dynamic04`)
  - Chapter 07: Data Sources & State Querying (`data01`–`data04`)
  - Chapter 08: Modular Infrastructure Architecture (`module01`–`module05`)
  - Chapter 09: State Refactoring & Surgery (`state01`–`state04`)
  - Chapter 10: Native Unit & Integration Testing (`test01`–`test04`)
  - Chapter 11: Production Infrastructure Patterns (`pattern01`–`pattern04`)
  - Chapter 12: OpenTofu Advanced Features (`tofu01`–`tofu03`)
  - Chapter 13: Governance & Security Policies (`gov01`–`gov03`)
- **CLI Commands**: Rich Cobra CLI (`watch`, `run`, `hint`, `list`, `verify`, `version`) with progressive hints and chapter progress tracking.
- **Terminal UI**: Stylized terminal formatting (`ui`) powered by Lip Gloss, including exercise banners, progress bars, and bordered hint cards.
- **End-to-End Test Suite**: Complete E2E tests (`test/e2e_test.go`) validating all 56 exercises, 56 reference solutions, and CLI subprocess workflows.
- **CI/CD Automation**: GitHub Actions workflow testing matrix across Ubuntu and macOS on both OpenTofu and Terraform engines.
