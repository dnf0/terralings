# Terralings

> An interactive, terminal-driven learning environment for mastering Terraform and OpenTofu through hands-on exercises.

`terralings` guides you through fixing broken configurations, writing declarative Infrastructure as Code (IaC), mastering HCL expressions and built-in functions, refactoring state with `moved` blocks, authoring `.tftest.hcl` unit and integration tests, and configuring OpenTofu advanced features such as state encryption and early variable evaluation.

Inspired by [`rustlings`](https://github.com/rust-lang/rustlings), [`ziglings`](https://github.com/ziglings/exercises), and [`raylings`](https://github.com/ray-project/raylings).

---

## Features

- **Turnkey Embedded Initialization (`terralings init`)**: The complete 56-exercise curriculum is embedded directly into the binary—run `terralings init` anywhere to start practicing immediately without cloning git repos.
- **Interactive Guided Tour (`terralings tour`)**: 5-step interactive terminal walkthrough introducing IaC principles, exercise structure, watch/TUI modes, progressive hints, and LSP editor integration.
- **Pre-Flight Diagnostics (`terralings doctor`)**: Comprehensive environment health check verifying OpenTofu/Terraform binary availability, exercise directory scaffold integrity, plugin cache permissions, and state health.
- **Interactive Full-Screen TUI Dashboard (`terralings tui` / `terralings watch -i`)**: Two-pane terminal dashboard featuring real-time compilation viewport, grouped chapter navigation, live exercise search modal (`/`), and expandable hints drawer (`h`).
- **Learning Analytics & Progress Persistence (`terralings stats`)**: Automatically tracks completion state, pass/fail attempt metrics, time invested, and hint usage stored in `.terralings/state.json`.
- **Language Server Protocol Daemon (`terralings lsp`)**: Built-in JSON-RPC 2.0 LSP server providing instant in-editor validation diagnostics, hover exercise descriptions with hints, and quick-fix actions for VS Code, Neovim, and Helix.
- **NDJSON Event Streaming (`terralings watch --json`)**: Emits structured newline-delimited JSON events on file changes and evaluations for headless integrations, CI test runners, and custom editor plugins.
- **Interactive Watch Mode (`terralings watch`)**: Automatically monitors your exercise files via `fsnotify` and re-evaluates and validates in real time on every file save.
- **Exercise Reset (`terralings reset <name>`)**: Instantly restore any exercise back to its clean starter template if you want to redo it or fix a mistake.
- **Curriculum Search (`terralings search <term>`)**: Fast full-text search across all chapters, topics, hints, and exercises with relevance scoring.
- **Shell Autocompletions (`terralings completion`)**: Rich tab completion for Bash, Zsh, Fish, and PowerShell, including interactive exercise name completion.
- **Dual Engine Support**: Seamlessly detects and runs against either **OpenTofu** (`tofu`) or **Terraform** (`terraform`) (version >= 1.6.0).
- **Sub-100ms Evaluation**: Shared provider plugin caching eliminates redundant network downloads during provider initialization.
- **Comprehensive 13-Chapter Curriculum**: 56 progressive exercises covering primitives, variables, collections, functions, meta-arguments, dynamic blocks, data sources, modules, state refactoring, native testing, production patterns, OpenTofu extensions, and policy governance.
- **Progressive Hints (`terralings hint`)**: Multi-level contextual guidance to nudge you forward when stuck without spoiling the answer.
- **Built-in Verification (`terralings verify`)**: Complete curriculum progress dashboard tracking your status across all 13 chapters.
- **Zero Heavy Cloud Dependencies**: All exercises execute purely locally using standard built-in providers (`local`, `random`, `archive`, `terraform_data`, `test_mock`) or fast in-memory execution—no AWS/GCP/Azure credentials required.

---

## Prerequisites

Before using `terralings`, ensure you have installed:

1. **OpenTofu** (>= 1.6.0) or **Terraform** (>= 1.6.0):
   ```bash
   tofu version
   # or
   terraform version
   ```
2. *(Optional)* **Go** (>= 1.22) if building from source.

---

## Installation & Getting Started

### Option 1: Quick Install (Linux & macOS)

Install the latest pre-built binary directly to `~/.local/bin` or `/usr/local/bin`:

```bash
curl -fsSL https://raw.githubusercontent.com/dnf0/terralings/main/install.sh | bash
```

### Option 2: Go Toolchain

If you have Go 1.22+ installed:

```bash
go install github.com/dnf0/terralings/cmd/terralings@latest
```

### Option 3: GitHub Releases Pre-built Binaries

Download pre-compiled archives for Linux (`amd64`, `arm64`), macOS (Apple Silicon `arm64`, Intel `amd64`), or Windows from [GitHub Releases](https://github.com/dnf0/terralings/releases).

### Option 4: Build from Source

```bash
git clone https://github.com/dnf0/terralings.git
cd terralings
make build
# Binary is located at ./bin/terralings
```

---

## Quickstart

Once installed, verify your environment, take the guided tour, and scaffold the exercises:

```bash
# 1. Run pre-flight health checks to verify your OpenTofu/Terraform setup
terralings doctor

# 2. Take the 2-minute interactive guided walkthrough
terralings tour

# 3. Initialize exercises in current directory (creates exercises/ folder)
terralings init

# 4. Start interactive learning loop (standard terminal stream)
terralings watch

# 5. Or launch the full-screen interactive TUI dashboard
terralings tui
# or
terralings watch -i
```

> 📖 **New to Terralings?** Check out the comprehensive [Onboarding Guide](docs/onboarding-guide.md) for step-by-step instructions, editor configurations, and learning tips.

---

## CLI Commands Reference

| Command | Description |
|---|---|
| `terralings doctor [--json]` | Run pre-flight diagnostics to verify environment and workspace readiness |
| `terralings tour [--step <n>] [--non-interactive] [--json]` | Start the interactive guided onboarding tour |
| `terralings init [dir] [-f]` | Extract and initialize embedded curriculum exercises into a directory |
| `terralings watch [-i] [--json]` | Start continuous watch mode (`-i` for TUI, `--json` for NDJSON stream) |
| `terralings tui` | Launch the interactive full-screen terminal UI dashboard |
| `terralings stats` | Display learning analytics, attempt counts, time invested, and progress |
| `terralings lsp` | Start Language Server Protocol (LSP) daemon over stdio |
| `terralings run <exercise>` | Run verification against a single exercise (e.g. `terralings run primitives01`) |
| `terralings hint <exercise> [-i <idx>]` | Display progressive hint(s) for the specified exercise |
| `terralings reset <exercise> [-d <dir>]` | Reset an exercise back to its initial starting template |
| `terralings search <term>` | Search exercises by keyword, concept, or chapter |
| `terralings list` | List all chapters and exercises with status indicators |
| `terralings verify` | Run sequential evaluation across the entire curriculum and display progress |
| `terralings completion <shell>` | Generate autocompletion scripts for `bash`, `zsh`, `fish`, or `powershell` (alias: `completions`) |
| `terralings version` | Print the Terralings CLI version and detected IaC binary |

### Global Flags

| Flag | Description |
|---|---|
| `--bin <path>` | Override binary auto-detection with an explicit path (or `export TERRALINGS_BIN=...`) |
| `--state <path>` | Custom path to state persistence file (default: `.terralings/state.json`) |

---

## Features & Usage Guide

### 1. Interactive Full-Screen TUI Dashboard

Launch a terminal dashboard with split-pane navigation, real-time evaluation feedback, and searchable curriculum overview:

```bash
terralings tui
# or
terralings watch -i
```

#### TUI Layout & Features

- **Sidebar (Left Pane)**: Visual tree of chapters and exercises with real-time status indicators (`✓` passed, `•` in progress, `·` not started).
- **Compiler Viewport (Right Pane)**: Formatted validation output, syntax and runtime diagnostics, error callouts, and animated spinner during execution.
- **Collapsible Hints Drawer**: View progressive hints on demand directly within the dashboard without leaving your terminal.
- **Interactive Search Modal**: Press `/` to trigger instant search across exercises and jump directly to any topic.

#### Keyboard Shortcuts

| Key | Action |
|---|---|
| `↑` / `k` | Move cursor up in exercise list / search results |
| `↓` / `j` | Move cursor down in exercise list / search results |
| `Tab` | Switch active focus between Sidebar and Viewport |
| `Enter` | Run/evaluate selected exercise (or select item in search) |
| `h` | Toggle hints drawer / cycle through hint levels |
| `r` | Reset current exercise back to original starter template |
| `/` | Open exercise search modal |
| `Esc` | Close search modal or cancel |
| `q` / `Ctrl+C` | Quit dashboard |

---

### 2. Learning Analytics & Progress Persistence

Terralings tracks your learning journey across all sessions and exercises. Progress is automatically persisted to `.terralings/state.json` (and automatically added to your `.gitignore`).

```bash
terralings stats
```

#### State Tracking Mechanics

- **Attempt Tracking**: Increments total and failed attempts per exercise each time validation runs.
- **Time Invested**: Tracks cumulative active time spent solving exercises.
- **Hint Consultation**: Records the depth of hints consulted.
- **Atomic & Resilient**: State updates are written atomically (`.tmp` + `sync` + `rename`) with automatic recovery and backup (`state.json.bak`) if corruption is detected.
- **Custom State Location**: Use `--state <path>` to store state in a custom directory.

#### Sample Output

```text
📊 TERRALINGS LEARNING ANALYTICS

Overall Progress: [████████████░░░░░░░░] 60% (34/56 completed)
Time Invested:    1h 45m
Total Attempts:   82 (avg 1.5 per exercise)
Hints Consulted:  12

Chapter Breakdown:
  • 01_primitives        100% (6/6) | Avg Attempts: 1.2 | Hints: 1
  • 02_variables         100% (5/5) | Avg Attempts: 1.4 | Hints: 2
  • 03_outputs_locals    100% (4/4) | Avg Attempts: 1.0 | Hints: 0
  • 04_functions         100% (5/5) | Avg Attempts: 1.8 | Hints: 3
  • 05_meta_arguments    100% (5/5) | Avg Attempts: 1.6 | Hints: 2
  • 06_dynamic_blocks     75% (3/4) | Avg Attempts: 2.0 | Hints: 2
  • 07_data_sources       50% (2/4) | Avg Attempts: 1.5 | Hints: 1
  • 08_modules             0% (0/5) | Avg Attempts: 0.0 | Hints: 0
```

---

### 3. Language Server Protocol (LSP) Daemon

Terralings includes a built-in Language Server Protocol (LSP) daemon communicating over JSON-RPC 2.0 (`stdio`). When configured in your editor, it provides:

- **Live Diagnostics (`publishDiagnostics`)**: Real-time syntax errors, semantic validation issues, and plan diagnostics.
- **Hover Documentation (`textDocument/hover`)**: Hover over any exercise file to read its objectives, chapter context, and multi-level hints rendered in rich markdown.
- **Code Actions (`textDocument/codeAction`)**: Quick-fix and source actions to reveal progressive hints.

```bash
terralings lsp
```

#### Editor Configuration

##### Neovim (`nvim-lspconfig`)

Add the following to your Neovim configuration (e.g. `~/.config/nvim/lua/plugins/lsp.lua` or `init.lua`):

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

##### Helix (`languages.toml`)

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

##### Visual Studio Code

You can use `terralings lsp` in VS Code with any generic LSP extension (such as [Generic LSP Client](https://marketplace.visualstudio.com/items?itemName=eyhn.vscode-generic-lsp) or custom extension settings in `.vscode/settings.json`):

```json
{
  "generic-lsp.serverConfigurations": [
    {
      "name": "terralings",
      "command": "terralings",
      "args": ["lsp"],
      "languages": ["terraform", "hcl", "terraform-vars"]
    }
  ]
}
```

---

### 4. NDJSON Event Streaming

For automated tooling, continuous integration environments, and custom IDE extensions, `terralings watch --json` streams newline-delimited JSON (NDJSON) events directly to `stdout`:

```bash
terralings watch --json
```

#### Event Types & Payload Schema

Each line emitted is an atomic JSON object:

1. **`exercise_start`**: Emitted when evaluation begins for an exercise.
2. **`exercise_result`**: Emitted when evaluation completes with diagnostics, pass/fail status, and raw CLI output.
3. **`completed`**: Emitted when all curriculum exercises pass.

```json
{
  "event": "exercise_result",
  "timestamp": "2026-08-27T12:00:00Z",
  "exercise": {
    "name": "primitives01",
    "title": "Terraform Configuration Block",
    "path": "exercises/01_primitives/primitives01.tf",
    "mode": "validate",
    "chapter_name": "01_primitives",
    "hints": [
      "Use required_version = \">= 1.6.0\"",
      "Configure required_providers block with local provider"
    ]
  },
  "passed": false,
  "diagnostics": [
    {
      "file": "exercises/01_primitives/primitives01.tf",
      "line": 4,
      "column": 1,
      "severity": "error",
      "summary": "Missing required argument"
    }
  ],
  "raw_output": "Error: Missing required argument...",
  "exit_code": 1,
  "current_index": 0,
  "total_count": 56
}
```

---

### 5. Command Examples & Standard Watch Mode

#### 1. Pre-Flight Diagnostics
```bash
terralings doctor
```
```text
🩺 Terralings Doctor Diagnostic Report
────────────────────────────────────────────────────────────

 ✓ IaC Engine Binary
   Found opentofu at /usr/local/bin/tofu (OpenTofu v1.8.0)

 ✓ Curriculum Scaffold
   Exercises directory present (56 configuration files found)

 ✓ Provider Plugin Cache
   Plugin cache directory ready at ~/.terralings/plugin-cache

 ✓ Git Ignore Integration
   .terralings directory is properly git-ignored.

 ✓ Progress Persistence Store
   State store healthy at .terralings/state.json (0 completed, 0 attempts)

────────────────────────────────────────────────────────────
 All diagnostics passed! Your environment is 100% ready for Terralings.
```

#### 2. Interactive Guided Tour
```bash
terralings tour
```
```text
 STEP 1 OF 5  Welcome & Core Philosophy
 Master Terraform & OpenTofu through interactive hands-on practice

   Terralings is designed to teach you Infrastructure-as-Code from first principles.
   All exercises run in isolated, sandboxed environments without requiring real cloud credentials.
   We follow the Ziglings / Rustlings v6 model: pure deterministic validation with zero magic comment friction.

   Example Command:
   terralings watch

   Key Takeaways:
   ✓ 100% local, safe evaluation with OpenTofu / Terraform.
   ✓ Real compiler errors & plan outputs guide your progress.

[Enter / n] Next | [p] Prev | [1-5] Jump | [q] Quit
>
```

#### 3. Interactive Watch Mode
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

#### 4. Progressive Hints
```bash
terralings hint primitives01
```
```text
╭───────────────────────────────────╮
│                                   │
│ 💡 Hint (1/2) for primitives01:   │
│ Use required_version = ">= 1.6.0" │
│                                   │
╰───────────────────────────────────╯
```

View the next hint level:
```bash
terralings hint primitives01 --index 1
```

#### 5. Curriculum Overview
```bash
terralings list
```
```text
⚡ TERRALINGS: Master Terraform & OpenTofu from Scratch ⚡

Chapter 01: HCL Foundations & Core Primitives - Blocks, attributes, provider requirements and first resources
  · primitives01     : Terraform Configuration Block (exercises/01_primitives/primitives01.tf)
  · primitives02     : First Resource Declaration (exercises/01_primitives/primitives02.tf)
  · primitives03     : Resource Dependencies (exercises/01_primitives/primitives03.tf)
  · primitives04     : String Interpolation & Heredoc (exercises/01_primitives/primitives04.tf)
  · primitives05     : Syntax & Formatting (exercises/01_primitives/primitives05.tf)
  · primitives06     : Lifecycle Mechanics (exercises/01_primitives/primitives06.tf)
...
```

#### 6. Curriculum Search
```bash
terralings search dynamic
```
```text
🔍 Search Results for 'dynamic' (4 matches):
  • dynamic01        Repeating Nested Blocks with dynamic
    exercises/06_dynamic_blocks/dynamic01.tf | matched in: name, chapter
  • dynamic02        Custom Iterator Variable
    exercises/06_dynamic_blocks/dynamic02.tf | matched in: name, chapter
...
```

#### 7. Resetting an Exercise
```bash
terralings reset primitives01
```
```text
🔄 Reset exercise 'primitives01' (exercises/01_primitives/primitives01.tf) back to original template.
```

#### 8. Shell Autocompletions
Enable shell autocompletions for quick tab navigation and exercise autocompletion:

```bash
# Bash (Linux)
terralings completion bash > /etc/bash_completion.d/terralings

# Zsh (macOS / Linux)
terralings completion zsh > "${fpath[1]}/_terralings"

# Fish
terralings completion fish > ~/.config/fish/completions/terralings.fish

# PowerShell
terralings completion powershell | Out-String | Invoke-Expression
```

#### 9. Progress Verification
```bash
terralings verify
```
```text
Progress: [████████████████████████████████████████] 56/56 (100.0%)

🎉 Congratulations! You have completed all Terralings exercises! 🎉
```

### Environment Variables & Global Flags

- `--bin <path>`: Override the automatic binary detector with an explicit executable path (e.g. `terralings watch --bin /opt/homebrew/bin/tofu`).
- `TERRALINGS_BIN`: Set the binary path via environment variable (e.g. `export TERRALINGS_BIN=/usr/local/bin/terraform`).

---

## Curriculum Map

The curriculum spans **13 structured chapters** containing **56 exercises**:

### [Chapter 01: HCL Foundations & Core Primitives](exercises/01_primitives/)
- `primitives01`: `terraform` configuration block, `required_version`, and `required_providers`.
- `primitives02`: First resource declaration (`local_file`), block labels, and attributes.
- `primitives03`: Implicit and explicit resource dependencies.
- `primitives04`: String interpolation, heredoc syntax (`<<-EOT`), and escaping.
- `primitives05`: Canonical HCL syntax rules and formatting conventions.
- `primitives06`: Resource lifecycle mechanics and stateful recreation triggers.

### [Chapter 02: Input Variables, Types & Validations](exercises/02_variables/)
- `variables01`: Primitive variable types (`string`, `number`, `bool`).
- `variables02`: Collection types (`list(string)`, `map(string)`, `set(string)`).
- `variables03`: Structural types (`object({...})`, `tuple([...])`).
- `variables04`: Default values and nullable handling.
- `variables05`: Custom `validation` blocks with condition expressions and error messages.

### [Chapter 03: Outputs, Locals & Expressions](exercises/03_outputs_locals/)
- `outputs01`: Defining output values and `sensitive = true` redaction.
- `locals01`: Intermediate local calculations for DRY configurations.
- `expr01`: Ternary conditional expressions (`condition ? true_val : false_val`).
- `expr02`: Splat expressions (`[*]`) and list transformations.

### [Chapter 04: Built-in Functions & Collections](exercises/04_functions/)
- `func01`: String manipulation functions (`format`, `join`, `lower`, `trimspace`).
- `func02`: Collection functions (`merge`, `concat`, `distinct`, `slice`).
- `func03`: Data serialization functions (`jsonencode`, `yamlencode`, `jsondecode`).
- `func04`: Filesystem and hashing functions (`file`, `fileexists`, `sha256`, `abspath`).
- `func05`: Safe evaluation functions (`can()`, `try()`).

### [Chapter 05: Meta-Arguments & Resource Scaling](exercises/05_meta_arguments/)
- `meta01`: Scaling resources with `count` and `count.index`.
- `meta02`: Idempotent resource mapping with `for_each` and `each.key`/`each.value`.
- `meta03`: Explicit dependency ordering with `depends_on`.
- `meta04`: Lifecycle control (`create_before_destroy`, `prevent_destroy`, `ignore_changes`).
- `meta05`: Chained dependencies across maps and resources.

### [Chapter 06: Dynamic Blocks & Advanced HCL](exercises/06_dynamic_blocks/)
- `dynamic01`: Repeating nested blocks using `dynamic "block_name"`.
- `dynamic02`: Custom block iterators with `iterator`.
- `dynamic03`: Hierarchical and multi-level nested dynamic blocks.
- `dynamic04`: Conditional block generation using empty collections.

### [Chapter 07: Data Sources & State Querying](exercises/07_data_sources/)
- `data01`: Reading local filesystem artifacts with `data "local_file"`.
- `data02`: Bundling archive files with `data "archive_file"`.
- `data03`: Querying external processes and JSON structures with `data "external"`.
- `data04`: Adding custom `precondition` and `postcondition` lifecycle validations.

### [Chapter 08: Modular Infrastructure Architecture](exercises/08_modules/)
- `module01`: Child module instantiation, input passing, and directory layout.
- `module02`: Consuming and bubbling child module outputs to root caller.
- `module03`: Composing multiple interdependent modules.
- `module04`: Multi-instance module provisioning with `for_each`.
- `module05`: Passing explicit provider configurations via `providers = { ... }` aliases.

### [Chapter 09: State Refactoring & Surgery](exercises/09_state_refactoring/)
- `state01`: Renaming resources without destruction using declarative `moved` blocks.
- `state02`: Migrating monolithic resources into encapsulated child modules via `moved` blocks.
- `state03`: Declarative resource adoption using `import` blocks.
- `state04`: Managing forced recreation triggers using `replace_triggered_by`.

### [Chapter 10: Native Unit & Integration Testing](exercises/10_testing/)
- `test01`: Writing native `.tftest.hcl` test suites with `command = plan` run blocks.
- `test02`: Asserting against execution side-effects with `command = apply`.
- `test03`: Mocking third-party providers with `mock_provider` blocks.
- `test04`: Testing negative validation scenarios using `expect_failures`.

### [Chapter 11: Production Infrastructure Patterns](exercises/11_patterns/)
- `pattern01`: Multi-tier infrastructure composition (networking -> compute -> storage).
- `pattern02`: Zero-downtime blue/green deployment switching with `terraform_data`.
- `pattern03`: Centralized tag and metadata inheritance across resources.
- `pattern04`: Idempotent resource cleanup and canary rollout patterns.

### [Chapter 12: OpenTofu Advanced Features](exercises/12_opentofu/)
- `tofu01`: Configuring OpenTofu state file encryption with `encryption` and `key_provider` blocks.
- `tofu02`: Early variable and local evaluation in provider and backend configurations.
- `tofu03`: Configuring custom registry mirrors and decentralized OpenTofu registries.

### [Chapter 13: Governance & Security Policies](exercises/13_governance/)
- `gov01`: Enforcing strict resource naming conventions and regex constraints.
- `gov02`: Validating mandatory security compliance guardrails (TLS, encryption at rest).
- `gov03`: Cost and quota governance guardrails restricting over-provisioning.

---

## Development & Testing

For contributors looking to build, test, or extend Terralings:

```bash
# Run all format checks, linters, and tests
make all

# Run unit and curriculum tests with race detection
make test-race

# Run formatting checks across Go code and solutions
make check

# Build the terralings binary
make build
```

---

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, development workflow, and pull request guidelines.

---

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
