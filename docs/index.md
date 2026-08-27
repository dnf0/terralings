# Terralings

<p align="center">
  <img src="assets/demo.svg" alt="Terralings Terminal Demo" width="800">
</p>

<p align="center">
  <strong>An interactive, terminal-driven learning environment for mastering Terraform and OpenTofu.</strong>
</p>

<p align="center">
  <a href="getting-started.md">Getting Started</a> •
  <a href="onboarding-guide.md">Onboarding Guide</a> •
  <a href="syllabus.md">Curriculum Syllabus</a> •
  <a href="cli-reference.md">CLI Reference</a> •
  <a href="vscode-extension.md">VS Code Companion</a> •
  <a href="contributing.md">Contributing</a>
</p>

---

## Mission

**Terralings** bridges the gap between passive syntax reading and real-world Infrastructure as Code (IaC) engineering. Inspired by [Rustlings](https://github.com/rust-lang/rustlings), [Ziglings](https://github.com/ziglings/exercises), and [Raylings](https://github.com/ray-project/raylings), Terralings provides an interactive, hands-on curriculum that takes you from foundational HCL syntax to complex production-grade architectures, state migrations, and enterprise governance.

Whether you are configuring your first resource block or orchestrating zero-downtime state refactors with `moved` blocks and `.tftest.hcl` suites, Terralings gives you an immediate, friction-free feedback loop right in your terminal.

---

## Pedagogical Philosophy

Terralings is built on five core educational pillars:

1. **Active Debugging over Passive Reading**: Each exercise presents realistic broken or incomplete configuration code with clear `# TODO:` instructions. You learn by identifying compilation errors, diagnosing plan discrepancies, and writing working declarative code.
2. **Instant Feedback Loop**: Powered by a sub-30ms file watcher (`fsnotify`), Terralings re-evaluates your changes immediately upon saving the file. Hotkeys allow instant navigation, progressive hints, and manual re-evaluation without restarting your session.
3. **Dual-Engine Compatibility**: Terralings natively detects and executes against either **OpenTofu** (`tofu` &ge; 1.6.0) or **Terraform** (`terraform` &ge; 1.5.0), ensuring universal applicability across open-source and enterprise toolchains.
4. **Progressive Hinting System**: When you get stuck, multi-tiered contextual hints provide gentle nudges and concept reminders before revealing targeted syntax patterns, preserving the learning challenge.
5. **Zero-Friction Validation**: Terralings eliminates cumbersome "magic comments" or manual completion markers. Your code is validated directly against canonical parser checks, execution plans, and native test runners.

---

## Architecture Overview

Terralings is engineered in Go for extreme performance, offline reliability, and zero cloud credential requirements. The complete 56-exercise curriculum is embedded directly in the standalone binary.

```
                            +-----------------------+
                            |     User Terminal     |
                            +-----------+-----------+
                                        |
                                        v
                            +-----------------------+
                            |  Terralings CLI (Go)  |
                            +-----------+-----------+
                                        |
               +------------------------+------------------------+
               |                                                 |
               v                                                 v
   +-----------------------+                         +-----------------------+
   |  File Watcher Engine  |                         |  Bubble Tea TUI & UI  |
   |       (fsnotify)      |                         | (diagnostics / tree)  |
   +-----------+-----------+                         +-----------------------+
               |
               v
   +-----------------------+
   |  Curriculum Manifest  |  (13 Chapters / 56 Exercises)
   +-----------+-----------+
               |
               v
   +-----------------------+
   |   Exercise Runner     |
   +-----------+-----------+
               |
   +-----------+-----------+
   |                       |
   v                       v
+---------------+   +------------------+
|  OpenTofu CLI |   |  Terraform CLI   |
| (tofu binary) |   | (terraform bin)  |
+---------------+   +------------------+
```

### Key Components

- **Embedded Assets (`exercises/`)**: Complete curriculum embedded via Go's `embed.FS`—extract instantly anywhere with `terralings init`.
- **Diagnostic Engine (`terralings doctor`)**: Pre-flight system probe checking binary availability, filesystem permissions, and cache health.
- **Interactive Tour (`terralings tour`)**: Built-in 5-step terminal walkthrough covering IaC fundamentals, TUI shortcuts, and LSP setup.
- **Bubble Tea TUI (`terralings tui` / `watch -i`)**: Full-screen dual-pane terminal interface with real-time compilation viewport, fuzzy search, and hints drawer.
- **Language Server Protocol Daemon (`terralings lsp`)**: JSON-RPC 2.0 language server delivering in-editor diagnostics, hover tooltips, and quick-fixes for VS Code, Neovim, and Helix.

---

## Quick Navigation

<div class="grid cards" markdown>

-   :material-rocket-launch: **[Getting Started](getting-started.md)**

    ---

    Prerequisites, installation options (1-line script, Go toolchain, source), pre-flight checks, and first exercise walkthrough.

-   :material-book-open-page-variant: **[Curriculum Syllabus](syllabus.md)**

    ---

    Explore all 13 chapters and 56 exercises covering primitives, variables, modules, testing, OpenTofu encryption, and architecture governance.

-   :material-console: **[CLI Reference](cli-reference.md)**

    ---

    Comprehensive guide to all 15 subcommands, interactive hotkeys, flags, NDJSON event streaming, and environment variables.

-   :material-microsoft-visual-studio-code: **[VS Code Companion](vscode-extension.md)**

    ---

    Install the official extension for curriculum tree views, live diagnostic badges, interactive walkthroughs, and one-click runners.

-   :material-compass: **[Onboarding Guide](onboarding-guide.md)**

    ---

    Deep-dive onboarding journey explaining pedagogical concepts, terminal workflows, analytics, and troubleshooting tips.

-   :material-account-group: **[Contributing](contributing.md)**

    ---

    Guidelines for authoring new exercises, maintaining solutions, running test suites, and submitting pull requests.

</div>

---

## Key Feature Highlights

| Feature | Description |
|---|---|
| **Zero Cloud Dependencies** | Executes purely locally using built-in providers (`local`, `random`, `archive`, `terraform_data`, `test_mock`)—no AWS/GCP/Azure credentials or cloud spend required. |
| **Sub-100ms Runner** | Shared plugin cache (`~/.terralings/plugin-cache`) eliminates redundant network downloads during provider initialization. |
| **Progress Persistence** | Automatically tracks attempts, pass/fail status, hint usage, and time invested in `.terralings/state.json`. |
| **Full-Text Curriculum Search** | Search topics, error messages, and HCL concepts across all exercises using `terralings search <term>`. |
| **Shell Autocompletions** | First-class tab completion for Bash, Zsh, Fish, and PowerShell with dynamic exercise name suggestions. |
