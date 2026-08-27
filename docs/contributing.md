# Contributing to Terralings

Thank you for your interest in contributing to Terralings! Whether you are fixing a typo, improving exercise instructions, adding new curriculum chapters, or enhancing the CLI or VS Code extension, your help is warmly welcomed.

---

## Development Prerequisites

Ensure the following development dependencies are installed on your machine:

1. **Go Toolchain**: &ge; 1.22
   ```bash
   go version
   ```
2. **OpenTofu or Terraform**:
   - **OpenTofu**: &ge; 1.6.0 *(Recommended)*
   - **Terraform**: &ge; 1.5.0
   ```bash
   tofu version || terraform version
   ```
3. **Node.js & npm** *(for VS Code extension development)*:
   - Node.js &ge; 18.0.0
   - npm &ge; 9.0.0
   ```bash
   node -v && npm -v
   ```
4. **Make**: Standard build automation utility.

---

## Local Setup & Quickstart

Clone the repository and run the complete check and build pipeline:

```bash
git clone https://github.com/dnf0/terralings.git
cd terralings

# Run all formatting, verification, build, and race-detector checks:
make all
```

### Common Makefile Targets

| Target | Description |
|---|---|
| `make all` | Runs `make check`, `make build`, `make test-race`, `make extension-check`, and `make extension-package` |
| `make build` | Compiles the Go binary into `bin/terralings` (with macOS codesigning) |
| `make test` | Executes the complete Go test suite |
| `make test-race` | Executes Go test suite with the race detector enabled (`-race`) |
| `make check` | Runs `go mod verify`, `gofmt` checks, solutions format checks, and `go vet` |
| `make fmt-check` | Verifies all Go source files are formatted with `gofmt` |
| `make fmt-solutions` | Verifies all HCL solution files are formatted with `tofu fmt` / `terraform fmt` |
| `make extension-check` | Compiles TypeScript, checks types, and runs VS Code extension tests |
| `make extension-package` | Bundles and packages the VS Code extension into `dist/terralings-vscode.vsix` |
| `make docs-serve` | Starts local MkDocs documentation development server on `localhost:8000` |
| `make docs-build` | Builds the static documentation site with `--strict` verification |
| `make clean` | Removes compiled binaries, temporary test state, and cache directories |

---

## Project Repository Structure

```
terralings/
├── cmd/
│   └── terralings/         # CLI main entrypoint and Cobra commands
├── internal/
│   ├── detector/           # OpenTofu/Terraform binary auto-detection
│   ├── doctor/             # Pre-flight diagnostic probe implementation
│   ├── lsp/                # Language Server Protocol daemon & handlers
│   ├── manifest/           # Curriculum manifest singleton & exercises map
│   ├── models/             # Core data structures (Exercise, Chapter, Mode)
│   ├── runner/             # Isolated exercise execution & validation engine
│   ├── search/             # Curriculum full-text search engine
│   ├── state/              # Progress persistence store (.terralings/state.json)
│   ├── tour/               # 5-step guided terminal tour
│   ├── tui/                # Bubble Tea full-screen terminal dashboard
│   ├── ui/                 # Lip Gloss terminal styling & formatting components
│   └── watcher/            # fsnotify file watcher engine & event loops
├── exercises/              # Embedded student-facing starter exercises
├── solutions/              # Canonical passing reference solutions (1:1 with exercises)
├── extensions/
│   └── vscode/             # Official VS Code Companion Extension (TypeScript)
├── docs/                   # MkDocs documentation site sources
└── test/                   # Comprehensive end-to-end and integration tests
```

---

## Authoring New Exercises

When adding a new exercise to the curriculum, follow this four-step workflow:

### Step 1: Create Starter Exercise (`exercises/<chapter>/<name>.tf`)

- Place the exercise file in the appropriate chapter directory.
- Include the `// I AM NOT DONE` or `# I AM NOT DONE` comment marker at the top of the file.
- Add clear `# TODO:` comments explaining what the learner needs to accomplish.
- Ensure the starter code fails validation (`ModeValidate`), planning (`ModePlan`), or testing (`ModeTest`) in its initial state.

```hcl
// I AM NOT DONE
// ============================================================================
// Exercise: myexercise01
// Chapter:  01_primitives
// ============================================================================

terraform {
  # TODO: Declare required_version and required_providers
}
```

### Step 2: Create Canonical Solution (`solutions/<chapter>/<name>.tf`)

- Create a corresponding file in `solutions/<chapter>/<name>.tf` with identical filename.
- Provide the complete, working, canonical HCL code with the completion marker removed.
- Verify the solution passes `tofu fmt` / `terraform fmt`.

### Step 3: Register in Manifest (`internal/manifest/manifest.go`)

Register the exercise inside the appropriate chapter in `internal/manifest/manifest.go`:

```go
{
    Name:        "myexercise01",
    Title:       "Descriptive Human-Readable Title",
    Path:        "exercises/01_primitives/myexercise01.tf",
    ChapterName: "01_primitives",
    Hints: []string{
        "First progressive conceptual hint.",
        "Second progressive syntax hint.",
    },
    Mode: models.ModePlan, // models.ModeValidate, models.ModePlan, or models.ModeTest
},
```

### Step 4: Run Verification Tests

Run the full verification suite to confirm that:
1. The new exercise starts in a failing state.
2. The reference solution passes validation/planning/testing cleanly.
3. The manifest and filesystem stay strictly synchronized.

```bash
make check
make test
```

---

## Code Formatting & Style Standards

- **Go Code**: All Go files must be formatted with `gofmt`. Run `make fmt-check` before committing.
- **HCL Code**: All `.tf` and `.tftest.hcl` files in `solutions/` must be canonically formatted with `tofu fmt` or `terraform fmt`. Run `make fmt-solutions`.
- **TypeScript Code**: Extension code in `extensions/vscode/` must pass `npm run check-types` without errors or warnings.

---

## Commit Guidelines

Terralings enforces the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification:

```
<type>(<optional scope>): <description>

[optional body]

[optional footer]
```

### Common Types

- `feat`: A new exercise, CLI feature, or extension capability
- `fix`: A bug fix in an exercise, runner, LSP server, or TUI
- `docs`: Documentation additions or improvements
- `test`: Adding or updating test cases
- `refactor`: Code refactoring without behavioral changes
- `chore`: Build scripts, dependencies, or toolchain updates

### Examples

- `feat(manifest): add new exercise for state encryption`
- `fix(runner): handle trailing whitespace in plan evaluation`
- `docs(syllabus): document chapter 13 architecture governance`

---

## Pull Request Guidelines

1. **Create a Feature Branch**: Always create a feature branch off `main` (e.g., `feat/new-exercise` or `fix/lsp-hover-issue`).
2. **Atomic Commits**: Keep commits logical, focused, and well-described.
3. **Pass All Quality Gates**: Ensure `make all` passes locally before submitting your pull request.
4. **Detailed PR Description**: Outline the problem solved, changes made, and testing evidence.
