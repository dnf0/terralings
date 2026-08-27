# Contributing to Terralings

Thank you for your interest in contributing to Terralings! We welcome contributions ranging from typo fixes and documentation updates to new exercises, tests, CLI enhancements, and VS Code extension features.

---

## Development Setup

### Prerequisites

- **Go**: 1.22 or higher
  ```bash
  go version
  ```
- **OpenTofu**: &ge; 1.6.0 *(Recommended)* or **Terraform**: &ge; 1.5.0
  ```bash
  tofu version || terraform version
  ```
- **Node.js & npm** *(for VS Code companion extension development)*:
  - Node.js &ge; 18.0.0
  - npm &ge; 9.0.0
  ```bash
  node -v && npm -v
  ```
- **Make**: Standard build automation utility.
- **Git**

### Clone & Build

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
| `make check` | Runs `go mod verify`, `gofmt` check, solutions format check, and `go vet` |
| `make build` | Compiles binary to `bin/terralings` (with macOS codesigning support) |
| `make test` | Runs Go test suite |
| `make test-race` | Runs Go test suite with the race detector enabled (`-race`) |
| `make fmt-check` | Verifies Go source files are formatted with `gofmt` |
| `make fmt-solutions` | Verifies HCL solution files are formatted with `tofu fmt` / `terraform fmt` |
| `make extension-check` | Type-checks, compiles, and tests the VS Code companion extension |
| `make extension-package` | Bundles and packages the VS Code extension into `dist/terralings-vscode.vsix` |
| `make docs-serve` | Starts local MkDocs documentation development server on `localhost:8000` |
| `make docs-build` | Builds static documentation site with strict verification |
| `make clean` | Removes compiled binaries, temporary test state, and cache directories |

---

## Project Structure

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
- Place the exercise in the appropriate chapter directory.
- Include clear `# TODO:` instructions explaining what the learner needs to accomplish.
- Ensure the starter code fails validation (`ModeValidate`), planning (`ModePlan`), or testing (`ModeTest`) in its initial state.

```hcl
// ============================================================================
// Exercise: myexercise01
// Chapter:  01_primitives
// ============================================================================

terraform {
  # TODO: Declare required_version and required_providers
}
```

### Step 2: Create Canonical Solution (`solutions/<chapter>/<name>.tf`)
- Create a corresponding file in `solutions/<chapter>/<name>.tf` with the identical filename.
- Provide the complete, working, canonical HCL code.
- Ensure the solution is properly formatted with `tofu fmt` or `terraform fmt`.

### Step 3: Register in Manifest (`internal/manifest/manifest.go`)
- Register the exercise in `internal/manifest/manifest.go` with its ID, title, chapter, path, mode (`ModeValidate`, `ModePlan`, or `ModeTest`), and progressive hints:

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
- Run `make all` (or `make test-race`) to confirm that the starter code fails and the solution passes cleanly.

---

## Code Formatting & Quality Standards

- **Go Code**: All Go files must be formatted with `gofmt`. Run `make fmt-check` before committing.
- **HCL Code**: All `.tf` and `.tftest.hcl` files in `solutions/` must be canonically formatted. Run `make fmt-solutions`.
- **TypeScript Code**: Extension code in `extensions/vscode/` must pass `npm run check-types` without errors.

---

## Commit Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification:

```
<type>(<optional scope>): <description>

[optional body]

[optional footer]
```

### Common Types

- `feat:` A new feature, command, or exercise
- `fix:` A bug fix in an exercise, runner, LSP server, TUI, or CLI
- `docs:` Documentation additions or improvements
- `test:` Adding or updating tests
- `refactor:` Code refactoring without behavioral changes
- `chore:` Build scripts, dependencies, or toolchain updates

---

## Pull Request Guidelines

1. Fork the repository and create a descriptive feature branch (`feat/my-feature` or `fix/issue-description`).
2. Run `make all` locally to ensure all formatting, linting, builds, race-detection tests, and extension checks pass cleanly.
3. Submit a pull request with a clear description of the changes, testing evidence, and related issue links.
