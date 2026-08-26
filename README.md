# Terralings

> An interactive, terminal-driven learning environment for mastering Terraform and OpenTofu through hands-on exercises.

`terralings` guides you through fixing broken configurations, writing declarative Infrastructure as Code (IaC), mastering HCL expressions and built-in functions, refactoring state with `moved` blocks, authoring `.tftest.hcl` unit and integration tests, and configuring OpenTofu advanced features such as state encryption and early variable evaluation.

Inspired by [`rustlings`](https://github.com/rust-lang/rustlings), [`ziglings`](https://github.com/ziglings/exercises), and [`raylings`](https://github.com/ray-project/raylings).

---

## Features

- **Interactive Watch Mode (`terralings watch`)**: Automatically monitors your exercise files via `fsnotify` and re-evaluates and validates in real time on every file save.
- **Dual Engine Support**: Seamlessly detects and runs against either **OpenTofu** (`tofu`) or **Terraform** (`terraform`) (version >= 1.6.0).
- **Sub-100ms Evaluation**: Shared provider plugin caching eliminates redundant network downloads during provider initialization.
- **Comprehensive 13-Chapter Curriculum**: 56 progressive exercises covering primitives, variables, collections, functions, meta-arguments, dynamic blocks, data sources, modules, state refactoring, native testing, production patterns, OpenTofu extensions, and policy governance.
- **Progressive Hints (`terralings hint`)**: Multi-level contextual guidance to nudge you forward when stuck without spoiling the answer.
- **Built-in Verification (`terralings verify`)**: Complete curriculum progress dashboard tracking your status across all 13 chapters.
- **Zero Heavy Cloud Dependencies**: All exercises execute purely locally using standard built-in providers (`local`, `random`, `archive`, `terraform_data`, `test_mock`) or fast in-memory execution—no AWS/GCP/Azure credentials required.

---

## Prerequisites

Before using `terralings`, ensure you have installed:

1. **Go** (version 1.22 or higher) — required for compiling `terralings` from source:
   ```bash
   go version
   ```
2. **OpenTofu** (>= 1.6.0) or **Terraform** (>= 1.6.0):
   ```bash
   tofu version
   # or
   terraform version
   ```

---

## Installation & Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/dnf0/terralings.git
cd terralings
```

### 2. Build the Binary

Build the executable binary into `./bin/terralings`:

```bash
make build
```

Or install directly to `$GOPATH/bin`:

```bash
go install ./cmd/terralings
```

### 3. Verify Installation

Check that `terralings` successfully discovers your OpenTofu or Terraform binary:

```bash
./bin/terralings version
```

Example output:
```text
terralings v0.1.0
Detected binary: /usr/local/bin/tofu (OpenTofu v1.6.2)
```

---

## How It Works

Each exercise file contains a deliberate syntax, type, expression, or logic error, marked with a comment header:

```hcl
// I AM NOT DONE
// Exercise: primitives01
// Fix the configuration below so that 'tofu validate' succeeds.
```

### The Learning Loop

1. **Start Watch Mode**: Run `./bin/terralings watch`.
2. **Read the Error Output**: `terralings` compiles the current exercise in an isolated workspace and prints the compiler / runtime diagnostic error.
3. **Edit the File**: Open the file in your preferred text editor (e.g., `exercises/01_primitives/primitives01.tf`), fix the broken configuration or implement the missing resources.
4. **Remove the Marker**: Delete the `// I AM NOT DONE` or `# I AM NOT DONE` line at the top of the file.
5. **Save and Advance**: As soon as you save the file, `terralings` detects the change, verifies the solution, displays a success notification, and immediately transitions to the next exercise.

---

## CLI Commands Reference

| Command | Description |
|---|---|
| `terralings watch` | Start interactive watch mode (automatically evaluates on file save) |
| `terralings run <exercise>` | Run verification against a single exercise (e.g. `terralings run primitives01`) |
| `terralings hint <exercise>` | Display progressive hint(s) for the specified exercise |
| `terralings list` | List all chapters and exercises with status indicators |
| `terralings verify` | Run sequential evaluation across the entire curriculum and display progress |
| `terralings version` | Print the Terralings CLI version and detected IaC binary |

### Command Examples & Sample Outputs

#### 1. Interactive Watch Mode
```bash
terralings watch
```
```text
======================================================================
  TERRALINGS: Continuous Watch Mode
  Watching directory: exercises
  Press Ctrl+C to exit
======================================================================

⌛ primitives01 still contains 'I AM NOT DONE' marker. Keep going!

Success! The configuration is valid.
```

#### 2. Progressive Hints
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

#### 3. Curriculum Overview
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

#### 4. Progress Verification
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
