# Chapter 01: HCL Foundations & Core Primitives

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; Syntax, Blocks, Providers, Implicit Dependencies, and Lifecycle Mechanics
-   :material-api: **Primary Primitives** &bull; `terraform`, `provider`, `resource`, `terraform_data`
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html?chapter=1){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Engine Mechanics

In Terraform and OpenTofu, **HashiCorp Configuration Language (HCL)** is a declarative domain-specific language engineered to express desired infrastructure state. The core engine constructs a Directed Acyclic Graph (DAG) of resources, computing execution plans based on explicit and implicit dependencies.

```text
┌─────────────────────────────────────────────────────────────┐
│                   HCL Execution Pipeline                    │
│                                                             │
│  ┌─────────────────┐             ┌───────────────────────┐  │
│  │  Root Config    │             │  Provider RPC Plugins │  │
│  │  (*.tf Files)   │ ──(Parse)──►│  (local, aws, google) │  │
│  └────────┬────────┘             └───────────┬───────────┘  │
│           │                                  │              │
│           ▼                                  ▼              │
│     [ AST Builder ]                     [ Schema Sync]      │
│     [ HCL Spec     ]                     [ Validate    ]    │
│           │                                  ▲              │
│           ▼                                  │              │
│  ┌───────────────────────────────────────────┴───────────┐  │
│  │           Directed Acyclic Graph (DAG) Engine         │  │
│  │           - Node Walk & Topological Ordering          │  │
│  │           - Cycle Detection & Implicit Referencing    │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

When configuration files are evaluated:
1. The lexer and parser read all `.tf` files in the directory concurrently into a single memory scope.
2. The engine analyzes attribute references (e.g. `local_file.config.id`) to construct dependency vertices in the DAG.
3. Provider plugins are initialized via gRPC handshake to fetch resource schemas.
4. The execution planner walks the DAG in topological order, running independent vertices concurrently while sequencing dependent nodes.

---

## 2. Annotated Production HCL Anatomy & Field Reference

Below is a production-grade declarative manifest demonstrating root configuration, provider pinning, implicit dependencies, and lifecycle hooks:

```hcl
terraform {
  required_version = ">= 1.6.0"

  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.4.0"
    }
  }
}

provider "local" {}

# Primary managed resource with formatted heredoc content
resource "local_file" "app_config" {
  filename        = "${path.module}/app.json"
  content         = <<-EOT
    {
      "environment": "production",
      "log_level": "info",
      "max_connections": 100
    }
  EOT
  file_permission = "0644"

  lifecycle {
    create_before_destroy = true
  }
}

