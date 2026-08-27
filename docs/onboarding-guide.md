# Terralings Onboarding Guide

> An interactive, terminal-driven learning journey to master Terraform and OpenTofu through hands-on practice.

Welcome to **Terralings**! This guide takes you from zero to confidently authoring, refactoring, testing, and governing production-grade Infrastructure as Code (IaC).

---

## Table of Contents

1. [Welcome & Core Philosophy](#1-welcome-core-philosophy)
2. [Step 0: Environment Health Check (`terralings doctor`)](#2-step-0-environment-health-check-terralings-doctor)
3. [Step 1: Workspace Initialization (`terralings init`)](#3-step-1-workspace-initialization-terralings-init)
4. [Step 2: The Interactive Guided Tour (`terralings tour`)](#4-step-2-the-interactive-guided-tour-terralings-tour)
5. [Step 3: Solving Your First Exercise (`primitives01`)](#5-step-3-solving-your-first-exercise-primitives01)
6. [Step 4: Choosing Your Learning Modality](#6-choosing-your-learning-modality)
   - [Modal A: Continuous Watch Mode (`terralings watch`)](#modal-a-continuous-watch-mode-terralings-watch)
   - [Modal B: Interactive Full-Screen TUI Dashboard (`terralings tui`)](#modal-b-interactive-full-screen-tui-dashboard-terralings-tui)
7. [Step 5: Using Progressive Hints, Search & Reset](#7-step-5-using-progressive-hints-search-reset)
   - [Progressive Hints (`terralings hint`)](#progressive-hints-terralings-hint)
   - [Curriculum Search (`terralings search`)](#curriculum-search-terralings-search)
   - [Resetting Exercises (`terralings reset`)](#resetting-exercises-terralings-reset)
8. [Step 6: Progress Tracking & Analytics (`terralings stats`)](#8-step-6-progress-tracking-analytics-terralings-stats)
9. [Step 7: Editor Integration & Language Server (`terralings lsp`)](#9-step-7-editor-integration-language-server-terralings-lsp)
   - [Visual Studio Code](#visual-studio-code)
   - [Neovim (`nvim-lspconfig`)](#neovim-nvim-lspconfig)
   - [Helix (`languages.toml`)](#helix-languagestoml)
10. [Curriculum Roadmap](#10-curriculum-roadmap)
11. [Troubleshooting & FAQ](#11-troubleshooting-faq)

---

## 1. Welcome & Core Philosophy

Terralings is inspired by [`rustlings`](https://github.com/rust-lang/rustlings), [`ziglings`](https://github.com/ziglings/exercises), and [`raylings`](https://github.com/ray-project/raylings). It provides a self-contained, offline-first feedback loop for learning HashiCorp Configuration Language (HCL), Terraform, and OpenTofu.

```
       ┌────────────────┐
       │ Read `# TODO:` │
       └───────┬────────┘
               │
       ┌───────▼────────┐
       │   Edit Code    │◄───────────┐
       └───────┬────────┘            │
               │ (Save file)         │ Fix
       ┌───────▼────────┐            │ Errors
       │  Auto-Evaluate ├─[ Failed ]─┘
       └───────┬────────┘
               │ [ Passed ]
       ┌───────▼────────┐
       │ Next Exercise! │
       └────────────────┘
```

### Core Design Principles

- **Zero Cloud Credentials or Costs**: All exercises run entirely locally using lightweight built-in providers (`local`, `random`, `archive`, `terraform_data`, `test_mock`) or fast in-memory execution. You never need AWS, GCP, or Azure accounts to complete the curriculum.
- **Sandboxed Runner**: Each exercise runs in an isolated workspace with shared provider plugin caching (`~/.terralings/plugin-cache`), executing evaluation subroutines in under 100 milliseconds.
- **Zero Magic Comment Friction**: Unlike older educational tools that require removing `// I AM NOT DONE` comments, Terralings validates your code deterministically against real compiler syntax, validation checks, execution plans, and `.tftest.hcl` assertions.
- **Dual Engine Compatibility**: Works seamlessly with either **OpenTofu** (`tofu`) or **Terraform** (`terraform`) (version >= 1.6.0).

---

## 2. Step 0: Environment Health Check (`terralings doctor`)

Before starting the curriculum, run `terralings doctor` to verify that your system is configured correctly:

```bash
terralings doctor
```

### Diagnostic Checks Performed

`terralings doctor` evaluates 5 critical pre-flight criteria:

| Check | What is Validated | Resolution if Warning/Failure |
|---|---|---|
| **IaC Engine Binary** | Discovers `tofu` or `terraform` in `$PATH` or via `--bin` / `$TERRALINGS_BIN` | Install OpenTofu (`brew install opentofu` or official curl script) or Terraform |
| **Curriculum Scaffold** | Verifies the presence and integrity of the `exercises/` folder | Run `terralings init` in your project folder |
| **Provider Plugin Cache** | Verifies read/write access to `~/.terralings/plugin-cache` | Ensure write permissions on `~/.terralings` or set `TERRALINGS_PLUGIN_CACHE_DIR` |
| **Progress Persistence Store** | Validates `.terralings/state.json` permissions and JSON schema | Check write permissions in `.terralings/` |
| **Git Ignore Integration** | Ensures `.terralings/` is ignored if working inside a Git repository | Add `.terralings/` to your `.gitignore` |

### Sample Output

```text
🩺 Terralings Doctor Diagnostic Report
────────────────────────────────────────────────────────────

 ✓ IaC Engine Binary
   Found opentofu at /usr/local/bin/tofu (OpenTofu v1.8.0)

 ✓ Curriculum Scaffold
   Exercises directory present (56 configuration files found)

 ✓ Provider Plugin Cache
   Plugin cache directory ready at /Users/username/.terralings/plugin-cache

 ✓ Git Ignore Integration
   .terralings directory is properly git-ignored.

 ✓ Progress Persistence Store
   State store healthy at .terralings/state.json (0 completed, 0 attempts)

────────────────────────────────────────────────────────────
 All diagnostics passed! Your environment is 100% ready for Terralings.
```

### Machine-Readable Diagnostics

To integrate environment checks into scripts or CI pipelines, pass `--json`:

```bash
terralings doctor --json
```

---

## 3. Step 1: Workspace Initialization (`terralings init`)

Terralings embeds all 56 exercises directly into the binary. You do not need to clone a git repository to start practicing.

Create a directory for your practice exercises and initialize:

```bash
mkdir my-terralings && cd my-terralings
terralings init
```

This creates the `exercises/` directory with 13 progressive chapters:

```text
exercises/
├── 01_primitives/
│   ├── primitives01.tf
│   ├── primitives02.tf
│   └── ...
├── 02_variables/
├── 03_outputs_locals/
├── 04_functions/
├── 05_meta_arguments/
├── 06_dynamic_blocks/
├── 07_data_sources/
├── 08_modules/
├── 09_state_refactoring/
├── 10_testing/
├── 11_patterns/
├── 12_opentofu/
└── 13_governance/
```

> [!TIP]
> If you ever want to re-extract fresh exercise files over existing ones, use the force flag:
> ```bash
> terralings init --force
> ```

---

## 4. Step 2: The Interactive Guided Tour (`terralings tour`)

To get an immediate 2-minute visual walkthrough of Terralings, run:

```bash
terralings tour
```

### The 5 Tour Steps

```
┌────────────────────────────────────────────────────────────┐
│ STEP 1 OF 5  Welcome & Core Philosophy                     │
│ Master Terraform & OpenTofu through interactive practice   │
├────────────────────────────────────────────────────────────┤
│ STEP 2 OF 5  Anatomy of an Exercise                        │
│ How exercises are structured and solved                    │
├────────────────────────────────────────────────────────────┤
│ STEP 3 OF 5  Continuous Watch & Verification               │
│ The rapid edit-save-verify feedback loop                   │
├────────────────────────────────────────────────────────────┤
│ STEP 4 OF 5  Interactive TUI, Hints & Analytics            │
│ Terminal dashboard, progressive hints, and progress stats  │
├────────────────────────────────────────────────────────────┤
│ STEP 5 OF 5  Editor Integration & LSP                      │
│ Live compiler diagnostics and hover docs in your editor    │
└────────────────────────────────────────────────────────────┘
```

### Interactive Tour Controls

While running `terralings tour`:
- Press <kbd>Enter</kbd> or <kbd>n</kbd> to advance to the next step.
- Press <kbd>p</kbd> to return to the previous step.
- Press <kbd>1</kbd> – <kbd>5</kbd> to jump directly to any step.
- Press <kbd>q</kbd> or <kbd>Ctrl+C</kbd> to exit the tour.

### Non-Interactive & JSON Tour Modes

Render all steps at once without waiting for keyboard input:
```bash
terralings tour --non-interactive
```

Render a specific step directly:
```bash
terralings tour --step 3 --non-interactive
```

Export tour content as JSON for tooling:
```bash
terralings tour --json
```

---

## 5. Step 3: Solving Your First Exercise (`primitives01`)

Let's walk through solving the very first exercise: `exercises/01_primitives/primitives01.tf`.

### 1. Open the Exercise File

Open `exercises/01_primitives/primitives01.tf` in your editor. You will see:

```terraform
# ==============================================================================
# Exercise: primitives01
# Chapter: 01_primitives (HCL Foundations & Core Primitives)
#
# Task:
# Every Terraform / OpenTofu project begins with a `terraform` configuration
# block. This block specifies the minimum required engine version and declares
# the providers needed by your configuration.
#
# In this exercise, complete the `terraform` block to require version ">= 1.6.0"
# and declare the `local` provider from "hashicorp/local" with version "~> 2.0".
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"

  required_providers {
    # TODO: Specify provider source "hashicorp/local" and version "~> 2.0"
    local = {
      source = ""
    }
  }
}
```

### 2. Run Single Exercise Verification

Verify the exercise from the CLI:

```bash
terralings run primitives01
```

Terralings evaluates the file and prints the compiler diagnostic output:

```text
✗ primitives01: Terraform Configuration Block (exercises/01_primitives/primitives01.tf)

Error: Invalid version constraint string

  on exercises/01_primitives/primitives01.tf line 16, in terraform:
  16:   required_version = "INVALID_VERSION"

The given string "INVALID_VERSION" is not a valid version constraint.
```

### 3. Apply the Fix

Edit `exercises/01_primitives/primitives01.tf` and resolve the TODO instructions:

```hcl
terraform {
  required_version = ">= 1.6.0"

  required_providers {
    local = {
      source  = "hashicorp/local"
      version = ">= 2.4.0"
    }
  }
}
```

### 4. Verify the Solution

Run `terralings run primitives01` again:

```text
✓ primitives01: Terraform Configuration Block (exercises/01_primitives/primitives01.tf)
  Configuration is valid.
```

Congratulations! You have solved your first Terralings exercise.

---

## 6. Choosing Your Learning Modality

Terralings offers two primary learning workflows to match your preference:

### Modal A: Continuous Watch Mode (`terralings watch`)

Standard continuous watch mode monitors the `exercises/` folder via filesystem notifications. Every time you save a `.tf` or `.hcl` file, Terralings immediately validates it.

```bash
terralings watch
```

```text
======================================================================
  TERRALINGS: Continuous Watch Mode
  Watching directory: exercises
  Press Ctrl+C to exit
======================================================================

✓ primitives01 passed!

[Enter / n] Next exercise (primitives02)  |  [p] Previous  |  [r] Rerun  |  [q] Quit
```

#### Interactive Watch Controls
When an exercise passes, the watcher pauses and offers interactive hotkeys:
- <kbd>Enter</kbd> or <kbd>n</kbd>: Advance to next exercise.
- <kbd>p</kbd>: Return to previous exercise.
- <kbd>r</kbd>: Re-run the current exercise.
- <kbd>q</kbd>: Quit watch mode.

---

### Modal B: Interactive Full-Screen TUI Dashboard (`terralings tui`)

For an immersive, terminal dashboard experience with split panes, real-time error viewports, collapsible hint drawers, and quick search:

```bash
terralings tui
# or
terralings watch -i
```

#### TUI Dashboard Layout

```
┌───────────────────────────────┬──────────────────────────────────────────────┐
│  TERRALINGS CURRICULUM        │  EVALUATION VIEWPORT: primitives01           │
├───────────────────────────────┼──────────────────────────────────────────────┤
│ ▼ Chapter 01: Primitives      │                                              │
│   ✓ primitives01              │ Error: Invalid version constraint string     │
│   • primitives02              │   on exercises/01_primitives/primitives01.tf │
│   · primitives03              │   required_version = "INVALID_VERSION"       │
│                               │                                              │
│ ► Chapter 02: Variables       │ ──────────────────────────────────────────── │
│ ► Chapter 03: Outputs & Locals│ 💡 Hint (1/2): Set required_version =        │
│                               │    ">= 1.6.0" inside the terraform block.    │
├───────────────────────────────┴──────────────────────────────────────────────┤
│ [Tab] Switch Pane | [Enter] Run | [h] Toggle Hints | [/] Search | [q] Quit   │
└──────────────────────────────────────────────────────────────────────────────┘
```

#### TUI Keyboard Shortcuts

| Key | Action |
|---|---|
| <kbd>↑</kbd> / <kbd>k</kbd> | Move cursor up in exercise list / search results |
| <kbd>↓</kbd> / <kbd>j</kbd> | Move cursor down in exercise list / search results |
| <kbd>Tab</kbd> | Toggle active focus between Chapter Sidebar and Viewport |
| <kbd>Enter</kbd> | Evaluate selected exercise (or select item in search) |
| <kbd>h</kbd> | Toggle progressive hint drawer / cycle hint levels |
| <kbd>r</kbd> | Reset exercise back to starter template |
| <kbd>/</kbd> | Open curriculum search modal |
| <kbd>Esc</kbd> | Close modal / cancel search |
| <kbd>q</kbd> / <kbd>Ctrl+C</kbd> | Quit TUI |

---

## 7. Step 5: Using Progressive Hints, Search & Reset

### Progressive Hints (`terralings hint`)

When you are stuck on an exercise, use progressive hints. Terralings provides tiered hints starting with conceptual nudges before revealing syntax examples:

```bash
# View the first hint
terralings hint primitives01
```

```text
╭────────────────────────────────────────────────╮
│                                                │
│ 💡 Hint (1/2) for primitives01:                │
│ Set required_version = ">= 1.6.0" inside the   │
│ terraform block.                               │
│                                                │
╰────────────────────────────────────────────────╯
```

To reveal deeper hints, pass the `--index` flag:

```bash
# View the second (more detailed) hint
terralings hint primitives01 -i 1
```

---

### Curriculum Search (`terralings search`)

Quickly find exercises covering specific HCL constructs, functions, or topics:

```bash
terralings search "for_each"
```

```text
🔍 Search Results for 'for_each' (3 matches):
  • meta02           Resource Mapping with for_each
    exercises/05_meta_arguments/meta02.tf | matched in: name, chapter, hints
  • module04         Multi-instance Module Provisioning with for_each
    exercises/08_modules/module04.tf | matched in: name, chapter, hints
  • dynamic01        Repeating Nested Blocks with dynamic
    exercises/06_dynamic_blocks/dynamic01.tf | matched in: hints
```

---

### Resetting Exercises (`terralings reset`)

If you want to start an exercise from scratch or test your memory, reset it back to its original starter template:

```bash
terralings reset primitives01
```

```text
🔄 Reset exercise 'primitives01' (exercises/01_primitives/primitives01.tf) back to original template.
```

---

## 8. Step 6: Progress Tracking & Analytics (`terralings stats`)

Terralings automatically records your learning metrics in `.terralings/state.json`. Run `terralings stats` at any time to inspect your progress:

```bash
terralings stats
```

### Sample Analytics Output

```text
📊 TERRALINGS LEARNING ANALYTICS

Overall Progress: [████████████████████░░░░░░░░░░] 67% (38/56 completed)
Time Invested:    2h 15m
Total Attempts:   94 (avg 1.6 per exercise)
Hints Consulted:  14

Chapter Breakdown:
  • 01_primitives        100% (6/6) | Avg Attempts: 1.2 | Hints: 1
  • 02_variables         100% (5/5) | Avg Attempts: 1.4 | Hints: 2
  • 03_outputs_locals    100% (4/4) | Avg Attempts: 1.0 | Hints: 0
  • 04_functions         100% (5/5) | Avg Attempts: 1.8 | Hints: 3
  • 05_meta_arguments    100% (5/5) | Avg Attempts: 1.6 | Hints: 2
  • 06_dynamic_blocks    100% (4/4) | Avg Attempts: 2.2 | Hints: 3
  • 07_data_sources       75% (3/4) | Avg Attempts: 1.5 | Hints: 1
  • 08_modules            40% (2/5) | Avg Attempts: 2.0 | Hints: 2
  • 09_state_refactoring   0% (0/4) | Avg Attempts: 0.0 | Hints: 0
  ...
```

---

## 9. Step 7: Editor Integration & Language Server (`terralings lsp`)

Terralings comes with a built-in Language Server Protocol (LSP) daemon communicating over JSON-RPC 2.0 (`stdio`). When configured in your editor, you get:

- **Inline Diagnostics**: Real-time syntax errors and validation failures right in your editor gutter.
- **Hover Documentation**: Hover over any exercise to read its objectives and hint cards rendered in rich markdown.
- **Code Actions**: Quick fixes to reveal progressive hints.

```bash
terralings lsp
```

### Visual Studio Code

Configure Terralings using any generic LSP client extension (such as [Generic LSP Client](https://marketplace.visualstudio.com/items?itemName=eyhn.vscode-generic-lsp)):

Add to `.vscode/settings.json`:

```json
{
  "generic-lsp.serverConfigurations": [
    {
      "name": "terralings",
      "command": "terralings",
      "args": ["lsp"],
      "languages": ["terraform", "hcl", "terraform-vars"],
      "rootMarkers": [".terralings", "exercises", ".git"]
    }
  ]
}
```

---

### Neovim (`nvim-lspconfig`)

Add the following to your Neovim configuration (e.g., `~/.config/nvim/lua/plugins/lsp.lua` or `init.lua`):

```lua
local lspconfig = require("lspconfig")
local configs = require("lspconfig.configs")

if not configs.terralings then
  configs.terralings = {
    default_config = {
      cmd = { "terralings", "lsp" },
      filetypes = { "terraform", "hcl" },
      root_dir = lspconfig.util.root_pattern(".terralings", "exercises", ".git"),
      settings = {},
    },
  }
end

lspconfig.terralings.setup({
  on_attach = function(client, bufnr)
    local opts = { buffer = bufnr, silent = true }
    vim.keymap.set("n", "K", vim.lsp.buf.hover, opts)
    vim.keymap.set("n", "<leader>ca", vim.lsp.buf.code_action, opts)
  end,
})
```

---

### Helix (`languages.toml`)

Add Terralings as a language server in `~/.config/helix/languages.toml`:

```toml
[language-server.terralings]
command = "terralings"
args = ["lsp"]

[[language]]
name = "hcl"
language-servers = [ "terralings" ]

[[language]]
name = "terraform"
language-servers = [ "terralings" ]
```

---

## 10. Curriculum Roadmap

The complete Terralings curriculum is structured into 13 chapters containing 56 progressive exercises:

| Chapter | Title | Focus Areas & Concepts | Exercises |
|---|---|---|:---:|
| **01** | **HCL Foundations & Primitives** | `terraform` block, providers, resources, dependencies, string interpolation, heredoc syntax, lifecycles | 6 |
| **02** | **Input Variables & Validations** | Primitive types, collection types (`list`, `map`, `set`), structural types (`object`, `tuple`), custom validations | 5 |
| **03** | **Outputs, Locals & Expressions** | Output values, `sensitive` redaction, local values, ternary conditionals, splat expressions (`[*]`) | 4 |
| **04** | **Built-in Functions & Collections** | String formatting, collection manipulation (`merge`, `slice`), JSON/YAML encoding, filesystem functions, `can()`/`try()` | 5 |
| **05** | **Meta-Arguments & Scaling** | Scaling with `count`, idempotent mapping with `for_each`, explicit `depends_on`, lifecycle management | 5 |
| **06** | **Dynamic Blocks & Nested HCL** | Repeating nested blocks with `dynamic`, custom iterators, nested dynamic blocks, conditional generation | 4 |
| **07** | **Data Sources & State Querying** | `data` sources (`local_file`, `archive_file`, `external`), lifecycle `precondition` & `postcondition` | 4 |
| **08** | **Modular Architecture** | Child module instantiation, input passing, bubbling outputs, composition, multi-instance `for_each`, provider aliases | 5 |
| **09** | **State Refactoring & Surgery** | Zero-downtime resource renaming via `moved` blocks, module extraction refactoring, `import` blocks, `replace_triggered_by` | 4 |
| **10** | **Native Unit & Integration Tests** | `.tftest.hcl` test suites, `command = plan` / `apply`, `mock_provider` blocks, negative testing with `expect_failures` | 4 |
| **11** | **Production Infrastructure Patterns** | Multi-tier composition, zero-downtime blue/green switching, tag inheritance, canary rollout patterns | 4 |
| **12** | **OpenTofu Advanced Features** | OpenTofu state file encryption (`key_provider`), early variable evaluation, custom registry mirrors | 3 |
| **13** | **Governance & Policy Guardrails** | Strict naming convention enforcement, security compliance rules (TLS, encryption at rest), resource quota guardrails | 3 |
| **Total** | | **Complete IaC Mastery** | **56** |

---

## 11. Troubleshooting & FAQ

### Q1: `terralings doctor` reports missing IaC engine binary. How do I install OpenTofu or Terraform?

**OpenTofu (Recommended):**
- **macOS (Homebrew):** `brew install opentofu`
- **Linux:** `curl -fsSL https://get.opentofu.org/install-opentofu.sh | sh`
- **Windows (Chocolatey / Scoop):** `choco install opentofu` or `scoop install opentofu`

**Terraform:**
- **macOS (Homebrew):** `brew tap hashicorp/tap && brew install hashicorp/tap/terraform`
- **Linux:** Install via HashiCorp official package repository.

If the binary is in a non-standard path, specify `--bin /path/to/binary` or set `export TERRALINGS_BIN=/path/to/binary`.

---

### Q2: Do I need AWS, GCP, or Azure credentials to complete the exercises?

**No.** All exercises execute completely locally and offline using standard local providers (`local`, `random`, `archive`, `terraform_data`, `test_mock`) or fast in-memory execution. You will never be billed for cloud infrastructure while learning with Terralings.

---

### Q3: How do I verify all exercises at once?

Run `terralings verify`:

```bash
terralings verify
```

This sequentially evaluates every exercise across all 13 chapters and displays a complete curriculum completion report.

---

### Q4: My editor is not showing error diagnostics. How do I troubleshoot `terralings lsp`?

1. Ensure `terralings` is accessible in your system `$PATH`.
2. Verify that your editor root is set to the directory containing the `exercises/` folder and `.terralings/` directory.
3. Start `terralings lsp` manually in a terminal to confirm it initializes without crashing.

---

### Q5: How do I generate shell completions for tab navigation?

Terralings supports automatic shell autocompletion for Bash, Zsh, Fish, and PowerShell:

```bash
# Zsh (macOS / Linux)
terralings completion zsh > "${fpath[1]}/_terralings"

# Bash (Linux)
terralings completion bash > /etc/bash_completion.d/terralings

# Fish
terralings completion fish > ~/.config/fish/completions/terralings.fish

# PowerShell
terralings completion powershell | Out-String | Invoke-Expression
```

---

### Q6: Where is my progress saved? Can I backup or transfer it?

Your progress and metrics are stored in `.terralings/state.json`. You can back up this file or use `--state <path>` to specify a custom storage location. If `.terralings/state.json` is ever corrupted, Terralings automatically restores from `.terralings/state.json.bak`.

---

*Happy Learning! For bug reports, suggestions, or contributions, visit [github.com/dnf0/terralings](https://github.com/dnf0/terralings).*
