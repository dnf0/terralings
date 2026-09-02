# Chapter 11: Production Patterns & Anti-Patterns

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; Environment Mapping, Feature Flags, Tagging Factories, and Production IaC Hardening
-   :material-play-circle: **Interactive Challenges** &bull; 4 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Hardening Patterns

Production Infrastructure as Code requires strategies for managing environment drift, feature toggles, centralized tagging policies, and eliminating common anti-patterns like hardcoded credentials or brittle conditional logic.

```text
    ┌──────────────────────────────┐
    │     Environment Strategy     │
    │  (lookup table / map matrix) │
    └──────────────┬───────────────┘
                   │
                   ▼
    ┌──────────────────────────────┐
    │  Standardized Tag Factory    │ ──► [ Compliance & Cost Attribution ]
    │  (CostCenter, Owner, Tier)   │
    └──────────────┬───────────────┘
                   │
                   ▼
    ┌──────────────────────────────┐
    │ Feature Toggles & Guards     │ ──► [ Safe Phased Rollouts ]
    └──────────────────────────────┘
```

---

## 2. Annotated HCL Anatomy & Common Patterns

```hcl
locals {
  # Environment configuration matrix pattern
  env_specs = {
    dev = {
      instance_count = 1
      backup_enabled = false
    }
    prod = {
      instance_count = 5
      backup_enabled = true
    }
  }

  current_spec = local.env_specs[var.environment]

  # Standardized enterprise tag factory
  mandatory_tags = {
    ManagedBy   = "Terraform"
    Environment = var.environment
    Repository  = "infra-core"
  }
}
```

---

## 3. Production Best Practices

1. **Use Lookup Tables over Deep Ternaries**: Model multi-environment configuration variations as structured maps rather than deeply nested `a ? b : (c ? d : e)` expressions.
2. **Centralize Tagging**: Enforce tag schemas at the root or provider level using `default_tags` or merged local tag factories.
3. **Avoid Unbounded Cardinality**: Guard collection expansions to avoid accidental runaway resource creation.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `patterns01` | Environment Lookup Matrix | `plan` | Implement clean environment-specific configuration maps. |
| `patterns02` | Boolean Feature Flags | `plan` | Control resource creation and optional features with flags. |
| `patterns03` | Centralized Tag Factory | `plan` | Merge organizational and resource-specific tags cleanly. |
| `patterns04` | Anti-Pattern Remediation | `plan` | Identify and fix hardcoded credentials and brittle conditional logic. |
