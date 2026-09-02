# Chapter 08: Modular Infrastructure Architecture

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; Child Modules, Input Encapsulation, Multi-Instance Calls, and Provider Aliasing
-   :material-play-circle: **Interactive Challenges** &bull; 5 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Module Encapsulation

Modules are the primary mechanism for packaging, reusing, and abstracting infrastructure components. A child module encapsulates implementation details behind a clean boundary of input variables and output values.

```text
    ┌─────────────────────────────────────────┐
    │              Root Module                │
    │  module "network" { ... }               │
    │  module "compute" { ... }               │
    └────────────────────┬────────────────────┘
                         │
           ┌─────────────┴─────────────┐
           ▼                           ▼
    ┌──────────────┐            ┌──────────────┐
    │ Child Module │            │ Child Module │
    │   (Network)  │            │  (Compute)   │
    └──────────────┘            └──────────────┘
```

---

## 2. Annotated HCL Anatomy & Schema Reference

```hcl
module "storage_vault" {
  source = "./modules/vault"

  # Module inputs
  vault_name  = "production-vault"
  replication = true

  # Provider inheritance and aliasing
  providers = {
    local = local
  }
}

# Consuming module outputs
output "vault_id" {
  value = module.storage_vault.vault_identifier
}
```

---

## 3. Production Best Practices

1. **Explicit Version Pinning**: For remote modules (Git / Terraform Registry), always pin git tags or semantic versions (`?ref=v1.2.0`).
2. **Narrow Input Surface**: Keep module input variables cohesive; do not expose internal implementation knobs that break abstraction.
3. **Avoid Deep Nesting**: Limit module nesting depth to 1–2 levels. Overly nested hierarchies make refactoring and state surgery painful.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `modules01` | Calling Local Child Modules | `validate` | Instantiate a child module and pass typed input variables. |
| `modules02` | Consuming Module Outputs | `plan` | Propagate child module outputs to root resources and outputs. |
| `modules03` | Multi-Instance Module Calls | `plan` | Scale whole modules across environments using `for_each`. |
| `modules04` | Module Provider Aliases | `validate` | Pass explicit provider instances and aliases into child modules. |
| `modules05` | Reusable Module Composition | `plan` | Compose multiple child modules into a unified multi-tier architecture. |
