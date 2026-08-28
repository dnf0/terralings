# Intelligent Path Resolution, 1-Click Auto-Init UX & Core Engine Modernization Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement robust multi-depth path resolution and 1-click workspace auto-initialization in the VS Code extension, modernize core engine placeholder validation to check both comments and blanks (`___`, `/* ??? */`, `<!-- ANSWER -->`), add JSON support and state awareness to `terralings list`, and release `v0.4.0` with full CI/CD verification and marketplace publishing.

**Architecture:**
1. VS Code Extension: Add `pathUtils.ts` containing `getEffectiveWorkspaceRoot()` (preventing no-workspace `/` system crashes) and `resolveExercisePath()` (handling repo root, `exercises/` subfolder, track subfolders, and parent directory traversal). Enhance `extension.ts` and `cliRunner.ts` with auto-init prompts and `openNextExercise`.
2. Core Engine: Modernize `internal/runner/runner.go` `CheckMarker` to detect both legacy comment markers and unfilled inline blanks (`___`, `/* ??? */`, etc.).
3. CLI & State: Add `--json` flag to `terralings list` and render actual progress from `.terralings/state.json`.
4. Version & Parity: Bump to `0.4.0` in `internal/version/version.go`, `package.json`, and `CHANGELOG.md`.

**Tech Stack:** Go 1.23, TypeScript 5.5, VS Code Extension API, Mocha, Cobra CLI, Lipgloss, MkDocs.

---

### Task 1: Core Engine Placeholder & Incomplete Comment Check

**Files:**
- Modify: `internal/runner/runner.go`
- Modify: `test/runner_test.go`

- [ ] **Step 1: Write failing tests for unfilled blank detection and comment variations**

In `test/runner_test.go`, add test cases for:
- `___` inline blanks
- `/* ??? */` placeholder
- `<!-- ANSWER -->` placeholder
- `<!-- I AM NOT DONE -->`
- `// I AM NOT DONE`
- `# I AM NOT DONE`
- Clean completed HCL without any markers

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v -run TestCheckMarker ./test`
Expected: FAIL on blank placeholders.

- [ ] **Step 3: Implement multi-pattern marker and blank detection in `internal/runner/runner.go`**

Update `notDoneCommentRegex` and `unfilledBlankRegex`:
```go
var notDoneCommentRegex = regexp.MustCompile(`(?i)(<!--\s*i\s+am\s+not\s+done\s*-->|//\s*i\s+am\s+not\s+done|#\s*i\s+am\s+not\s+done|i\s+am\s+not\s+done)`)
var unfilledBlankRegex = regexp.MustCompile(`(___|\/\*\s*\?\?\?\s*\*\/|<!--\s*ANSWER\s*-->)`)

