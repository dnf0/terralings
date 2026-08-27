# Terralings VS Code Extension & Shell Autocompletions Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a native VS Code companion extension (`extensions/vscode`) with LSP client, exercises explorer tree view, interactive walkthroughs, and status bar, plus CLI shell autocompletion for bash, zsh, fish, and powershell in `terralings`.

**Architecture:**
- Go CLI: Native Cobra shell completion generation (`terralings completion`) and dynamic argument completion functions for all 56 exercises and 13 chapters across commands.
- VS Code Extension: TypeScript extension using `@vscode/test-electron`, `vscode-languageclient`, and `esbuild`. Implements a `TreeDataProvider` reading static exercise metadata and `.terralings/state.json`, an embedded LSP client bridging to `terralings lsp`, and declarative `contributes.walkthroughs`.

**Tech Stack:** Go 1.22+, Cobra, TypeScript, VS Code Extension API, esbuild, Mocha.

---

### Task 1: CLI Shell Autocompletions & Dynamic Argument Completion

**Files:**
- Modify: `cmd/terralings/main.go`
- Create: `test/cli_completion_test.go`

- [ ] **Step 1: Write the failing tests for shell completions**

Create `test/cli_completion_test.go`:
```go
package test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestCLI_Completion_Bash(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/terralings", "completion", "bash")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("completion bash failed: %v, stderr: %s", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "bash completion") && !strings.Contains(out, "__terralings_") && !strings.Contains(out, "complete -o default -F") {
		t.Errorf("expected bash completion script, got:\n%s", out)
	}
}

func TestCLI_Completion_Zsh(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/terralings", "completion", "zsh")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("completion zsh failed: %v, stderr: %s", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "compdef") && !strings.Contains(out, "zsh completion") {
		t.Errorf("expected zsh completion script, got:\n%s", out)
	}
}

func TestCLI_Completion_Fish(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/terralings", "completion", "fish")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("completion fish failed: %v, stderr: %s", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "fish completion") && !strings.Contains(out, "complete -c terralings") {
		t.Errorf("expected fish completion script, got:\n%s", out)
	}
}

func TestCLI_Completion_PowerShell(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/terralings", "completion", "powershell")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("completion powershell failed: %v, stderr: %s", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "Register-ArgumentCompleter") && !strings.Contains(out, "powershell completion") {
		t.Errorf("expected powershell completion script, got:\n%s", out)
	}
}

func TestCLI_DynamicExerciseCompletion(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/terralings", "__complete", "run", "prim")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("__complete run prim failed: %v, stderr: %s", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "primitives01") {
		t.Errorf("expected dynamic completion to return 'primitives01', got:\n%s", out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/cli_completion_test.go`  
Expected: FAIL (`completion` command not recognized)

- [ ] **Step 3: Implement `completion` subcommand and dynamic `ValidArgsFunction` in `cmd/terralings/main.go`**

Add completion command and hook `ValidArgsFunction` to `runCmd`, `hintCmd`, `resetCmd`, `searchCmd`:
```go
var completionNoDesc bool

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion script for Terralings.
To load completions:

Bash:
  $ source <(terralings completion bash)

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ terralings completion zsh > "${fpath[1]}/_terralings"

Fish:
  $ terralings completion fish | source

PowerShell:
  PS> terralings completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletionV2(os.Stdout, !completionNoDesc)
		case "zsh":
			if completionNoDesc {
				return cmd.Root().GenZshCompletionNoDesc(os.Stdout)
			}
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, !completionNoDesc)
		case "powershell":
			if completionNoDesc {
				return cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return fmt.Errorf("unsupported shell: %s", args[0])
		}
	},
}

func exerciseArgsCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	m, err := manifest.LoadManifest()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var completions []string
	for _, ex := range m.Exercises {
		if strings.HasPrefix(ex.Name, toComplete) {
			completions = append(completions, fmt.Sprintf("%s\t%s (%s)", ex.Name, ex.Title, ex.Chapter))
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
```

