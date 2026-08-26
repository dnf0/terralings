# Terralings: An Interactive Hands-On Terraform & OpenTofu Learning Environment

**Date:** 2026-08-26  
**Status:** Approved  
**Target Repository:** `terralings`  

---

## 1. Executive Summary & Vision

`terralings` is an interactive, terminal-driven educational tool and comprehensive hands-on curriculum designed to teach HashiCorp Terraform and OpenTofu from the ground up—inspired by the pedagogy and developer experience of `rustlings`, `ziglings`, and `raylings`.

Learners progress through progressively challenging exercises where they fix broken code, fill in missing infrastructure declarations, master complex HCL expressions and functions, write modular architectures, refactor state with `moved` blocks, run unit tests using `.tftest.hcl`, and explore modern OpenTofu state encryption.

---

## 2. System Architecture

```
                                +-----------------------------+
                                |      Learner Terminal       |
                                |       terralings watch      |
                                +--------------+--------------+
                                               |
                                               v
+------------------------------------------------------------------------------------------+
|                                  Terralings CLI Engine (Go)                              |
|                                                                                          |
|  +--------------------+   +-----------------------+   +-------------------------------+  |
|  | File Watcher       |   | Manifest & Registry   |   | Subprocess Runner & Evaluator |  |
|  | (fsnotify)         |-->| (exercise catalogue & |-->| (tofu/terraform detection,    |  |
|  |                    |   |  marker parser)       |   |  init, validate, plan, test)  |  |
|  +--------------------+   +-----------------------+   +---------------+---------------+  |
|                                                                       |                  |
+-----------------------------------------------------------------------|------------------+
                                                                        |
                                       +--------------------------------+
                                       |
                                       v
+------------------------------------------------------------------------------------------+
|                               Runtime & Provider Caching Layer                           |
|                                                                                          |
|  +---------------------------------------------+  +-----------------------------------+  |
|  | Shared Local Provider Cache                 |  | Offline Mocking & HCL Test Engine |  |
|  | - TF_PLUGIN_CACHE_DIR (~/.terralings/cache) |  | - OpenTofu / Terraform test suite |  |
|  | - Instant sub-100ms initialization          |  | - mock_provider & override blocks |  |
|  +---------------------------------------------+  +-----------------------------------+  |
|                                                                                          |
+------------------------------------------------------------------------------------------+
```

### 2.1 Core Components

1. **CLI Engine (`cmd/terralings/main.go`, `internal/ui`, `internal/runner`, `internal/watcher`)**:
   - Written in idiomatic Go (1.22+).
   - Built on `spf13/cobra` and `charmbracelet/lipgloss` / `bubbletea`.
   - Manages interactive command modes (`watch`, `run`, `test`, `hint`, `list`, `verify`).
   - Renders color-coded terminal diagnostics, syntax highlights, progress bars, and ASCII completion banners.

2. **Binary Detector (`internal/detector`)**:
   - Searches `$PATH` for `tofu` (OpenTofu) first; falls back to `terraform` if `tofu` is not found.
   - Respects `TERRALINGS_BIN` environment variable or `--bin` CLI flag for explicit overrides.
   - Validates minimum required version (>= 1.6.0).

3. **Provider Cache & Fast Initialization (`internal/runner`)**:
   - Manages a shared local plugin cache directory (`~/.terralings/plugin-cache` or project `.cache/`) via `TF_PLUGIN_CACHE_DIR`.
   - Guarantees sub-100ms `init` execution time per exercise without redundant network roundtrips.

4. **Exercise Manifest & Registry (`internal/manifest`)**:
   - Declarative catalogue of all 13 chapters and 50+ exercises with titles, paths, prerequisites, and progressive hints.
   - Parses `# I AM NOT DONE` / `// I AM NOT DONE` markers at the top of `.tf` files.

