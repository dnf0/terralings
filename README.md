# Terralings

> An interactive, terminal-driven learning environment for mastering Terraform and OpenTofu through hands-on exercises.

`terralings` guides you through fixing broken configurations, writing declarative infrastructure as code (IaC), mastering HCL expressions and built-in functions, refactoring state with `moved` blocks, authoring `.tftest.hcl` unit tests, and configuring OpenTofu features such as state encryption.

Inspired by `rustlings`, `ziglings`, and `raylings`.

---

## Features

- **Interactive Watch Mode (`terralings watch`)**: Automatically monitors your exercise files via `fsnotify` and recompiles/validates on every save.
- **Dual Engine Support**: Seamlessly detects and runs against either **OpenTofu** (`tofu`) or **Terraform** (`terraform`) (version >= 1.6.0).
- **Sub-100ms Evaluation**: Shared provider plugin caching eliminates redundant network roundtrips during `init`.
- **Comprehensive 13-Chapter Curriculum**: 50+ exercises covering primitives, variables, collections, functions, meta-arguments, dynamic blocks, data sources, modules, state refactoring, testing, design patterns, OpenTofu extensions, and policy governance.
- **Progressive Hints (`terralings hint`)**: Multi-level contextual guidance when you get stuck.
- **Built-in Solutions & Verification**: Verify entire exercise suites or single exercises with automated grading.

---

## Prerequisites

Before using `terralings`, ensure you have installed:

1. **Go** (version 1.22 or higher) — for compiling `terralings` from source:
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

### 1. Clone the repository

```bash
git clone https://github.com/dnf0/terralings.git
cd terralings
```

### 2. Build the binary

```bash
make build
```

Or install directly to `$GOPATH/bin`:

```bash
go install ./cmd/terralings
```

### 3. Start Learning

Start watch mode to begin with the first exercise:

```bash
terralings watch
```

`terralings` will display the error for the current exercise. Open the indicated file in your editor, fix the issue, and remove the `// I AM NOT DONE` or `# I AM NOT DONE` marker at the top of the file. As soon as you save the file, `terralings` will verify your solution and advance to the next exercise!

---

## CLI Commands

| Command | Description |
|---|---|
| `terralings watch` | Run interactive watcher mode (auto-validates on file save) |
| `terralings run <exercise>` | Run and validate a specific exercise (e.g. `terralings run primitives01`) |
| `terralings hint [exercise]` | Print progressive hints for the current or specified exercise |
| `terralings list` | List all chapters and exercises along with completion status |
| `terralings verify` | Run validation against all exercises sequentially |
| `terralings version` | Print the terralings CLI version and detected IaC binary |

### Environment Variables & Flags

- `TERRALINGS_BIN`: Explicitly specify the binary path (e.g., `TERRALINGS_BIN=/usr/local/bin/tofu terralings watch`).
- `--bin <path>`: CLI flag override for OpenTofu/Terraform binary.

---

## Curriculum Map

The curriculum spans 13 structured chapters:

1. **HCL Foundations & Core Primitives** (`exercises/01_primitives/`): `terraform` block, providers, resources, dependencies, strings, heredocs, and syntax formatting.
2. **Input Variables, Types & Validations** (`exercises/02_variables/`): Primitives, collections (`list`, `map`, `set`), structural types (`object`, `tuple`), defaults, and custom `validation` rules.
3. **Outputs, Locals & Expressions** (`exercises/03_outputs_locals/`): Output values, `sensitive` markers, `locals`, ternary expressions, and splat operators (`[*]`).
4. **Built-in Functions & Collections** (`exercises/04_functions/`): String formatting, map merging, JSON/YAML encoding, filesystem functions, `can()`, and `try()`.
5. **Meta-Arguments & Scaling** (`exercises/05_meta_arguments/`): `count`, `for_each`, explicit `depends_on`, and `lifecycle` blocks (`create_before_destroy`, `prevent_destroy`, `ignore_changes`).
6. **Dynamic Blocks & Advanced HCL** (`exercises/06_dynamic_blocks/`): Iterating dynamic nested blocks, custom iterators, nested dynamic blocks, and conditional emission.
7. **Data Sources & State Querying** (`exercises/07_data_sources/`): Local filesystem queries, archive bundling, external data sources, and custom pre/postconditions.
8. **Modular Infrastructure Architecture** (`exercises/08_modules/`): Child module encapsulation, parameters, composition, module `for_each`, and provider aliases.
9. **State Refactoring & Surgery** (`exercises/09_state_refactoring/`): Declarative refactoring with `moved` blocks, resource migration, declarative `import`, and `replace_triggered_by`.
10. **Native Unit & Integration Testing** (`exercises/10_testing/`): Native `.tftest.hcl` files, `plan` vs `apply` test runs, `mock_provider`, and `expect_failures`.
11. **Production Infrastructure Patterns** (`exercises/11_patterns/`): Multi-tier composition, zero-downtime blue/green patterns, and tag inheritance.
12. **OpenTofu Advanced Features** (`exercises/12_opentofu/`): OpenTofu state encryption (`key_provider`), early variable evaluation, and OpenTofu registry.
13. **Governance & Security Policies** (`exercises/13_governance/`): Enforcing security guardrails, naming conventions, and compliance rules.

---

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, development workflow, and pull request guidelines.

---

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
