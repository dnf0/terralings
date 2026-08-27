# Terralings VS Code Extension & Shell Autocompletions Design Specification

**Date:** 2026-08-27  
**Status:** Approved  
**Target:** Terralings VS Code Companion Extension (`extensions/vscode`) & CLI Shell Autocompletions  

---

## 1. Overview & Objectives

This specification defines the architecture, components, and verification strategy for two core developer experience enhancements to Terralings:
1. **Official VS Code Extension (`extensions/vscode`)**: A dedicated companion extension providing an embedded LSP client (connecting to `terralings lsp`), an interactive 13-chapter exercise explorer tree view, progress status bar indicator, integrated terminal launchers (`watch`, `tui`), command palette actions, and an embedded 5-step VS Code interactive walkthrough.
2. **CLI Shell Autocompletions (`terralings completion`)**: Native Cobra completion generators for `bash`, `zsh`, `fish`, and `powershell` with dynamic argument completion for all 56 exercise names and 13 chapters across `run`, `hint`, `reset`, and `search`.

---

## 2. Component Specifications

### 2.1 VS Code Extension (`extensions/vscode`)

#### 2.1.1 Directory Structure
```text
extensions/vscode/
├── package.json               # Extension manifest (contributes: views, commands, walkthroughs, configs)
├── tsconfig.json              # TypeScript compiler options
├── esbuild.js                 # Production bundle configuration
├── .vscodeignore              # Exclusion rules for .vsix packaging
├── README.md                  # Marketplace and repository extension guide
├── icons/                     # Extension branding and tree icons
│   ├── terralings.svg
│   ├── check.svg
│   ├── error.svg
│   └── play.svg
├── walkthrough/               # 5-step markdown guides
│   ├── 01_welcome.md
│   ├── 02_anatomy.md
│   ├── 03_watch.md
│   ├── 04_tui_hints.md
│   └── 05_editor_lsp.md
├── src/
│   ├── extension.ts           # Extension activation & deactivation lifecycle
│   ├── lspClient.ts           # LanguageClient management (stdio connection to terralings lsp)
│   ├── treeProvider.ts        # Curriculum TreeDataProvider (13 chapters, 56 exercises)
│   ├── statusBar.ts           # Status bar progress widget
│   ├── stateWatcher.ts        # Watches .terralings/state.json for real-time tree refresh
│   └── cliRunner.ts           # Terminal & OutputChannel helpers for CLI subcommands
└── test/
    ├── suite/
    │   ├── extension.test.ts
    │   └── treeProvider.test.ts
    └── runTest.ts
```

#### 2.1.2 Extension Manifest (`package.json`)
- **Activation Events**:
  - `onLanguage:terraform`
  - `onLanguage:hcl`
  - `workspaceContains:exercises`
  - `onView:terralings-exercises`
  - `onCommand:terralings.*`
- **Contributed Configuration**:
  - `terralings.binaryPath`: Custom path to `terralings` executable (defaults to `terralings` or `./bin/terralings`).
  - `terralings.enableLsp`: Boolean flag to toggle language server features (default: `true`).
  - `terralings.autoOpenWalkthrough`: Open onboarding walkthrough on initial workspace load (default: `true`).
- **Contributed Views & View Containers**:
  - Activity Bar container: `terralings-exercises` with icon `icons/terralings.svg`.
  - View: `terralings-tree` ("Curriculum & Exercises") supporting collapsible chapters and exercise leaf nodes.
- **Contributed Commands**:
  - `terralings.openExercise`: Opens selected exercise file.
  - `terralings.runCurrent`: Runs `terralings run` on active editor file.
  - `terralings.watch`: Opens a dedicated VS Code terminal executing `terralings watch`.
  - `terralings.tui`: Opens a dedicated VS Code terminal executing `terralings tui`.
  - `terralings.hint`: Displays progressive hints for active or selected exercise.
  - `terralings.reset`: Confirms and resets selected exercise template.
  - `terralings.doctor`: Runs diagnostics into "Terralings Diagnostics" OutputChannel.
  - `terralings.tour`: Opens the interactive tour walkthrough.
- **Contributed Walkthrough (`contributes.walkthroughs`)**:
  - Walkthrough ID: `terralings.welcome`
  - Steps:
    1. Welcome & Philosophy (`walkthrough/01_welcome.md`)
    2. Anatomy of an Exercise (`walkthrough/02_anatomy.md`)
    3. Continuous Watch Mode (`walkthrough/03_watch.md`)
    4. TUI Dashboard & Hints (`walkthrough/04_tui_hints.md`)
    5. Editor LSP & Completion (`walkthrough/05_editor_lsp.md`)

#### 2.1.3 LSP Client Integration (`src/lspClient.ts`)
- Spawns `terralings lsp` using `vscode-languageclient/node`.
- Connects stdin/stdout streams.
- Monitors Terraform (`*.tf`) and HCL (`*.hcl`) documents.
- Passes workspace URI and auto-restarts on unexpected crashes.

#### 2.1.4 Tree Data Provider & Real-Time State Sync (`src/treeProvider.ts`, `src/stateWatcher.ts`)
- Reads static exercise definitions and hierarchy.
- Reads `.terralings/state.json` (or defaults to empty state).
- Decorates items with status badges:
  - Completed: `ThemeIcon("pass", new ThemeColor("testing.iconPassed"))`
  - Failed: `ThemeIcon("error", new ThemeColor("testing.iconFailed"))`
  - In-Progress: `ThemeIcon("play", new ThemeColor("testing.iconQueued"))`
  - Unattempted: `ThemeIcon("circle-outline")`
- `stateWatcher` registers a `vscode.workspace.createFileSystemWatcher("**/.terralings/state.json")` to auto-fire `onDidChangeTreeData` when exercises pass or fail.

---

### 2.2 CLI Shell Autocompletions (`cmd/terralings/main.go`)

#### 2.2.1 `completion` Subcommand
- Syntax: `terralings completion [bash|zsh|fish|powershell]`
- Generates standard shell script to `stdout`.
- Supports `--no-descriptions` flag to omit description annotations where appropriate.

#### 2.2.2 Dynamic Cobra Argument Completions
- `runCmd`: `ValidArgsFunction` dynamically returns list of all exercise names with chapter annotations (e.g. `primitives01\t01_primitives`).
- `hintCmd`: `ValidArgsFunction` dynamically returns exercise names.
- `resetCmd`: `ValidArgsFunction` dynamically returns exercise names.
- `searchCmd`: `ValidArgsFunction` dynamically returns exercise names and chapter titles.

---

## 3. Verification & Testing Strategy

### 3.1 VS Code Extension Tests
- Unit tests verifying:
  - Extension registration and command exports.
  - TreeDataProvider hierarchy construction (13 chapter roots, 56 exercise nodes).
  - State parser correctly assigning status icons based on mock `.terralings/state.json`.
  - Binary locator resolving `terralings` paths.
  - Walkthrough assets existence and non-empty markdown content.

### 3.2 Go CLI Completion Tests
- Unit / Integration tests in `test/cli_completion_test.go`:
  - `TestCLI_Completion_Bash`: Generates valid bash completion without panic.
  - `TestCLI_Completion_Zsh`: Generates valid zsh completion.
  - `TestCLI_Completion_Fish`: Generates valid fish completion.
  - `TestCLI_Completion_PowerShell`: Generates valid powershell completion.
  - `TestCLI_DynamicCompletion_Exercises`: Tests `ValidArgsFunction` returns all 56 exercises with chapter descriptions.