# Downstream state entity implicitly depending on local_file.app_config
resource "terraform_data" "config_audit" {
  input = {
    file_id      = local_file.app_config.id
    sha256_hash  = local_file.app_config.content_sha256
    generated_at = timestamp()
  }

  lifecycle {
    triggers_replace = [
      local_file.app_config.content
    ]
  }
}
```

### Key Field Schema Reference

| Field / Block | Type | Description |
| :--- | :--- | :--- |
| `terraform.required_version` | `String` | Semantic version constraint for the Terraform/OpenTofu CLI binary (e.g., `>= 1.6.0, < 2.0.0`). |
| `terraform.required_providers` | `Map(Object)` | Declares upstream provider registry addresses (`source`) and pessimistic version bounds (`version`). |
| `resource.<type>.<name>` | `Block` | Declares an infrastructure component managed by a provider plugin. |
| `lifecycle.create_before_destroy` | `Boolean` | Inverts the default destruction order so new resources are provisioned before existing ones are torn down. |
| `lifecycle.triggers_replace` | `List(Any)` | Forces resource recreation whenever any evaluated value in the list changes. |
| `terraform_data` | `Built-in Resource` | Native state tracking resource replacing deprecated `null_resource` without third-party plugins. |

---

## 3. Real-World Architectural Patterns

### Pattern 1: Deterministic Multi-Line Heredoc Configuration

```hcl
resource "local_file" "bootstrap_script" {
  filename        = "${path.module}/scripts/bootstrap.sh"
  file_permission = "0755"

  # Indented heredoc (<<-EOT) strips leading whitespace before output
  content = <<-EOT
    #!/usr/bin/env bash
    set -euo pipefail

    echo "Bootstrapping service in ${path.module}..."
    export ENVIRONMENT="production"
  EOT
}
```

### Pattern 2: Graph Decoupling with `terraform_data`

```hcl
resource "terraform_data" "pipeline_trigger" {
  input = local_file.bootstrap_script.content_sha256

  lifecycle {
    triggers_replace = [
      local_file.bootstrap_script.content_sha256
    ]
  }
}
```

---

## 4. Production Hardening & Operational Governance

- **Pin Exact Provider Minor Versions**: Use pessimistic constraints (`~> 2.4.0`) in `required_providers` to accept patch releases while blocking breaking API modifications.
- **Enforce Canonical Formatting**: Always run `tofu fmt -check` or `terraform fmt -check` in continuous integration to eliminate diff noise.
- **Prefer Implicit Dependencies**: Reference attributes directly (e.g. `local_file.app_config.id`) instead of relying on explicit `depends_on`, allowing the DAG walker to optimize parallel execution.
- **Migrate from `null_resource`**: Replace all legacy `null_resource` instances with the native `terraform_data` primitive to eliminate external provider overhead.

---

## 5. Failure Modes & Diagnostic Triage Tree

??? failure "Error: Cycle in Dependency Graph (`Cycle: ...`)"
    **Root Cause:** Two or more resources reference each other's computed attributes directly or indirectly, creating an unresolvable circular dependency.

    **Diagnostic Triage Sequence:**
    1. Run `terraform graph | dot -Tpng > graph.png` to render the visual execution graph.
    2. Identify the circular link (e.g., `Resource A` references `Resource B` while `Resource B` references `Resource A`).
    3. Break the cycle by extracting shared data into a `locals` block or decoupling resources using `terraform_data`.

??? failure "Error: Unsupported argument / attribute"
    **Root Cause:** An argument or attribute is misspelled or does not exist in the pinned provider schema version.

    **Diagnostic Triage Sequence:**
    1. Inspect the exact provider version with `tofu version` or `terraform version`.
    2. Review provider schema docs for the target resource type.
    3. Check whether the field belongs inside a nested block or directly under the top-level resource block.

??? failure "Error: Provider configuration not present"
    **Root Cause:** A resource specifies a provider that has not been declared in `required_providers` or configured via a `provider` block.

    **Diagnostic Triage Sequence:**
    1. Check `required_providers` under `terraform { ... }` for the missing provider alias.
    2. Run `terraform init` or `tofu init` to download and link the required provider plugin.

---

## 6. Interactive Practice Matrix

Practice concepts from this chapter directly in the interactive WebAssembly sandbox:

| Exercise ID | Challenge Description | Direct Link | Action |
| :--- | :--- | :--- | :--- |
| **`primitives01`** | Terraform Configuration Block | [`../playground/index.html?exercise=primitives01`](../playground/index.html?exercise=primitives01) | [**⚡ Solve in Playground →**](../playground/index.html?exercise=primitives01){ .md-button .md-button--primary } |
| **`primitives02`** | First Resource Declaration | [`../playground/index.html?exercise=primitives02`](../playground/index.html?exercise=primitives02) | [**⚡ Solve in Playground →**](../playground/index.html?exercise=primitives02){ .md-button .md-button--primary } |
| **`primitives03`** | Resource Dependencies | [`../playground/index.html?exercise=primitives03`](../playground/index.html?exercise=primitives03) | [**⚡ Solve in Playground →**](../playground/index.html?exercise=primitives03){ .md-button .md-button--primary } |
| **`primitives04`** | String Interpolation & Heredoc | [`../playground/index.html?exercise=primitives04`](../playground/index.html?exercise=primitives04) | [**⚡ Solve in Playground →**](../playground/index.html?exercise=primitives04){ .md-button .md-button--primary } |
| **`primitives05`** | Syntax & Formatting | [`../playground/index.html?exercise=primitives05`](../playground/index.html?exercise=primitives05) | [**⚡ Solve in Playground →**](../playground/index.html?exercise=primitives05){ .md-button .md-button--primary } |
| **`primitives06`** | Lifecycle Mechanics | [`../playground/index.html?exercise=primitives06`](../playground/index.html?exercise=primitives06) | [**⚡ Solve in Playground →**](../playground/index.html?exercise=primitives06){ .md-button .md-button--primary } |