5. **Canonical Solutions & Testing Harness (`solutions/`, `test/`)**:
   - `solutions/` mirrors every exercise file in `exercises/` with fully functional, verified implementations.
   - Automated Go test suite validates that:
     - All solutions pass `tofu validate` and `tofu test` / `tofu plan`.
     - All starter exercises in `exercises/` fail before user modification (preventing accidental no-op exercises).
     - The CLI runner and watcher components function correctly.

---

## 3. Curriculum & Syllabus Specification

The curriculum is divided into 13 numbered chapters spanning fundamental syntax to enterprise patterns and OpenTofu features:

### Chapter 1: HCL Foundations & Core Primitives
- `exercises/01_primitives/primitives01.tf`: Terraform configuration block, `required_version`, and `required_providers`.
- `exercises/01_primitives/primitives02.tf`: Declaring first resource (`local_file` / `terraform_data`).
- `exercises/01_primitives/primitives03.tf`: Resource dependencies, implicit references, and dependency graphs.
- `exercises/01_primitives/primitives04.tf`: String interpolation (`${...}`) and heredoc multi-line strings (`<<-EOT`).
- `exercises/01_primitives/primitives05.tf`: Comments, identifier naming rules, and canonical code formatting (`tofu fmt`).
- `exercises/01_primitives/primitives06.tf`: The declarative lifecycle: understanding `validate`, `plan`, and `apply`.

### Chapter 2: Input Variables, Types & Validations
- `exercises/02_variables/variables01.tf`: Primitive variable declarations (`string`, `number`, `bool`).
- `exercises/02_variables/variables02.tf`: Collection types: `list(string)`, `map(string)`, and `set(string)`.
- `exercises/02_variables/variables03.tf`: Structural types: `object({...})` and `tuple([...])`.
- `exercises/02_variables/variables04.tf`: Variable defaults, `nullable = false`, and `.tfvars` priority order.
- `exercises/02_variables/variables05.tf`: Custom variable validation blocks with conditions and `error_message`.

### Chapter 3: Outputs, Locals & Expressions
- `exercises/03_outputs_locals/outputs01.tf`: Defining outputs, descriptions, and `sensitive = true` redaction.
- `exercises/03_outputs_locals/locals01.tf`: Declaring `locals` blocks for DRY intermediate calculations.
- `exercises/03_outputs_locals/expr01.tf`: Ternary conditional expressions (`condition ? true_val : false_val`).
- `exercises/03_outputs_locals/expr02.tf`: Splat expressions (`[*]`) and legacy attribute referencing.

### Chapter 4: Built-in Functions & Collections
- `exercises/04_functions/func01.tf`: String manipulation functions (`format`, `join`, `split`, `replace`, `lower`).
- `exercises/04_functions/func02.tf`: Collection operations (`merge`, `lookup`, `distinct`, `slice`, `flatten`, `zipmap`).
- `exercises/04_functions/func03.tf`: Encodings: `jsonencode`, `yamlencode`, `base64encode`, and `jsondecode`.
- `exercises/04_functions/func04.tf`: Filesystem functions: `file()`, `templatefile()`, `fileset()`, and `fileexists()`.
- `exercises/04_functions/func05.tf`: Safe evaluation expressions: `can()`, `try()`, and fallback defaults.

### Chapter 5: Meta-Arguments & Resource Scaling
- `exercises/05_meta_arguments/meta01.tf`: Scaling resources with `count` and indexing with `count.index`.
- `exercises/05_meta_arguments/meta02.tf`: Idempotent resource mapping with `for_each` over sets and maps (`each.key`, `each.value`).
- `exercises/05_meta_arguments/meta03.tf`: Explicit dependency ordering with `depends_on`.
- `exercises/05_meta_arguments/meta04.tf`: Resource lifecycle blocks: `create_before_destroy` and `prevent_destroy`.
- `exercises/05_meta_arguments/meta05.tf`: Dynamic drift handling with lifecycle `ignore_changes`.

