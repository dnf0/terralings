# Curriculum Syllabus

Terralings features a comprehensive 13-chapter curriculum comprising **56 progressive exercises**. The curriculum covers the entire spectrum of modern Infrastructure as Code (IaC)—from foundational syntax and variables to advanced state surgery, native `.tftest.hcl` testing, OpenTofu state encryption, and architectural governance.

---

## Curriculum Overview Matrix

| # | Chapter Name | Focus Area | Exercises | Reference Guide | Verification |
|---|---|---|:---:|:---:|:---:|
| 01 | [HCL Foundations & Core Primitives](#chapter-01-hcl-foundations-core-primitives) | Blocks, providers, resources, heredocs, lifecycle | 6 | [Ch 01 Guide →](guides/01-primitives.md) | `validate`, `plan` |
| 02 | [Input Variables, Types & Validations](#chapter-02-input-variables-types-validations) | Primitives, collections, objects, custom validation rules | 5 | [Ch 02 Guide →](guides/02-variables.md) | `validate`, `plan` |
| 03 | [Outputs, Locals & Expressions](#chapter-03-outputs-locals-expressions) | Outputs, sensitive masking, locals, ternaries, splats | 4 | [Ch 03 Guide →](guides/03-outputs-locals.md) | `plan` |
| 04 | [Built-in Functions & Collections](#chapter-04-built-in-functions-collections) | String ops, collection math, encoding, files, `try`/`can` | 5 | [Ch 04 Guide →](guides/04-functions.md) | `plan` |
| 05 | [Meta-Arguments & Scaling](#chapter-05-meta-arguments-resource-scaling) | `count`, `for_each`, `depends_on`, lifecycle rules, drift | 5 | [Ch 05 Guide →](guides/05-meta-arguments.md) | `plan` |
| 06 | [Dynamic Blocks & Advanced HCL](#chapter-06-dynamic-blocks-advanced-hcl) | `dynamic` blocks, custom iterators, nested loops | 4 | [Ch 06 Guide →](guides/06-dynamic-blocks.md) | `plan` |
| 07 | [Data Sources & State Querying](#chapter-07-data-sources-state-querying) | Local files, archive zips, external scripts, pre/postconditions | 4 | [Ch 07 Guide →](guides/07-data-sources.md) | `plan` |
| 08 | [Modular Infrastructure Architecture](#chapter-08-modular-infrastructure-architecture) | Child modules, multi-instance calling, provider aliases | 5 | [Ch 08 Guide →](guides/08-modules.md) | `validate`, `plan` |
| 09 | [State Management & Refactoring](#chapter-09-state-management-refactoring-surgery) | Declarative `moved` blocks, `import` blocks, replacements | 4 | [Ch 09 Guide →](guides/09-state-management.md) | `plan` |
| 10 | [Native Unit & Integration Testing](#chapter-10-native-unit-integration-testing-tftesthcl) | `.tftest.hcl`, `run` blocks, mocks, `expect_failures` | 4 | [Ch 10 Guide →](guides/10-testing.md) | `test` |
| 11 | [Production Patterns & Anti-Patterns](#chapter-11-production-patterns-anti-patterns) | Environment mapping, feature flags, tagging factories | 4 | [Ch 11 Guide →](guides/11-production-patterns.md) | `plan` |
| 12 | [OpenTofu Innovations & Enterprise](#chapter-12-opentofu-innovations-enterprise-features) | State encryption at rest, early evaluation, open registries | 3 | [Ch 12 Guide →](guides/12-opentofu.md) | `validate`, `plan` |
| 13 | [Architecture Governance & Standards](#chapter-13-architecture-governance-enterprise-standards) | Root encapsulation, policy scoping (ADR-0005), ephemeral isolation | 3 | [Ch 13 Guide →](guides/13-governance.md) | `plan` |

---

## [Chapter 01: HCL Foundations & Core Primitives](guides/01-primitives.md)

**Reference Guide**: [📖 Chapter 01 Reference Guide & Architecture Specs](guides/01-primitives.md)  
**Directory**: `exercises/01_primitives/`  
**Description**: Master the foundational structure of HashiCorp Configuration Language (HCL). Understand root `terraform` configuration blocks, provider requirements, basic resource declarations, implicit dependency graphs, multi-line heredoc string formatting, and resource lifecycle triggers.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`primitives01`**](playground/index.html?exercise=primitives01) | Terraform Configuration Block | `exercises/01_primitives/primitives01.tf` | `validate` | [⚡ Solve in Playground →](playground/index.html?exercise=primitives01) |
| [**`primitives02`**](playground/index.html?exercise=primitives02) | First Resource Declaration | `exercises/01_primitives/primitives02.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=primitives02) |
| [**`primitives03`**](playground/index.html?exercise=primitives03) | Resource Dependencies | `exercises/01_primitives/primitives03.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=primitives03) |
| [**`primitives04`**](playground/index.html?exercise=primitives04) | String Interpolation & Heredoc | `exercises/01_primitives/primitives04.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=primitives04) |
| [**`primitives05`**](playground/index.html?exercise=primitives05) | Syntax & Formatting | `exercises/01_primitives/primitives05.tf` | `validate` | [⚡ Solve in Playground →](playground/index.html?exercise=primitives05) |
| [**`primitives06`**](playground/index.html?exercise=primitives06) | Lifecycle Mechanics | `exercises/01_primitives/primitives06.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=primitives06) |

---

## [Chapter 02: Input Variables, Types & Validations](guides/02-variables.md)

**Reference Guide**: [📖 Chapter 02 Reference Guide & Architecture Specs](guides/02-variables.md)  
**Directory**: `exercises/02_variables/`  
**Description**: Build robust, parameterizable infrastructure. Learn to define type constraints, manage primitive types, collection types (`list`, `map`, `set`), complex structural types (`object`, `tuple`), optional attributes with defaults, and custom variable validation condition blocks.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`variables01`**](playground/index.html?exercise=variables01) | Primitive Variable Declarations | `exercises/02_variables/variables01.tf` | `validate` | [⚡ Solve in Playground →](playground/index.html?exercise=variables01) |
| [**`variables02`**](playground/index.html?exercise=variables02) | Collection Types | `exercises/02_variables/variables02.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=variables02) |
| [**`variables03`**](playground/index.html?exercise=variables03) | Structural Types | `exercises/02_variables/variables03.tf` | `validate` | [⚡ Solve in Playground →](playground/index.html?exercise=variables03) |
| [**`variables04`**](playground/index.html?exercise=variables04) | Defaults and Nullable | `exercises/02_variables/variables04.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=variables04) |
| [**`variables05`**](playground/index.html?exercise=variables05) | Custom Variable Validations | `exercises/02_variables/variables05.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=variables05) |

---

## [Chapter 03: Outputs, Locals & Expressions](guides/03-outputs-locals.md)

**Reference Guide**: [📖 Chapter 03 Reference Guide & Architecture Specs](guides/03-outputs-locals.md)  
**Directory**: `exercises/03_outputs_locals/`  
**Description**: Implement clean DRY (Don't Repeat Yourself) configurations. Learn to export state via output values, mask confidential secrets with sensitive redaction, compute reusable intermediate values using `locals`, and write ternary conditional and splat projection expressions.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`outputs01`**](playground/index.html?exercise=outputs01) | Defining Outputs & Sensitive Redaction | `exercises/03_outputs_locals/outputs01.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=outputs01) |
| [**`locals01`**](playground/index.html?exercise=locals01) | Locals for Intermediate Calculations | `exercises/03_outputs_locals/locals01.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=locals01) |
| [**`expr01`**](playground/index.html?exercise=expr01) | Ternary Conditional Expressions | `exercises/03_outputs_locals/expr01.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=expr01) |
| [**`expr02`**](playground/index.html?exercise=expr02) | Splat Expressions | `exercises/03_outputs_locals/expr02.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=expr02) |

---

## [Chapter 04: Built-in Functions & Collections](guides/04-functions.md)

**Reference Guide**: [📖 Chapter 04 Reference Guide & Architecture Specs](guides/04-functions.md)  
**Directory**: `exercises/04_functions/`  
**Description**: Leverage HCL's powerful built-in function library. Transform strings, manipulate and flatten nested data structures, encode/decode JSON and YAML documents, interact with the local filesystem and templates, and handle runtime errors safely using `try()` and `can()`.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`func01`**](playground/index.html?exercise=func01) | String Manipulation Functions | `exercises/04_functions/func01.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=func01) |
| [**`func02`**](playground/index.html?exercise=func02) | Collection Operations | `exercises/04_functions/func02.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=func02) |
| [**`func03`**](playground/index.html?exercise=func03) | Data Encodings | `exercises/04_functions/func03.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=func03) |
| [**`func04`**](playground/index.html?exercise=func04) | Filesystem Functions | `exercises/04_functions/func04.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=func04) |
| [**`func05`**](playground/index.html?exercise=func05) | Safe Evaluation Expressions | `exercises/04_functions/func05.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=func05) |

---

## [Chapter 05: Meta-Arguments & Resource Scaling](guides/05-meta-arguments.md)

**Reference Guide**: [📖 Chapter 05 Reference Guide & Architecture Specs](guides/05-meta-arguments.md)  
**Directory**: `exercises/05_meta_arguments/`  
**Description**: Scale and govern resource instantiations declaratively. Explore index-based scaling with `count`, idempotent key-based mapping with `for_each`, explicit graph dependency ordering with `depends_on`, zero-downtime lifecycle management, and drift handling with `ignore_changes`.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`meta01`**](playground/index.html?exercise=meta01) | Scaling Resources with Count | `exercises/05_meta_arguments/meta01.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=meta01) |
| [**`meta02`**](playground/index.html?exercise=meta02) | Idempotent Mapping with For Each | `exercises/05_meta_arguments/meta02.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=meta02) |
| [**`meta03`**](playground/index.html?exercise=meta03) | Explicit Dependency Ordering | `exercises/05_meta_arguments/meta03.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=meta03) |
| [**`meta04`**](playground/index.html?exercise=meta04) | Resource Lifecycle Blocks | `exercises/05_meta_arguments/meta04.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=meta04) |
| [**`meta05`**](playground/index.html?exercise=meta05) | Dynamic Drift Handling | `exercises/05_meta_arguments/meta05.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=meta05) |

---

## [Chapter 06: Dynamic Blocks & Advanced HCL](guides/06-dynamic-blocks.md)

**Reference Guide**: [📖 Chapter 06 Reference Guide & Architecture Specs](guides/06-dynamic-blocks.md)  
**Directory**: `exercises/06_dynamic_blocks/`  
**Description**: Construct flexible, deeply configurable resources without code duplication. Generate repeated nested configuration blocks dynamically from collections, configure custom iterators, manage multi-level nested dynamic blocks, and conditionally omit blocks based on logic.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`dynamic01`**](playground/index.html?exercise=dynamic01) | Basic Dynamic Block Iteration | `exercises/06_dynamic_blocks/dynamic01.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=dynamic01) |
| [**`dynamic02`**](playground/index.html?exercise=dynamic02) | Dynamic Blocks with Custom Iterator | `exercises/06_dynamic_blocks/dynamic02.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=dynamic02) |
| [**`dynamic03`**](playground/index.html?exercise=dynamic03) | Nested Dynamic Blocks | `exercises/06_dynamic_blocks/dynamic03.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=dynamic03) |
| [**`dynamic04`**](playground/index.html?exercise=dynamic04) | Conditional Dynamic Block Emission | `exercises/06_dynamic_blocks/dynamic04.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=dynamic04) |

---

## [Chapter 07: Data Sources & State Querying](guides/07-data-sources.md)

**Reference Guide**: [📖 Chapter 07 Reference Guide & Architecture Specs](guides/07-data-sources.md)  
**Directory**: `exercises/07_data_sources/`  
**Description**: Integrate existing environment data and external systems into your IaC configurations. Query filesystem state, generate runtime ZIP/tarball archives, execute local scripts via external JSON data sources, and enforce strict execution contracts with custom preconditions and postconditions.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`data01`**](playground/index.html?exercise=data01) | Local Filesystem Data Sources | `exercises/07_data_sources/data01.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=data01) |
| [**`data02`**](playground/index.html?exercise=data02) | Archive File Data Sources | `exercises/07_data_sources/data02.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=data02) |
| [**`data03`**](playground/index.html?exercise=data03) | External Data Source Queries | `exercises/07_data_sources/data03.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=data03) |
| [**`data04`**](playground/index.html?exercise=data04) | Custom Preconditions and Postconditions | `exercises/07_data_sources/data04.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=data04) |

---

## [Chapter 08: Modular Infrastructure Architecture](guides/08-modules.md)

**Reference Guide**: [📖 Chapter 08 Reference Guide & Architecture Specs](guides/08-modules.md)  
**Directory**: `exercises/08_modules/`  
**Description**: Architect scalable, maintainable infrastructure through modular composition. Encapsulate reusable child modules with clean inputs and outputs, instantiate modules across environments with `for_each`, pass provider configurations and aliases, and enforce clean submodule boundaries.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`module01`**](playground/index.html?exercise=module01) | Building a Clean Child Module | `exercises/08_modules/module01` | `validate` | [⚡ Solve in Playground →](playground/index.html?exercise=module01) |
| [**`module02`**](playground/index.html?exercise=module02) | Calling Local Child Modules | `exercises/08_modules/module02` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=module02) |
| [**`module03`**](playground/index.html?exercise=module03) | Multi-Instance Module Deployment | `exercises/08_modules/module03` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=module03) |
| [**`module04`**](playground/index.html?exercise=module04) | Passing Provider Configurations & Aliases | `exercises/08_modules/module04` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=module04) |
| [**`module05`**](playground/index.html?exercise=module05) | Submodule Boundaries & Clean Architecture | `exercises/08_modules/module05` | `validate` | [⚡ Solve in Playground →](playground/index.html?exercise=module05) |

---

## [Chapter 09: State Management & Refactoring Surgery](guides/09-state-management.md)

**Reference Guide**: [📖 Chapter 09 Reference Guide & Architecture Specs](guides/09-state-management.md)  
**Directory**: `exercises/09_state_refactoring/`  
**Description**: Perform fearless state refactoring without destroying production resources. Master declarative `moved` blocks for zero-downtime resource renames and module migrations, refactor resources from `count` to `for_each`, adopt existing infrastructure via `import` blocks, and configure `replace_triggered_by`.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`state01`**](playground/index.html?exercise=state01) | Declarative Refactoring with Moved Blocks | `exercises/09_state_refactoring/state01.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=state01) |
| [**`state02`**](playground/index.html?exercise=state02) | Migrating Count to For-Each with Moved Blocks | `exercises/09_state_refactoring/state02.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=state02) |
| [**`state03`**](playground/index.html?exercise=state03) | Declarative Import Blocks | `exercises/09_state_refactoring/state03.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=state03) |
| [**`state04`**](playground/index.html?exercise=state04) | Controlled Resource Replacement | `exercises/09_state_refactoring/state04.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=state04) |

---

## [Chapter 10: Native Unit & Integration Testing (.tftest.hcl)](guides/10-testing.md)

**Reference Guide**: [📖 Chapter 10 Reference Guide & Architecture Specs](guides/10-testing.md)  
**Directory**: `exercises/10_testing/`  
**Description**: Master native Infrastructure as Code automated testing using `.tftest.hcl`. Author `run` blocks for plan-time and apply-time verification, write custom assertion condition blocks, create mock providers to test configurations without real cloud credentials, and validate failure scenarios with `expect_failures`.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`test01`**](playground/index.html?exercise=test01) | Basic Test Assertions with Run Blocks | `exercises/10_testing/test01` | `test` | [⚡ Solve in Playground →](playground/index.html?exercise=test01) |
| [**`test02`**](playground/index.html?exercise=test02) | Validating Applied Resources in Tests | `exercises/10_testing/test02` | `test` | [⚡ Solve in Playground →](playground/index.html?exercise=test02) |
| [**`test03`**](playground/index.html?exercise=test03) | Mocking Providers and Resources | `exercises/10_testing/test03` | `test` | [⚡ Solve in Playground →](playground/index.html?exercise=test03) |
| [**`test04`**](playground/index.html?exercise=test04) | Testing Failure Cases with Expect Failures | `exercises/10_testing/test04` | `test` | [⚡ Solve in Playground →](playground/index.html?exercise=test04) |

---

## [Chapter 11: Production Patterns & Anti-Patterns](guides/11-production-patterns.md)

**Reference Guide**: [📖 Chapter 11 Reference Guide & Architecture Specs](guides/11-production-patterns.md)  
**Directory**: `exercises/11_patterns/`  
**Description**: Apply battle-tested enterprise design patterns to real-world infrastructure problems. Implement multi-environment configuration maps, safely toggle features with conditional count patterns, enforce organizational tagging standards with merge factories, and build self-service input contracts.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`pattern01`**](playground/index.html?exercise=pattern01) | Multi-Environment Configuration Mapping | `exercises/11_patterns/pattern01.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=pattern01) |
| [**`pattern02`**](playground/index.html?exercise=pattern02) | Feature Flags & Conditional Resource Creation | `exercises/11_patterns/pattern02.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=pattern02) |
| [**`pattern03`**](playground/index.html?exercise=pattern03) | Tagging Factory Pattern | `exercises/11_patterns/pattern03.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=pattern03) |
| [**`pattern04`**](playground/index.html?exercise=pattern04) | Self-Service Input Contracts | `exercises/11_patterns/pattern04.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=pattern04) |

---

## [Chapter 12: OpenTofu Innovations & Enterprise Features](guides/12-opentofu.md)

**Reference Guide**: [📖 Chapter 12 Reference Guide & Architecture Specs](guides/12-opentofu.md)  
**Directory**: `exercises/12_opentofu/`  
**Description**: Explore the modern capabilities of the open-source OpenTofu engine. Configure native state encryption at rest (using PBKDF2 and AES-GCM), leverage early variable evaluation in provider/backend blocks, and configure provider dependencies with the OpenTofu public registry.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`tofu01`**](playground/index.html?exercise=tofu01) | State Encryption at Rest | `exercises/12_opentofu/tofu01.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=tofu01) |
| [**`tofu02`**](playground/index.html?exercise=tofu02) | Early Variable Evaluation | `exercises/12_opentofu/tofu02.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=tofu02) |
| [**`tofu03`**](playground/index.html?exercise=tofu03) | OpenTofu Public Registry Integration | `exercises/12_opentofu/tofu03.tf` | `validate` | [⚡ Solve in Playground →](playground/index.html?exercise=tofu03) |

---

## [Chapter 13: Architecture Governance & Enterprise Standards](guides/13-governance.md)

**Reference Guide**: [📖 Chapter 13 Reference Guide & Architecture Specs](guides/13-governance.md)  
**Directory**: `exercises/13_governance/`  
**Description**: Implement enterprise-grade infrastructure governance and architecture standards. Eliminate loose root resources through root module encapsulation, enforce policy encapsulation (ADR-0005: "The module that owns the resource owns the policies that talk to it"), and isolate ephemeral compute workloads.

### Exercises

| Exercise ID | Title | File Path | Mode | Action |
|---|---|---|:---:|:---:|
| [**`gov01`**](playground/index.html?exercise=gov01) | Root Module Encapsulation | `exercises/13_governance/gov01.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=gov01) |
| [**`gov02`**](playground/index.html?exercise=gov02) | Policy Encapsulation (ADR-0005) | `exercises/13_governance/gov02.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=gov02) |
| [**`gov03`**](playground/index.html?exercise=gov03) | Ephemeral Workload Isolation | `exercises/13_governance/gov03.tf` | `plan` | [⚡ Solve in Playground →](playground/index.html?exercise=gov03) |
