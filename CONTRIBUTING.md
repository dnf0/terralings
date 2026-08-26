# Contributing to Terralings

Thank you for your interest in contributing to Terralings! We welcome contributions ranging from typo fixes and documentation updates to new exercises, tests, and CLI enhancements.

---

## Development Setup

### Prerequisites

- **Go**: 1.22 or higher
- **OpenTofu**: >= 1.6.0 (or **Terraform** >= 1.6.0)
- **Make**: (optional, for convenience targets)
- **Git**

### Clone & Build

```bash
git clone https://github.com/dnf0/terralings.git
cd terralings

# Run all formatting, verification, build, and race-detector checks
make all

# Or run individual targets:
make check      # Dependency verification, gofmt check, and go vet
make build      # Build binary to bin/terralings
make build-all  # Build all packages
make test       # Run test suite
make test-race  # Run test suite with race detector
```

---

## Project Structure

- `cmd/terralings/`: Entry point for the CLI application.
- `internal/`: Core engine packages (`detector`, `models`, `manifest`, `runner`, `ui`, `watcher`).
- `exercises/`: Learner-facing exercise files organized by chapter. Each exercise starts with a `// I AM NOT DONE` or `# I AM NOT DONE` marker and contains intentional errors or missing declarations.
- `solutions/`: Canonical, passing implementations corresponding 1:1 with files in `exercises/`.
- `test/`: Integration and unit test suite verifying manifest integrity, solutions validity, CLI runner, and watcher behavior.

---

## Authoring New Exercises

When adding a new exercise or updating an existing one:

1. **Exercise File (`exercises/<chapter>/<name>.tf`)**:
   - Must include the `// I AM NOT DONE` or `# I AM NOT DONE` comment marker at the top.
   - Should focus on a single concept, with clear explanatory comments guiding the learner.
   - Must fail initial validation (`tofu validate` or `tofu test`) in its default state.

2. **Solution File (`solutions/<chapter>/<name>.tf`)**:
   - Must contain the working, canonical HCL code with the completion marker removed.
   - Must pass `tofu validate` (or `tofu test`) cleanly.

3. **Manifest Entry (`internal/manifest/manifest.go`)**:
   - Register the exercise with its ID, title, chapter, file path, mode (`validate`, `test`, `plan`), and progressive hints.

4. **Verify Solutions**:
   - Run `make all` (or `make test`) to ensure all tests pass.

---

## Commit Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` A new feature or exercise
- `fix:` A bug fix in an exercise, runner, or CLI
- `docs:` Documentation changes
- `test:` Adding or updating tests
- `refactor:` Code refactoring without behavioral change
- `chore:` Maintenance, dependencies, build configurations

---

## Pull Request Guidelines

1. Fork the repository and create a descriptive feature branch (`feat/my-feature` or `fix/issue-description`).
2. Run `make all` to ensure all linting, formatting, builds, and tests pass cleanly.
3. Submit a pull request with a clear description of the changes and testing evidence.