func CheckMarker(path string) bool {
    // Stat and read file/dir; return true if notDoneCommentRegex.Match(data) || unfilledBlankRegex.Match(data)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -v -run TestCheckMarker ./test`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/runner/runner.go test/runner_test.go
git commit --no-gpg-sign -m "feat(runner): detect inline blanks and incomplete comment markers"
```

---

### Task 2: CLI `terralings list` Improvements (`--json` & State Awareness)

**Files:**
- Modify: `cmd/terralings/main.go`
- Modify: `internal/ui/ui.go`
- Test: `test/cli_test.go` or `test/manifest_test.go`

- [ ] **Step 1: Write tests for `terralings list --json` and state-aware listing**

Add tests asserting `terralings list --json` outputs valid JSON containing all 56 exercises with `name`, `title`, `chapter`, `path`, `mode`, and `status`.

- [ ] **Step 2: Implement `--json` flag and state loading in `cmd/terralings/main.go`**

Add `--json` flag to `listCmd`:
```go
var listJSON bool
listCmd := &cobra.Command{
    Use: "list",
    Short: "List all curriculum chapters and exercises",
    RunE: func(cmd *cobra.Command, args []string) error {
        store, _ := state.NewStore(stateOverride)
        m := manifest.GetManifest()
        var statuses map[string]models.ExerciseStatus
        if store != nil {
            statuses = store.GetExerciseStatuses(m)
        }
        if listJSON {
            // output JSON array of exercises with status
        }
        // render styled table
    },
}
listCmd.Flags().BoolVar(&listJSON, "json", false, "Output exercise listing as JSON")
```

- [ ] **Step 3: Verify CLI list commands**

Run: `go run ./cmd/terralings list --json`
Expected: Valid JSON array with 56 exercises.

- [ ] **Step 4: Commit**

```bash
git add cmd/terralings/main.go internal/ui/ui.go test/
git commit --no-gpg-sign -m "feat(cli): add --json output and state awareness to terralings list"
```

---

### Task 3: Intelligent Path Resolution (`pathUtils.ts`)

**Files:**
- Create: `extensions/vscode/src/pathUtils.ts`
- Test: `extensions/vscode/test/suite/pathUtils.test.ts`

- [ ] **Step 1: Write comprehensive unit tests for `pathUtils.ts`**

Test:
- `getEffectiveWorkspaceRoot()` returns workspace folder when open.
- `getEffectiveWorkspaceRoot()` falls back safely to homedir or candidate repo directories, never `/`.
- `resolveExercisePath()` handles absolute paths, direct workspace root resolution, `exercises/` prefix stripping, and parent directory traversal up to 6 levels.

- [ ] **Step 2: Implement `extensions/vscode/src/pathUtils.ts`**

Implement:
```typescript
export function getEffectiveWorkspaceRoot(repoName: string = 'terralings'): string;
export function resolveExercisePath(exercisePath: string, workspaceRoot?: string): string;
```

- [ ] **Step 3: Run Mocha tests to verify**

Run: `cd extensions/vscode && npm test`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add extensions/vscode/src/pathUtils.ts extensions/vscode/test/suite/pathUtils.test.ts
git commit --no-gpg-sign -m "feat(vscode): implement intelligent multi-depth path resolution"
```

---

### Task 4: 1-Click Auto-Init UX, `openNextExercise`, and Sidebar Integration

**Files:**
- Modify: `extensions/vscode/src/treeProvider.ts`
- Modify: `extensions/vscode/src/extension.ts`
- Modify: `extensions/vscode/src/cliRunner.ts`
- Modify: `extensions/vscode/src/lspClient.ts`
- Modify: `extensions/vscode/package.json`
- Modify: `extensions/vscode/test/suite/extension.test.ts`
- Modify: `extensions/vscode/test/suite/treeProvider.test.ts`

- [ ] **Step 1: Update `treeProvider.ts` to use `resolveExercisePath` & robust state resolution**

Use `resolveExercisePath` in `resolveExerciseUri` and update `refresh()` to look for `.terralings/state.json` using `getEffectiveWorkspaceRoot()` and ancestor directory traversal.

- [ ] **Step 2: Implement 1-click workspace auto-initialization and `openNextExercise` in `extension.ts` and `cliRunner.ts`**

In `terralings.openExercise`:
- Check `fs.existsSync(resolvedPath)`.
- If missing, prompt user: `"Terralings exercises were not found on disk. Would you like to initialize them in \"${targetDir}\"?"` with buttons `['Initialize Exercises', 'Cancel']`.
- If confirmed, execute `terralings init` via `initExercises(targetDir)`, refresh tree/status bar, and open the document.
- Register commands `terralings.initExercises` and `terralings.openNextExercise`.

- [ ] **Step 3: Update `extensions/vscode/package.json` command declarations**

Add `terralings.initExercises` and `terralings.openNextExercise` to `contributes.commands`.

- [ ] **Step 4: Build and test extension**

Run: `cd extensions/vscode && npm run check-types && npm run build && npm test`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add extensions/vscode/
git commit --no-gpg-sign -m "feat(vscode): add 1-click auto-initialization UX and openNextExercise command"
```

---

### Task 5: Version Bump to v0.4.0 & Comprehensive Verification

**Files:**
- Modify: `internal/version/version.go`
- Modify: `extensions/vscode/package.json`
- Modify: `extensions/vscode/package-lock.json`
- Modify: `CHANGELOG.md`
- Modify: `docs/vscode-extension.md`
- Modify: `docs/cli-reference.md`
- Modify: `test/version_test.go`

- [ ] **Step 1: Bump version to `0.4.0` across all files**

Update:
- `internal/version/version.go`: `const Number = "0.4.0"`
- `extensions/vscode/package.json`: `"version": "0.4.0"`
- `extensions/vscode/package-lock.json`: `"version": "0.4.0"`
- `CHANGELOG.md`: Add `## [0.4.0] - 2026-08-28` section documenting path resolution, auto-init UX, placeholder checks, and CLI list `--json`.
- `docs/vscode-extension.md`: Document `terralings.initExercises` and `terralings.openNextExercise`.

- [ ] **Step 2: Run all local checks**

Run:
```bash
make all
make docs-build
```
Expected: All linters, unit tests, race tests, VSIX packaging, and documentation builds PASS.

- [ ] **Step 3: Commit, branch, PR, and release verification**

Create branch `feat/path-resolution-auto-init`, commit, push, create PR, verify all 8 CI checks pass, merge to `main`, and verify automated `release.yaml` workflow publishes `v0.4.0` to GitHub Releases and VS Code Marketplace.
