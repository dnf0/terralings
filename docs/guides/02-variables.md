# Chapter 02: Input Variables, Types & Validations

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; Primitives, Collections, Structural Objects, Optional Defaults, and Custom Validations
-   :material-play-circle: **Interactive Challenges** &bull; 5 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Type System

Input variables parameterize Terraform/OpenTofu configurations without hardcoding environments. The HCL type system enforces static type checking during the planning phase before any provider APIs are contacted.

```text
    ┌──────────────────────────────┐
    │     terraform.tfvars / CLI   │ (Untrusted External Inputs)
    └──────────────┬───────────────┘
                   │
                   ▼ (Static Type Checking & Coercion)
    ┌──────────────────────────────┐
    │  HCL Variable Type Engine    │ ──► [ Type Mismatch Exception ]
    └──────────────┬───────────────┘
                   │
                   ▼ (Custom Validation Conditions)
    ┌──────────────────────────────┐
    │    Condition Evaluation      │ ──► [ Custom error_message Fail ]
    └──────────────┬───────────────┘
                   │
                   ▼
    ┌──────────────────────────────┐
    │ Validated Module Execution   │
    └──────────────────────────────┘
```

---

## 2. Annotated HCL Anatomy & Schema Reference

```hcl
variable "environment" {
  type        = string
  description = "Target deployment environment tier."
  default     = "dev"

  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "The environment must be one of: 'dev', 'staging', or 'prod'."
  }
}

variable "server_cluster" {
  type = object({
    cluster_name = string
    node_count   = number
    tags         = optional(map(string), {})
    subnets      = list(string)
  })
  description = "Structural configuration schema for cluster orchestration."
}
```

### Type Matrix

- **Primitives**: `string`, `number`, `bool`.
- **Collections**: `list(<TYPE>)`, `map(<TYPE>)`, `set(<TYPE>)`.
- **Structural**: `object({ <ATTR> = <TYPE>, ... })`, `tuple([<TYPE>, ...])`.
- **Modifiers**: `optional(<TYPE>, <DEFAULT_VAL>)`, `nullable = false`.

---

## 3. Production Best Practices

1. **Always Specify Type & Description**: Avoid untyped `variable "foo" {}`. Always declare strict type signatures and descriptive documentation.
2. **Fail Fast with Custom Validations**: Enforce boundary conditions using `validation` blocks with clear, actionable `error_message` strings.
3. **Use Optional Attributes with Defaults**: Modern HCL (`optional(type, default)`) allows nested attributes in complex objects to remain concise without requiring repetitive boilerplate.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `variables01` | Primitive Variable Declarations | `validate` | Define `string`, `number`, and `bool` variables with type constraints. |
| `variables02` | Collection Types | `plan` | Work with `list(string)`, `map(string)`, and `set(string)` indexing. |
| `variables03` | Structural Types | `validate` | Declare structured schemas with `object({...})` and `optional(...)`. |
| `variables04` | Defaults and Nullable | `plan` | Master default fallback behaviors and `nullable = false` enforcement. |
| `variables05` | Custom Variable Validations | `plan` | Implement custom validation condition blocks and custom error messages. |