### Chapter 6: Dynamic Blocks & Advanced HCL
- `exercises/06_dynamic_blocks/dynamic01.tf`: Basic `dynamic` block iteration over nested maps.
- `exercises/06_dynamic_blocks/dynamic02.tf`: Dynamic blocks with custom `iterator` naming.
- `exercises/06_dynamic_blocks/dynamic03.tf`: Multi-level nested dynamic blocks.
- `exercises/06_dynamic_blocks/dynamic04.tf`: Conditional dynamic block emission using empty list/map tricks.

### Chapter 7: Data Sources & State Querying
- `exercises/07_data_sources/data01.tf`: Querying local filesystem data sources and environment attributes.
- `exercises/07_data_sources/data02.tf`: Creating and inspecting compressed bundles with `archive_file`.
- `exercises/07_data_sources/data03.tf`: External data source queries and structured JSON contracts.
- `exercises/07_data_sources/data04.tf`: Custom preconditions and postconditions on data sources and resources.

### Chapter 8: Modular Infrastructure Architecture
- `exercises/08_modules/module01/`: Building a clean child module (inputs, internal encapsulation, outputs).
- `exercises/08_modules/module02/`: Calling local child modules and passing variables.
- `exercises/08_modules/module03/`: Module composition and multi-instance module deployment with `for_each`.
- `exercises/08_modules/module04/`: Passing provider configurations and aliases (`providers = { ... }`).
- `exercises/08_modules/module05/`: Submodule boundaries and avoiding deep nesting antipatterns.

### Chapter 9: State Management & Refactoring Surgery
- `exercises/09_state_refactoring/state01.tf`: Declarative refactoring using `moved` blocks (renaming resources).
- `exercises/09_state_refactoring/state02.tf`: Migrating resources from root module into child modules via `moved`.
- `exercises/09_state_refactoring/state03.tf`: Declarative `import` blocks for onboarding unmanaged resources into state.
- `exercises/09_state_refactoring/state04.tf`: Controlled resource replacement with `replace_triggered_by`.

### Chapter 10: Native Unit & Integration Testing (`.tftest.hcl`)
- `exercises/10_testing/test01.tftest.hcl`: Basic test assertions with `run` blocks and `command = plan`.
- `exercises/10_testing/test02.tftest.hcl`: Validating applied local resources in `command = apply` test runs.
- `exercises/10_testing/test03.tftest.hcl`: Mocking cloud providers and resources (`mock_provider "aws" { ... }`).
- `exercises/10_testing/test04.tftest.hcl`: Testing module failure cases and error assertions with `expect_failures`.

### Chapter 11: Production Patterns & Anti-Patterns
- `exercises/11_patterns/pattern01.tf`: Decomposing monolithic root modules into domain-bounded layers.
- `exercises/11_patterns/pattern02.tf`: Minimizing blast radius and side-effects with `terraform_data`.
- `exercises/11_patterns/pattern03.tf`: Managing secrets and sensitive values cleanly with `ephemeral` variables.
- `exercises/11_patterns/pattern04.tf`: Advanced list and map comprehensions (`for` expressions with filtering).

### Chapter 12: OpenTofu Innovations & Enterprise Features
- `exercises/12_opentofu/tofu01.tf`: State encryption at rest with `key_provider` and `method` blocks.
- `exercises/12_opentofu/tofu02.tf`: Early variable evaluation in backend and provider configurations.
- `exercises/12_opentofu/tofu03.tf`: Configuring and consuming the OpenTofu public provider registry.

### Chapter 13: Linting, CI/CD & Governance
- `exercises/13_governance/gov01.tf`: Strict formatting rules (`tofu fmt -check`) and canonical layout.
- `exercises/13_governance/gov02.tf`: Policy assertions and custom validation rules across environments.
- `exercises/13_governance/gov03.tf`: Structuring pull-request automated validation pipelines for CI.

---

## 4. Exercise & Solution Specification

Every exercise file follows a strict, consistent convention:

```terraform
# Exercise: exercises/01_primitives/primitives02.tf
# Topic: First Resource Declaration (local_file)
#
# Instructions:
# Fix the resource block below so it creates a file named "welcome.txt"
# with the content "Welcome to Terralings!".

# I AM NOT DONE

terraform {
  required_version = ">= 1.6.0"
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5"
    }
  }
}

resource "local_file" "welcome" {
  # FIX ME: Add filename and content
}
```

Every exercise in `exercises/` has a matching file in `solutions/` containing the canonical solution:

```terraform
# Solution: solutions/01_primitives/primitives02.tf
# Topic: First Resource Declaration (local_file)

terraform {
  required_version = ">= 1.6.0"
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5"
    }
  }
}

resource "local_file" "welcome" {
  filename = "${path.module}/welcome.txt"
  content  = "Welcome to Terralings!"
}
```

---

## 5. Offline-First Provider & Testing Strategy

To ensure zero cost and frictionless setup for learners:
1. **Local Providers**: Core exercises utilize `hashicorp/local`, `hashicorp/random`, `hashicorp/tls`, `hashicorp/archive`, `hashicorp/time`, and the built-in `terraform_data` resource.
2. **Mocking Cloud Resources**: For cloud-oriented chapters (AWS/GCP patterns), exercises use OpenTofu/Terraform's native `mock_provider` and `override_resource` blocks inside `.tftest.hcl` files.
3. **No Cloud Credentials Needed**: 100% of the curriculum runs completely offline on the learner's local workstation.

---

## 6. Repository Layout & CI/CD

```
terralings/
├── .github/workflows/ci.yml   # GitHub Actions multi-OS / binary matrix (tofu & terraform)
├── .gitignore                 # Excludes .agents/, .terraform/, state files, plugin caches
├── go.mod                     # Go 1.22+ module definition
├── go.sum                     # Dependency checksums
├── Makefile                   # build, test, lint, verify targets
├── README.md                  # Comprehensive documentation, quickstart & curriculum map
├── CONTRIBUTING.md            # Guidelines for adding exercises & local development
├── CHANGELOG.md               # Version history following keepachangelog standard
├── LICENSE                    # Apache 2.0 license
├── cmd/
│   └── terralings/
│       └── main.go            # Entry point parsing CLI flags & subcommands
├── internal/
│   ├── detector/              # Discovers and validates tofu / terraform on PATH
│   ├── manifest/              # Registry of all 13 chapters, 50+ exercises & hints
│   ├── models/                # Domain models (Chapter, Exercise, Status, RunResult)
│   ├── runner/                # Subprocess executor with TF_PLUGIN_CACHE_DIR
│   ├── ui/                    # Lipgloss formatting, banners, diffs, progress tables
│   └── watcher/               # Debounced fsnotify file watcher loop
├── exercises/                 # 13 chapters of starter exercises
├── solutions/                 # 13 chapters of verified working reference solutions
└── test/
    ├── detector_test.go
    ├── manifest_test.go
    ├── runner_test.go
    ├── cli_test.go
    └── solutions_test.go      # Validates all solutions pass & starter exercises fail
```

---

## 7. Verification and Acceptance Criteria

1. **CLI Commands**:
   - `terralings watch`: Interactive watcher detecting edits and re-evaluating in <100ms.
   - `terralings run <name>`: Runs validation/test on a single exercise or solution.
   - `terralings hint <name>`: Displays progressive hints for an exercise.
   - `terralings test`: Runs the automated test harness across all reference solutions.
   - `terralings list`: Displays all chapters, exercises, and completion progress.
   - `terralings verify`: Full curriculum verification suite.
2. **Curriculum Completeness**:
   - 13 chapters, 50+ exercises covering HCL syntax to advanced OpenTofu state encryption and `.tftest.hcl` testing.
   - 100% of solutions pass verification.
   - 100% of starter exercises fail before user solution.
3. **Continuous Integration**:
   - Automated GitHub Actions workflow running on Ubuntu and macOS across both `tofu` and `terraform`.