In `init()` / command setup:
```go
completionCmd.Flags().BoolVar(&completionNoDesc, "no-descriptions", false, "disable completion descriptions")
rootCmd.AddCommand(completionCmd)

runCmd.ValidArgsFunction = exerciseArgsCompletion
hintCmd.ValidArgsFunction = exerciseArgsCompletion
resetCmd.ValidArgsFunction = exerciseArgsCompletion
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -v ./test/cli_completion_test.go`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/terralings/main.go test/cli_completion_test.go
git commit -m "feat(cli): add shell completion command and dynamic argument completions"
```

---

### Task 2: VS Code Extension Manifest, Assets & Interactive Walkthrough

**Files:**
- Create: `extensions/vscode/package.json`
- Create: `extensions/vscode/tsconfig.json`
- Create: `extensions/vscode/esbuild.js`
- Create: `extensions/vscode/.vscodeignore`
- Create: `extensions/vscode/icons/terralings.svg`
- Create: `extensions/vscode/walkthrough/01_welcome.md`
- Create: `extensions/vscode/walkthrough/02_anatomy.md`
- Create: `extensions/vscode/walkthrough/03_watch.md`
- Create: `extensions/vscode/walkthrough/04_tui_hints.md`
- Create: `extensions/vscode/walkthrough/05_editor_lsp.md`

- [ ] **Step 1: Create `extensions/vscode/package.json`**

Configure metadata, scripts (`build`, `watch`, `compile`, `test`), dependencies (`vscode-languageclient`), views container `terralings-exercises`, view `terralings-tree`, commands, and walkthrough cards.

- [ ] **Step 2: Create `tsconfig.json`, `esbuild.js`, and `.vscodeignore`**

Setup fast TypeScript compiler and bundle configuration.

- [ ] **Step 3: Create Walkthrough Markdown Files & Icons**

Populate all 5 walkthrough steps in `extensions/vscode/walkthrough/` and SVG brand icon in `extensions/vscode/icons/terralings.svg`.

- [ ] **Step 4: Install npm dependencies in `extensions/vscode`**

Run: `cd extensions/vscode && npm install`

- [ ] **Step 5: Commit**

```bash
git add extensions/vscode/
git commit -m "feat(vscode): scaffold extension manifest, walkthroughs, and build scripts"
```

---

### Task 3: VS Code Extension LSP Client, CLI Bridge & Status Bar

**Files:**
- Create: `extensions/vscode/src/lspClient.ts`
- Create: `extensions/vscode/src/cliRunner.ts`
- Create: `extensions/vscode/src/statusBar.ts`

- [ ] **Step 1: Implement `extensions/vscode/src/lspClient.ts`**

Implement `LanguageClient` wrapper that:
- Finds `terralings` executable via configuration `terralings.binaryPath`, workspace `./bin/terralings`, or system PATH.
- Spawns `terralings lsp` with stdio transport.
- Configures documentSelector for `terraform` and `hcl`.

- [ ] **Step 2: Implement `extensions/vscode/src/cliRunner.ts`**

Helper functions to:
- Open integrated terminal running `terralings watch` or `terralings tui`.
- Execute `terralings run <exercise>` or `terralings doctor` and output to dedicated `vscode.OutputChannel`.

- [ ] **Step 3: Implement `extensions/vscode/src/statusBar.ts`**

Status bar widget displaying progress count and percentage with click trigger opening the tree explorer.

- [ ] **Step 4: Compile TypeScript**

Run: `npm run compile` in `extensions/vscode`  
Expected: Clean compile with zero errors.

- [ ] **Step 5: Commit**

```bash
git add extensions/vscode/src/
git commit -m "feat(vscode): implement LSP client, CLI runner bridge, and status bar"
```

---

### Task 4: Curriculum Tree Data Provider & Real-Time State Sync

**Files:**
- Create: `extensions/vscode/src/treeProvider.ts`
- Create: `extensions/vscode/src/stateWatcher.ts`
- Create: `extensions/vscode/src/extension.ts`

- [ ] **Step 1: Implement `extensions/vscode/src/treeProvider.ts`**

Implement `vscode.TreeDataProvider<TerralingsTreeItem>`:
- Encapsulates 13 chapters and 56 exercises.
- Reads `.terralings/state.json` to assign status icons (passed, failed, in-progress, unattempted).
- Clicking an exercise executes `vscode.commands.executeCommand('vscode.open', fileUri)`.

- [ ] **Step 2: Implement `extensions/vscode/src/stateWatcher.ts`**

Listens to filesystem changes on `.terralings/state.json` and invokes `treeProvider.refresh()` and `statusBar.update()`.

- [ ] **Step 3: Implement `extensions/vscode/src/extension.ts`**

Activate/deactivate lifecycle:
- Registers all commands (`terralings.watch`, `terralings.tui`, `terralings.runCurrent`, `terralings.hint`, `terralings.reset`, `terralings.doctor`, `terralings.tour`).
- Starts LSP client if enabled.
- Registers TreeDataProvider and FileSystemWatcher.

- [ ] **Step 4: Build extension bundle with esbuild**

Run: `npm run build` in `extensions/vscode`  
Expected: Generates `dist/extension.js`.

- [ ] **Step 5: Commit**

```bash
git add extensions/vscode/src/ extensions/vscode/dist/
git commit -m "feat(vscode): implement curriculum tree provider and state sync"
```

---

### Task 5: VS Code Extension Tests, Documentation & Makefile Integration

**Files:**
- Create: `extensions/vscode/test/suite/extension.test.ts`
- Create: `extensions/vscode/test/suite/treeProvider.test.ts`
- Create: `extensions/vscode/README.md`
- Modify: `Makefile`
- Modify: `README.md`

- [ ] **Step 1: Write extension unit tests in `extensions/vscode/test/`**

Test tree provider hierarchy creation, mock state evaluation, binary path resolution, and walkthrough files existence.

- [ ] **Step 2: Run extension tests and Go test suite**

Run: `cd extensions/vscode && npm test` and `make all`  
Expected: All tests pass.

- [ ] **Step 3: Update `extensions/vscode/README.md`, top-level `README.md`, and `Makefile`**

Add `make extension-build` / `make extension-test` targets to `Makefile`. Document the VS Code extension installation and features in `README.md`.

- [ ] **Step 4: Commit**

```bash
git add extensions/vscode/ Makefile README.md
git commit -m "docs(vscode): add extension tests, documentation, and makefile targets"
```

---
