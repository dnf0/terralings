# Curriculum Syllabus

Terralings features a comprehensive 13-chapter curriculum comprising **56 progressive exercises**. The curriculum covers the entire spectrum of modern Infrastructure as Code (IaC)—from foundational syntax and variables to advanced state surgery, native `.tftest.hcl` testing, OpenTofu state encryption, and architectural governance.

---

## Curriculum Overview Matrix

| # | Chapter Name | Focus Area | Exercises | Verification Modes |
|---|---|---|:---:|:---:|
| 01 | [HCL Foundations & Core Primitives](#chapter-01-hcl-foundations-core-primitives) | Blocks, providers, resources, heredocs, lifecycle | 6 | `validate`, `plan` |
| 02 | [Input Variables, Types & Validations](#chapter-02-input-variables-types-validations) | Primitives, collections, objects, custom validation rules | 5 | `validate`, `plan` |
| 03 | [Outputs, Locals & Expressions](#chapter-03-outputs-locals-expressions) | Outputs, sensitive masking, locals, ternaries, splats | 4 | `plan` |
| 04 | [Built-in Functions & Collections](#chapter-04-built-in-functions-collections) | String ops, collection math, encoding, files, `try`/`can` | 5 | `plan` |
| 05 | [Meta-Arguments & Scaling](#chapter-05-meta-arguments-resource-scaling) | `count`, `for_each`, `depends_on`, lifecycle rules, drift | 5 | `plan` |
| 06 | [Dynamic Blocks & Advanced HCL](#chapter-06-dynamic-blocks-advanced-hcl) | `dynamic` blocks, custom iterators, nested loops | 4 | `plan` |
| 07 | [Data Sources & State Querying](#chapter-07-data-sources-state-querying) | Local files, archive zips, external scripts, pre/postconditions | 4 | `plan` |
| 08 | [Modular Infrastructure Architecture](#chapter-08-modular-infrastructure-architecture) | Child modules, multi-instance calling, provider aliases | 5 | `validate`, `plan` |
| 09 | [State Management & Refactoring](#chapter-09-state-management-refactoring-surgery) | Declarative `moved` blocks, `import` blocks, replacements | 4 | `plan` |
| 10 | [Native Unit & Integration Testing](#chapter-10-native-unit-integration-testing-tftesthcl) | `.tftest.hcl`, `run` blocks, mocks, `expect_failures` | 4 | `test` |
| 11 | [Production Patterns & Anti-Patterns](#chapter-11-production-patterns-anti-patterns) | Environment mapping, feature flags, tagging factories | 4 | `plan` |
| 12 | [OpenTofu Innovations & Enterprise](#chapter-12-opentofu-innovations-enterprise-features) | State encryption at rest, early evaluation, open registries | 3 | `validate`, `plan` |
| 13 | [Architecture Governance & Standards](#chapter-13-architecture-governance-enterprise-standards) | Root encapsulation, policy scoping (ADR-0005), ephemeral isolation | 3 | `plan` |

---

## Chapter 01: HCL Foundations & Core Primitives

**Directory**: `exercises/01_primitives/`  
**Description**: Master the foundational structure of HashiCorp Configuration Language (HCL). Understand root `terraform` configuration blocks, provider requirements, basic resource declarations, implicit dependency graphs, multi-line heredoc string formatting, and resource lifecycle triggers.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `primitives01` | Terraform Configuration Block | `exercises/01_primitives/primitives01.tf` | `validate` | `terraform` block, `required_version`, `required_providers`, provider sources |
| `primitives02` | First Resource Declaration | `exercises/01_primitives/primitives02.tf` | `plan` | `resource` block syntax, resource types, labels, required attributes |
| `primitives03` | Resource Dependencies | `exercises/01_primitives/primitives03.tf` | `plan` | Implicit dependency graphs, attribute referencing (`<type>.<name>.<attr>`) |
| `primitives04` | String Interpolation & Heredoc | `exercises/01_primitives/primitives04.tf` | `plan` | String templates (`${...}`), indented heredocs (`<<-EOT`), escaping |
| `primitives05` | Syntax & Formatting | `exercises/01_primitives/primitives05.tf` | `validate` | Canonical HCL formatting, alignment, valid quotes, syntax validation |
| `primitives06` | Lifecycle Mechanics | `exercises/01_primitives/primitives06.tf` | `plan` | `terraform_data`, `triggers_replace`, forced resource recreation mechanics |

---

## Chapter 02: Input Variables, Types & Validations

**Directory**: `exercises/02_variables/`  
**Description**: Build robust, parameterizable infrastructure. Learn to define type constraints, manage primitive types, collection types (`list`, `map`, `set`), complex structural types (`object`, `tuple`), optional attributes with defaults, and custom variable validation condition blocks.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `variables01` | Primitive Variable Declarations | `exercises/02_variables/variables01.tf` | `validate` | `variable` blocks, `string`, `number`, `bool`, `description` metadata |
| `variables02` | Collection Types | `exercises/02_variables/variables02.tf` | `plan` | `list(string)`, `map(string)`, `set(string)`, indexing (`[0]`, `["key"]`) |
| `variables03` | Structural Types | `exercises/02_variables/variables03.tf` | `validate` | `object({...})`, `optional(type, default)`, `tuple([...])` constraints |
| `variables04` | Defaults and Nullable | `exercises/02_variables/variables04.tf` | `plan` | Default variable values, `nullable = false`, handling null assignments |
| `variables05` | Custom Variable Validations | `exercises/02_variables/variables05.tf` | `plan` | `validation` blocks, `condition` expressions, `error_message` formatting |

---

## Chapter 03: Outputs, Locals & Expressions

**Directory**: `exercises/03_outputs_locals/`  
**Description**: Implement clean DRY (Don't Repeat Yourself) configurations. Learn to export state via output values, mask confidential secrets with sensitive redaction, compute reusable intermediate values using `locals`, and write ternary conditional and splat projection expressions.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `outputs01` | Defining Outputs & Sensitive Redaction | `exercises/03_outputs_locals/outputs01.tf` | `plan` | `output` blocks, `value`, `description`, `sensitive = true` log masking |
| `locals01` | Locals for Intermediate Calculations | `exercises/03_outputs_locals/locals01.tf` | `plan` | `locals { ... }`, computed values, referencing `local.<name>` |
| `expr01` | Ternary Conditional Expressions | `exercises/03_outputs_locals/expr01.tf` | `plan` | Ternary operators (`cond ? true_val : false_val`), type unification |
| `expr02` | Splat Expressions | `exercises/03_outputs_locals/expr02.tf` | `plan` | Splat operators (`[*]`), extracting attribute lists across collections |

---

## Chapter 04: Built-in Functions & Collections

**Directory**: `exercises/04_functions/`  
**Description**: Leverage HCL's powerful built-in function library. Transform strings, manipulate and flatten nested data structures, encode/decode JSON and YAML documents, interact with the local filesystem and templates, and handle runtime errors safely using `try()` and `can()`.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `func01` | String Manipulation Functions | `exercises/04_functions/func01.tf` | `plan` | `format()`, `join()`, `split()`, `replace()`, `lower()`, `upper()` |
| `func02` | Collection Operations | `exercises/04_functions/func02.tf` | `plan` | `merge()`, `lookup()`, `distinct()`, `slice()`, `flatten()`, `zipmap()` |
| `func03` | Data Encodings | `exercises/04_functions/func03.tf` | `plan` | `jsonencode()`, `jsondecode()`, `yamlencode()`, `base64encode()` |
| `func04` | Filesystem Functions | `exercises/04_functions/func04.tf` | `plan` | `file()`, `templatefile()`, `fileset()`, `fileexists()`, path modules |
| `func05` | Safe Evaluation Expressions | `exercises/04_functions/func05.tf` | `plan` | `try()`, `can()`, dynamic fallback expressions, guarding edge cases |

---

## Chapter 05: Meta-Arguments & Resource Scaling

**Directory**: `exercises/05_meta_arguments/`  
**Description**: Scale and govern resource instantiations declaratively. Explore index-based scaling with `count`, idempotent key-based mapping with `for_each`, explicit graph dependency ordering with `depends_on`, zero-downtime lifecycle management, and drift handling with `ignore_changes`.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `meta01` | Scaling Resources with Count | `exercises/05_meta_arguments/meta01.tf` | `plan` | `count = <n>`, `count.index`, referencing indexed resource instances |
| `meta02` | Idempotent Mapping with For Each | `exercises/05_meta_arguments/meta02.tf` | `plan` | `for_each = toset(...)`, `each.key`, `each.value`, avoiding index churn |
| `meta03` | Explicit Dependency Ordering | `exercises/05_meta_arguments/meta03.tf` | `plan` | `depends_on = [...]`, managing non-attribute execution ordering |
| `meta04` | Resource Lifecycle Blocks | `exercises/05_meta_arguments/meta04.tf` | `plan` | `lifecycle { create_before_destroy, prevent_destroy }` zero-downtime rules |
| `meta05` | Dynamic Drift Handling | `exercises/05_meta_arguments/meta05.tf` | `plan` | `lifecycle { ignore_changes = [...] }`, ignoring out-of-band updates |

---

## Chapter 06: Dynamic Blocks & Advanced HCL

**Directory**: `exercises/06_dynamic_blocks/`  
**Description**: Construct flexible, deeply configurable resources without code duplication. Generate repeated nested configuration blocks dynamically from collections, configure custom iterators, manage multi-level nested dynamic blocks, and conditionally omit blocks based on logic.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `dynamic01` | Basic Dynamic Block Iteration | `exercises/06_dynamic_blocks/dynamic01.tf` | `plan` | `dynamic "block" { for_each = ... content { ... } }`, `.value` access |
| `dynamic02` | Dynamic Blocks with Custom Iterator | `exercises/06_dynamic_blocks/dynamic02.tf` | `plan` | `iterator = custom_name`, referencing `custom_name.key` and `.value` |
| `dynamic03` | Nested Dynamic Blocks | `exercises/06_dynamic_blocks/dynamic03.tf` | `plan` | Hierarchical dynamic blocks, nesting content blocks within parents |
| `dynamic04` | Conditional Dynamic Block Emission | `exercises/06_dynamic_blocks/dynamic04.tf` | `plan` | Conditionally passing empty collections `[]` to emit zero blocks |

---

## Chapter 07: Data Sources & State Querying

**Directory**: `exercises/07_data_sources/`  
**Description**: Integrate existing environment data and external systems into your IaC configurations. Query filesystem state, generate runtime ZIP/tarball archives, execute local scripts via external JSON data sources, and enforce strict execution contracts with custom preconditions and postconditions.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `data01` | Local Filesystem Data Sources | `exercises/07_data_sources/data01.tf` | `plan` | `data "local_file"`, reading files into configuration state |
| `data02` | Archive File Data Sources | `exercises/07_data_sources/data02.tf` | `plan` | `data "archive_file"`, building lambda/deploy packages on the fly |
| `data03` | External Data Source Queries | `exercises/07_data_sources/data03.tf` | `plan` | `data "external"`, executing local scripts and parsing JSON outputs |
| `data04` | Custom Preconditions and Postconditions | `exercises/07_data_sources/data04.tf` | `plan` | `lifecycle { precondition { ... } postcondition { ... } }` |

---

## Chapter 08: Modular Infrastructure Architecture

**Directory**: `exercises/08_modules/`  
**Description**: Architect scalable, maintainable infrastructure through modular composition. Encapsulate reusable child modules with clean inputs and outputs, instantiate modules across environments with `for_each`, pass provider configurations and aliases, and enforce clean submodule boundaries.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `module01` | Building a Clean Child Module | `exercises/08_modules/module01` | `validate` | Module structure (`main.tf`, `variables.tf`, `outputs.tf`), encapsulation |
| `module02` | Calling Local Child Modules | `exercises/08_modules/module02` | `plan` | `module "<name>" { source = "./..." }`, passing arguments, referencing outputs |
| `module03` | Multi-Instance Module Deployment | `exercises/08_modules/module03` | `plan` | `for_each` inside module blocks, addressing multi-instance outputs |
| `module04` | Passing Provider Configurations & Aliases | `exercises/08_modules/module04` | `plan` | `providers = { ... }`, `configuration_aliases` in `required_providers` |
| `module05` | Submodule Boundaries & Clean Architecture | `exercises/08_modules/module05` | `validate` | Decoupled module hierarchy, flat interfaces, antipattern prevention |

---

## Chapter 09: State Management & Refactoring Surgery

**Directory**: `exercises/09_state_refactoring/`  
**Description**: Perform fearless state refactoring without destroying production resources. Master declarative `moved` blocks for zero-downtime resource renames and module migrations, refactor resources from `count` to `for_each`, adopt existing infrastructure via `import` blocks, and configure `replace_triggered_by`.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `state01` | Declarative Refactoring with Moved Blocks | `exercises/09_state_refactoring/state01.tf` | `plan` | `moved { from = ... to = ... }`, non-destructive resource renaming |
| `state02` | Migrating Count to For-Each with Moved Blocks | `exercises/09_state_refactoring/state02.tf` | `plan` | Converting `resource[0]` to `resource["key"]` without recreation |
| `state03` | Declarative Import Blocks | `exercises/09_state_refactoring/state03.tf` | `plan` | `import { to = ... id = "..." }`, onboarding unmanaged cloud resources |
| `state04` | Controlled Resource Replacement | `exercises/09_state_refactoring/state04.tf` | `plan` | `lifecycle { replace_triggered_by = [...] }`, dependency-driven teardowns |

---

## Chapter 10: Native Unit & Integration Testing (.tftest.hcl)

**Directory**: `exercises/10_testing/`  
**Description**: Master native Infrastructure as Code automated testing using `.tftest.hcl`. Author `run` blocks for plan-time and apply-time verification, write custom assertion condition blocks, create mock providers to test configurations without real cloud credentials, and validate failure scenarios with `expect_failures`.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `test01` | Basic Test Assertions with Run Blocks | `exercises/10_testing/test01` | `test` | `.tftest.hcl`, `run { command = plan }`, `assert { condition, error_message }` |
| `test02` | Validating Applied Resources in Tests | `exercises/10_testing/test02` | `test` | `run { command = apply }`, testing provisioned outputs and contract values |
| `test03` | Mocking Providers and Resources | `exercises/10_testing/test03` | `test` | `mock_provider`, `override_resource`, zero-cloud unit testing suites |
| `test04` | Testing Failure Cases with Expect Failures | `exercises/10_testing/test04` | `test` | `expect_failures = [var.<name>]`, validating defensive validation blocks |

---

## Chapter 11: Production Patterns & Anti-Patterns

**Directory**: `exercises/11_patterns/`  
**Description**: Apply battle-tested enterprise design patterns to real-world infrastructure problems. Implement multi-environment configuration maps, safely toggle features with conditional count patterns, enforce organizational tagging standards with merge factories, and build self-service input contracts.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `pattern01` | Multi-Environment Configuration Mapping | `exercises/11_patterns/pattern01.tf` | `plan` | Environment dictionaries in `locals`, dynamic sizing & tiered config |
| `pattern02` | Feature Flags & Conditional Resource Creation | `exercises/11_patterns/pattern02.tf` | `plan` | `count = var.flag ? 1 : 0`, `one()`, `try()`, optional resource patterns |
| `pattern03` | Tagging Factory Pattern | `exercises/11_patterns/pattern03.tf` | `plan` | Merging base tags, environment tags, and custom metadata via `merge()` |
| `pattern04` | Self-Service Input Contracts | `exercises/11_patterns/pattern04.tf` | `plan` | Structured variable contracts, comprehensions with filter predicates |

---

## Chapter 12: OpenTofu Innovations & Enterprise Features

**Directory**: `exercises/12_opentofu/`  
**Description**: Explore the modern capabilities of the open-source OpenTofu engine. Configure native state encryption at rest (using PBKDF2 and AES-GCM), leverage early variable evaluation in provider/backend blocks, and configure provider dependencies with the OpenTofu public registry.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `tofu01` | State Encryption at Rest | `exercises/12_opentofu/tofu01.tf` | `plan` | `terraform { encryption { key_provider, method, state } }` |
| `tofu02` | Early Variable Evaluation | `exercises/12_opentofu/tofu02.tf` | `plan` | Referencing `var.<name>` in early-evaluated provider and backend blocks |
| `tofu03` | OpenTofu Public Registry Integration | `exercises/12_opentofu/tofu03.tf` | `validate` | Open registry provider sources, interoperability, required versions |

---

## Chapter 13: Architecture Governance & Enterprise Standards

**Directory**: `exercises/13_governance/`  
**Description**: Implement enterprise-grade infrastructure governance and architecture standards. Eliminate loose root resources through root module encapsulation, enforce policy encapsulation (ADR-0005: "The module that owns the resource owns the policies that talk to it"), and isolate ephemeral compute workloads.

### Exercises

| Exercise ID | Title | File Path | Mode | Key Concepts Tested |
|---|---|---|:---:|---|
| `gov01` | Root Module Encapsulation | `exercises/13_governance/gov01.tf` | `plan` | Zero loose resources in root environments, modular workload encapsulation |
| `gov02` | Policy Encapsulation (ADR-0005) | `exercises/13_governance/gov02.tf` | `plan` | Resource-owned policy ARNs, eliminating inline IAM wildcard grants |
| `gov03` | Ephemeral Workload Isolation | `exercises/13_governance/gov03.tf` | `plan` | Encapsulating batch/tooling compute to prevent root state pollution |
