# Chapter 13: Architecture Governance & Enterprise Standards

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; Root Module Encapsulation, Policy Scoping, Ephemeral Isolation, and Architectural Governance
-   :material-play-circle: **Interactive Challenges** &bull; 3 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Governance Boundaries

Enterprise Infrastructure as Code demands rigorous boundaries: root module encapsulation, policy-as-code integration, blast radius minimization, and strict separation between transient workspace resources and durable foundation state.

```text
    ┌─────────────────────────────────────────┐
    │     Enterprise Governance Framework     │
    └────────────────────┬────────────────────┘
                         │
     ┌───────────────────┼───────────────────┐
     ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│ Encapsulated │   │ Blast Radius │   │  Lifecycle   │
│ Root Modules │   │ Isolation    │   │  Governance  │
└──────────────┘   └──────────────┘   └──────────────┘
```

---

## 2. Annotated HCL Anatomy & Architecture Standards

```hcl
# Standard root module encapsulation pattern
terraform {
  required_version = ">= 1.6.0"
  
  # Remote backend with strict state locking
  backend "local" {
    path = "terraform.tfstate"
  }
}

# Core boundary: variables strictly validated at edge
variable "governance_tier" {
  type        = string
  description = "Enterprise compliance tier classification."
  
  validation {
    condition     = contains(["bronze", "silver", "gold"], var.governance_tier)
    error_message = "Compliance tier must be bronze, silver, or gold."
  }
}
```

---

## 3. Production Best Practices

1. **Minimize Blast Radius**: Divide infrastructure into small, decoupled workspaces rather than a single monolithic root module.
2. **Standardize Naming & Tagging Contracts**: Enforce organization-wide governance contracts via shared modules and automated policy checks.
3. **Automate Drift Detection**: Run continuous scheduled plans in CI/CD pipelines to detect unmanaged cloud changes.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `governance01` | Root Module Encapsulation | `plan` | Structure root configurations according to enterprise encapsulation rules. |
| `governance02` | Policy-Driven Variable Scoping | `plan` | Enforce organizational compliance rules directly in variable conditions. |
| `governance03` | Ephemeral vs Foundation Isolation | `plan` | Separate durable shared state from short-lived preview infrastructure. |
