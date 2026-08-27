# CLI Reference

Terralings provides a versatile command-line interface featuring 14 subcommands, interactive terminal user interfaces, JSON streaming modes, and Language Server Protocol daemons.

---

## Global Options & Environment Variables

These flags and environment variables apply across all Terralings commands:

### Global Flags

| Flag | Shorthand | Type | Default | Description |
|---|:---:|:---:|---|---|
| `--bin` | | `string` | *(auto-detected)* | Explicit path to the `tofu` or `terraform` executable |
| `--state` | | `string` | `.terralings/state.json` | Path to the progress and analytics state persistence file |
| `--help` | `-h` | `bool` | `false` | Display help and usage information for any command |

### Environment Variables

| Variable | Description |
|---|---|
| `TERRALINGS_BIN` or `TERRALINGS_BINARY` | Path override for the IaC engine binary (same as `--bin`). |
| `TERRALINGS_STATE_PATH` | Path to state persistence file (same as `--state`). |
| `TERRALINGS_PLUGIN_CACHE_DIR` | Shared provider plugin cache directory (default: `~/.terralings/plugin-cache`). |
| `NO_COLOR` | When set (any non-empty value), disables ANSI color escape sequences in terminal output. |

---

## Commands Summary

| Command | Synopsis | Description |
|---|---|---|
| [`terralings watch`](#terralings-watch) | `terralings watch [-i] [--json]` | Start continuous file watcher and live re-evaluation loop |
| [`terralings tui`](#terralings-tui) | `terralings tui` | Open full-screen interactive Bubble Tea terminal dashboard |
| [`terralings tour`](#terralings-tour) | `terralings tour [--step <n>] [--non-interactive] [--json]` | Launch the interactive 5-step guided onboarding walkthrough |
| [`terralings doctor`](#terralings-doctor) | `terralings doctor [--json]` | Execute pre-flight diagnostic health checks |
| [`terralings init`](#terralings-init) | `terralings init [target_dir] [-f]` | Extract the complete embedded exercise curriculum |
| [`terralings run`](#terralings-run) | `terralings run <exercise_name>` | Run standalone evaluation against a specific exercise |
| [`terralings hint`](#terralings-hint) | `terralings hint <exercise_name> [-i <index>]` | Display progressive hints for an exercise |
| [`terralings reset`](#terralings-reset) | `terralings reset <exercise_name> [-d <dir>]` | Reset an exercise file back to its initial template |
| [`terralings stats`](#terralings-stats) | `terralings stats` | View completion statistics, attempt metrics, and time invested |
| [`terralings search`](#terralings-search) | `terralings search <query>` | Full-text search across chapters, topics, hints, and exercises |
| [`terralings list`](#terralings-list) | `terralings list` | List all curriculum chapters, exercises, and status |
| [`terralings verify`](#terralings-verify) | `terralings verify` | Run full curriculum verification pass and show progress bar |
| [`terralings lsp`](#terralings-lsp) | `terralings lsp` | Start JSON-RPC 2.0 Language Server Protocol daemon on stdio |
| [`terralings completion`](#terralings-completion) | `terralings completion <shell>` | Generate shell autocompletion script (`bash`, `zsh`, `fish`, `powershell`) |
| [`terralings version`](#terralings-version) | `terralings version` | Print Terralings version and detected IaC engine binary |

---

## Command Details

### `terralings watch`

Start the continuous file watcher. Terralings watches `exercises/` using `fsnotify` and re-evaluates the active exercise automatically whenever changes are saved to disk.

```bash
terralings watch [flags]
```

#### Flags

- `-i, --interactive`: Launch the full-screen interactive Bubble Tea terminal dashboard (equivalent to `terralings tui`).
- `--json`: Emit structured Newline-Delimited JSON (NDJSON) event streams instead of formatted terminal text. Ideal for IDE plugins, CI pipelines, and external tooling.

#### Interactive Hotkeys (Standard Terminal Mode)

When running in standard watch mode, the following interactive single-key hotkeys are available:

| Key | Action |
|:---:|---|
| `h` | Display the next progressive hint for the active exercise |
| `r` | Manually re-run verification on the current exercise |
| `n` or `Enter` | Skip to the next exercise |
| `p` | Go back to the previous exercise |
| `q` or `Ctrl+C` | Quit watch mode and return to shell |

---

### `terralings tui`

Launch the full-screen interactive terminal dashboard powered by Bubble Tea and Lip Gloss.

```bash
terralings tui
```

#### Dashboard Keybindings

| Key | Action |
|:---:|---|
| `↑` / `k` | Navigate to previous exercise in tree |
| `↓` / `j` | Navigate to next exercise in tree |
| `Enter` | Select and evaluate highlighted exercise |
| `Tab` | Switch focus between curriculum tree and diagnostic output pane |
| `/` | Open live curriculum search filter modal |
| `h` | Toggle hints drawer for current exercise |
| `r` | Force re-evaluation of current exercise |
| `q` / `Esc` | Exit TUI dashboard |

---

### `terralings tour`

Launch the 5-step guided onboarding walkthrough. Introduces core IaC concepts, the Terralings feedback loop, watch and TUI modes, progressive hinting, and editor LSP integration.

```bash
terralings tour [flags]
```

#### Flags

- `--step <int>`: Render a specific tour step directly (1 to 5).
- `--non-interactive`: Render steps sequentially without waiting for interactive keypresses.
- `--json`: Emit tour content and step metadata as structured JSON.

---

### `terralings doctor`

Execute pre-flight diagnostic health checks to verify your system setup before starting exercises.

```bash
terralings doctor [flags]
```

#### Flags

- `--json`: Output diagnostic checks and results as machine-readable JSON.

#### Validated Checks

1. **IaC Engine Binary**: Checks for `tofu` or `terraform` in `$PATH` or custom `--bin` path.
2. **Curriculum Scaffold**: Verifies existence and count of files in `exercises/`.
3. **Provider Plugin Cache**: Verifies read/write access to `~/.terralings/plugin-cache`.
4. **Progress Persistence Store**: Tests read/write access to `.terralings/state.json`.
5. **Git Ignore Integration**: Confirms `.terralings/` is ignored in Git repositories.

---

### `terralings init`

Extract the entire embedded curriculum into your local workspace.

```bash
terralings init [target_dir] [flags]
```

#### Arguments

- `target_dir` *(optional, default: `exercises`)*: Directory where exercise chapters will be extracted.

#### Flags

- `-f, --force`: Overwrite existing files if the target directory already contains files.

---

### `terralings run`

Run verification against a single specified exercise without starting watch mode. Returns exit code `0` on success and `1` on failure.

```bash
terralings run <exercise_name>
```

#### Example

```bash
terralings run primitives01
# or with file path
terralings run exercises/01_primitives/primitives01.tf
```

---

### `terralings hint`

Display progressive hints for an exercise. Progressive hints are tiered: earlier hints provide conceptual guidance, while later hints provide concrete syntax examples.

```bash
terralings hint <exercise_name> [flags]
```

#### Flags

- `-i, --index <int>`: Zero-based index of the specific hint to display (default: `0`).

---

### `terralings reset`

Restore an exercise file back to its initial starter template from embedded assets. Useful if you want to restart an exercise from scratch.

```bash
terralings reset <exercise_name> [flags]
```

#### Flags

- `-d, --dir <string>`: Base exercises directory (default: `exercises`).

---

### `terralings stats`

Display detailed learning analytics and curriculum progress stored in `.terralings/state.json`.

```bash
terralings stats
```

#### Output Metrics

- Completed exercises count and overall completion percentage.
- Total verification attempts (passed vs. failed).
- Total hints requested across chapters.
- Estimated total time invested.
- Chapter-by-chapter progress breakdown.

---

### `terralings search`

Perform fast full-text search across all chapter titles, exercise names, descriptions, and progressive hints.

```bash
terralings search <query>
```

#### Example

```bash
terralings search "dynamic blocks"
terralings search "encryption"
terralings search "moved"
```

---

### `terralings list`

List all 13 chapters and 56 exercises with their current completion status indicators (`[✓]` completed, `[ ]` pending).

```bash
terralings list
```

---

### `terralings verify`

Run a full evaluation pass across the entire curriculum and display an aggregated progress bar and chapter scorecard.

```bash
terralings verify
```

---

### `terralings lsp`

Start the Language Server Protocol daemon over standard input/output (`stdio`). Implements JSON-RPC 2.0 LSP specifications for code diagnostics, hover tooltips, and code action commands.

```bash
terralings lsp
```

---

### `terralings completion`

Generate shell autocompletion scripts.

```bash
terralings completion <bash|zsh|fish|powershell>
```

*Alias*: `terralings completions`

---

### `terralings version`

Print the Terralings CLI version and details of the detected IaC engine binary.

```bash
terralings version
```

#### Output Example

```text
terralings v0.2.0
Detected binary: /usr/local/bin/tofu (OpenTofu v1.8.0)
```
