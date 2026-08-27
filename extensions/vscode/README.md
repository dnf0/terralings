# Terralings for Visual Studio Code

[![Visual Studio Code](https://img.shields.io/badge/VS%20Code-v1.80.0+-blue.svg)](https://code.visualstudio.com/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

The official **Terralings** companion extension for Visual Studio Code. Learn, practice, and master Terraform and OpenTofu through hands-on exercises with real-time feedback, interactive sidebar navigation, embedded Language Server Protocol diagnostics, and terminal dashboard integration.

---

## Overview

[Terralings](https://github.com/dnf0/terralings) is an interactive, terminal-driven learning environment for Infrastructure as Code (IaC). This extension brings the full Terralings curriculum directly into your VS Code editor, combining instant feedback loops with rich editor intelligence.

---

## Features

### 1. Language Server Protocol (LSP) Daemon
Connects automatically via `stdio` to `terralings lsp` to provide real-time editor feedback:
- **Live Diagnostics**: Immediate compiler and plan validation errors as you write HCL.
- **Hover Documentation**: Hover over any exercise file to inspect learning objectives, chapter context, and multi-level hints rendered in rich markdown.
- **Code Actions & Quick Fixes**: In-editor actions to request progressive hints when stuck.

### 2. Curriculum & Exercise Explorer (TreeView)
A dedicated sidebar panel in the Activity Bar displaying the complete 13-chapter, 56-exercise curriculum:
- **Visual Status Badges**: Real-time indicators for passed (✅), failed (❌), in-progress (▶️), and not-started (⚪) exercises.
- **Chapter Progress Counters**: Displays completed vs total exercises per chapter (e.g. `5/5`).
- **One-Click Navigation**: Click any exercise in the tree to immediately open its source file.
- **Context Actions**: Run verification, request hints, or reset templates directly from the tree items.

### 3. Interactive Guided Walkthroughs
Integrated welcome tour using VS Code's native Walkthrough UI (`terralings.welcome`):
- **5-Step Onboarding**: Philosophy, exercise anatomy, continuous watch mode, TUI dashboard, and LSP capabilities.
- **Interactive Action Buttons**: Directly launch commands and explore features from within the walkthrough steps.

### 4. Terminal & TUI Integration
- **Continuous Watch Mode (`terralings.watch`)**: Opens a dedicated VS Code terminal running `terralings watch` for sub-100ms re-evaluation on file save.
- **Full-Screen TUI Dashboard (`terralings.tui`)**: Launches the split-pane terminal dashboard (`terralings tui`) directly inside VS Code's integrated terminal.

### 5. Status Bar Progress Indicator
- Displays live curriculum completion count and percentage (e.g. `$(mortar-board) Terralings: 42/56 (75%)`).
- Click the status bar item to trigger the Exercise QuickPick search and jump to any exercise.

---

## Requirements

To use this extension, ensure you have installed:

1. **Terralings CLI**:
   - Install via curl: `curl -fsSL https://raw.githubusercontent.com/dnf0/terralings/main/install.sh | bash`
   - Or install via Go: `go install github.com/dnf0/terralings/cmd/terralings@latest`
   - Or compile locally in the repository root: `make build` (creates `./bin/terralings`).
2. **OpenTofu** (>= 1.6.0) or **Terraform** (>= 1.6.0) available in your system `$PATH`.

---

## Extension Settings

This extension contributes the following configurable settings:

| Setting | Type | Default | Description |
|---|---|---|---|
| `terralings.binaryPath` | `string` | `"terralings"` | Path or executable name for the `terralings` CLI binary. Can be an absolute path or relative to workspace root. |
| `terralings.enableLsp` | `boolean` | `true` | Enable or disable the built-in Terralings Language Server client. |
| `terralings.autoOpenWalkthrough` | `boolean` | `true` | Automatically open the welcome walkthrough when opening a Terralings workspace for the first time. |

---

## Commands Reference

| Command ID | Title | Description |
|---|---|---|
| `terralings.openExercise` | `Terralings: Open Exercise` | Show a searchable QuickPick of all 56 exercises with live status badges and jump to selection. |
| `terralings.runCurrent` | `Terralings: Verify Current Exercise` | Run verification on the currently active `.tf` exercise file or selected tree item. |
| `terralings.watch` | `Terralings: Start Watch Mode` | Launch continuous watch mode in an integrated terminal. |
| `terralings.tui` | `Terralings: Open Terminal Dashboard (TUI)` | Launch the full-screen terminal UI dashboard in an integrated terminal. |
| `terralings.hint` | `Terralings: Show Hint` | Display progressive hints for the current exercise in an output channel and popup notification. |
| `terralings.reset` | `Terralings: Reset Exercise Template` | Reset an exercise back to its clean starting template after confirmation. |
| `terralings.doctor` | `Terralings: Run Environment Diagnostics (Doctor)` | Run pre-flight checks and diagnostics against the environment. |
| `terralings.tour` | `Terralings: Open Guided Tour` | Open the interactive welcome walkthrough. |
| `terralings.refreshTree` | `Terralings: Refresh Exercises` | Reload `.terralings/state.json` and refresh tree items and status bar. |

---

## Development Guide

### Prerequisites
- Node.js (>= 18.0.0)
- npm (>= 9.0.0)

### Setup & Build
```bash
# Navigate to extension directory
cd extensions/vscode

# Install dependencies
npm install

# Type-check TypeScript sources
npm run check-types

# Bundle extension using esbuild
npm run build

# Run continuous watch compilation
npm run watch
```

### Running Tests
```bash
# Run Mocha unit test suite
npm test
```

### Packaging & Installation
To create a `.vsix` extension package:
```bash
npx @vscode/vsce package
# Install into VS Code
code --install-extension terralings-vscode-0.3.0.vsix
```

---

## License
 
Apache 2.0. See [LICENSE](LICENSE) for details.
