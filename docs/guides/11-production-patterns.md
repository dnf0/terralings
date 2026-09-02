# Chapter 11: Production Patterns & Anti-Patterns

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; Environment Maps, Feature Flags, Tagging Factories, and Self-Service Contracts
-   :material-api: **Primary Patterns** &bull; Lookup Matrices, Optional Resource Toggles, `merge()` Factories, Dynamic Contracts
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html?chapter=11){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Enterprise Design Patterns

In production Infrastructure as Code, codebases scale to support hundreds of engineers and diverse environments. **Enterprise Design Patterns** provide robust architectural blueprints to maintain clean separation of concerns, eliminate copy-pasted configuration, and enforce organizational governance.

```text
┌─────────────────────────────────────────────────────────────┐
│                 Enterprise Pattern Architecture             │
│                                                             │
│   1. Environment Matrix Pattern                             │
│   ┌───────────────────────────────────────────────────────┐ │
│   │ env_config = {                                        │ │
│   │   dev  = { tier = "t4g.nano",  replicas = 1 }         │ │
│   │   prod = { tier = "r6g.xlarge", replicas = 5 }        │ │
│   │ }                                                     │ │
│   └──────────────────────────┬────────────────────────────┘ │
│                              │                              │
│                              ▼                              │
│   2. Tagging Factory Pattern (Hierarchical Merging)         │
│   [ Global Org Tags ] ──► [ Env Tags ] ──► [ Resource Tags] │
│                              │                              │
│                              ▼                              │
│   3. Self-Service Input Contract & Filtering                │
│   [ Developer Spec ] ──► [ Schema Sanitizer ] ──► [ Engine] │
└─────────────────────────────────────────────────────────────┘
```

Core Enterprise Patterns:
1. **Environment Configuration Matrix**: Centralizes per-environment sizing, feature flags, and topology in structured maps inside `locals`.
2. **Deterministic Tagging Factory**: Uses layered `merge()` calls to guarantee required compliance tags are always attached.
3. **Safe Feature Toggles**: Implements clean 0/1 conditional patterns with safe scalar extraction via `one()`.
4. **Self-Service Input Contracts**: Validates, normalizes, and filters complex developer specifications into safe infrastructure blueprints.

---

## 2. Annotated Production HCL Anatomy & Field Reference

Below is a production-grade configuration demonstrating environment mapping, tagging factories, and feature toggles:

```hcl
variable "environment" {
  type        = string
  description = "Target deployment environment tier."
  default     = "production"

  validation {
    condition     = contains(["dev", "staging", "production"], var.environment)
    error_message = "Allowed environments: dev, staging, production."
  }
}

variable "enable_waf" {
  type        = bool
  description = "Toggle Web Application Firewall integration."
  default     = true
}

variable "custom_tags" {
  type    = map(string)
  default = { Service = "Checkout" }
}

locals {
  # 1. Environment Configuration Matrix
  env_matrix = {
    dev = {
      instance_count = 1
      tier           = "basic"
      backup_retention = 7
    }
    staging = {
      instance_count = 2
      tier           = "standard"
      backup_retention = 14
    }
    production = {
      instance_count = 5
      tier           = "enterprise"
      backup_retention = 90
    }
  }

  active_config = local.env_matrix[var.environment]

  # 2. Layered Tagging Factory
  base_tags = {
    ManagedBy = "Terralings"
    Repo      = "infra-core"
  }

  environment_tags = {
    Environment = var.environment
    Tier        = local.active_config.tier
  }

  effective_tags = merge(
    local.base_tags,
    local.environment_tags,
    var.custom_tags
  )
}

# 3. Primary Scaled Workload
resource "terraform_data" "service_nodes" {
  count = local.active_config.instance_count

  input = {
    hostname = "app-${var.environment}-${count.index + 1}"
    tier     = local.active_config.tier
    tags     = local.effective_tags
  }
}

# 4. Feature-Flagged Optional Resource
resource "terraform_data" "waf_shield" {
  count = var.enable_waf ? 1 : 0

  input = {
    protection_tier = "advanced"
    target_nodes    = terraform_data.service_nodes[*].id
  }
}
```

---

## 3. Real-World Architectural Patterns

### Pattern 1: Safe Optional Extraction with `one()`

```hcl
# Safely extract optional single resource attribute without index error
output "waf_protection_id" {
  description = "WAF ID if enabled, or null if disabled."
  value       = one(terraform_data.waf_shield[*].id)
}
```

### Pattern 2: Self-Service Contract with Filter Predicates

```hcl
variable "developer_service_specs" {
  type = list(object({
    name     = string
    enabled  = bool
    protocol = string
  }))
  default = [
    { name = "auth",   enabled = true,  protocol = "grpc" },
    { name = "legacy", enabled = false, protocol = "http" }
  ]
}

# Filter only active services into deployment map
locals {
  active_services = {
    for s in var.developer_service_specs : s.name => s
    if s.enabled && s.protocol == "grpc"
  }
}
```

---

## 4. Production Hardening & Operational Governance

- **Avoid the "Blast Radius" Antipattern**: Do not place an entire corporate cloud estate into a single root module. Segment state files by lifecycle and domain (e.g., networking vs compute vs database).
- **Enforce Mandatory Tags**: Use the tagging factory pattern to guarantee `CostCenter`, `Owner`, and `Environment` tags are never omitted.
- **Never Use Workspaces for Long-Lived Environments**: Terraform/OpenTofu CLI workspaces share the same backend configuration. Use distinct directory structures (`envs/dev/`, `envs/prod/`) for true isolation.

---

## 5. Failure Modes & Diagnostic Triage Tree

??? failure "Error: Invalid index on optional resource (`resource[0].id` when count = 0)"
    **Root Cause:** Directly indexing `[0]` on a resource that was conditionally disabled with `count = 0`.

    **Diagnostic Triage Sequence:**
    1. Replace `resource[0].attr` with `one(resource[*].attr)`.
    2. `one()` safely returns `null` if the list is empty, and returns the single element if `count = 1`.

??? failure "Error: Key not found in environment matrix map"
    **Root Cause:** `var.environment` contains a value that does not match any key in the `env_matrix` dictionary.

    **Diagnostic Triage Sequence:**
    1. Add a `validation` block to `var.environment` restricting allowed values with `contains(["dev", ...], var.environment)`.
    2. Use `lookup(local.env_matrix, var.environment, local.env_matrix["dev"])` as a safe fallback.

---

## 6. Interactive Practice Matrix

Practice concepts from this chapter directly in the interactive WebAssembly sandbox:

| Exercise ID | Challenge Description | Direct Link | Action |
| :--- | :--- | :--- | :--- |
| **`pattern01`** | Multi-Environment Configuration Mapping | [`../playground/index.html?exercise=pattern01`](../playground/index.html?exercise=pattern01) | [**⚡ Solve in Playground →**](../playground/index.html?exercise=pattern01){ .md-button .md-button--primary } |
| **`pattern02`** | Feature Flags & Conditional Resource Creation | [`../playground/index.html?exercise=pattern02`](../playground/index.html?exercise=pattern02) | [**⚡ Solve in Playground →**](../playground/index.html?exercise=pattern02){ .md-button .md-button--primary } |
| **`pattern03`** | Tagging Factory Pattern | [`../playground/index.html?exercise=pattern03`](../playground/index.html?exercise=pattern03) | [**⚡ Solve in Playground →**](../playground/index.html?exercise=pattern03){ .md-button .md-button--primary } |
| **`pattern04`** | Self-Service Input Contracts | [`../playground/index.html?exercise=pattern04`](../playground/index.html?exercise=pattern04) | [**⚡ Solve in Playground →**](../playground/index.html?exercise=pattern04){ .md-button .md-button--primary } |
