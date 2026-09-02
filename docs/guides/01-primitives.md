# Chapter 01: HCL Foundations & Core Primitives

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; Syntax, Blocks, Providers, Implicit Dependencies, and Lifecycle Mechanics
-   :material-play-circle: **Interactive Challenges** &bull; 6 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Engine Mechanics

In Terraform and OpenTofu, **HashiCorp Configuration Language (HCL)** is a declarative configuration language engineered to express desired infrastructure state. The core engine constructs a Directed Acyclic Graph (DAG) of resources, computing execution plans based on explicit and implicit dependencies.

```text
    ┌──────────────────────────────┐
    │  Declarative HCL Config (.tf)│
    └──────────────┬───────────────┘
                   │
                   ▼ (Parser & AST Construction)
    ┌──────────────────────────────┐
    │  Implicit Dependency Graph   │
    │  (local_file.a ──► file.b)   │
    └──────────────┬───────────────┘
                   │
                   ▼ (Provider Protocol gRPC)
    ┌──────────────────────────────┐
    │      Provider Plugin         │ ────► [ Target Cloud / OS API ]
    │ (hashicorp/local, aws, etc.) │
    └──────────────────────────────┘
```

The engine reads all `.tf` files in the root module directory simultaneously, building the dependency graph before evaluating individual expressions or issuing API calls.

---

## 2. Annotated HCL Anatomy & Schema Reference

Below is a production-grade configuration demonstrating root blocks, provider requirements, resource dependencies, and lifecycle controls:

```hcl
terraform {
  required_version = ">= 1.6.0"

  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.4"
    }
  }
}

provider "local" {}

# Resource with implicit attribute reference
resource "local_file" "config" {
  filename        = "${path.module}/app.conf"
  content         = <<-EOT
    # Generated application configuration
    env = "production"
    log_level = "info"
  EOT
  file_permission = "0644"

  lifecycle {
    create_before_destroy = true
  }
}

# Downstream resource depending implicitly on local_file.config
resource "terraform_data" "notify" {
  input = local_file.config.id

  lifecycle {
    triggers_replace = [
      local_file.config.content
    ]
  }
}
```

### Key Syntax Elements

- **`terraform` block**: Declares engine version constraints and required providers with registry addresses.
- **`provider` block**: Configures credentials and global parameters for target provider plugins.
- **`resource` block**: Defines a managed infrastructure component (`type` + `name` label).
- **`${path.module}`**: Built-in filesystem path reference resolving to the current module directory.
- **`<<-EOT`**: Indented heredoc syntax stripping leading whitespace for clean multi-line text.

---

## 3. Production Best Practices

1. **Pin Provider Versions Strictly**: Always declare explicit pessimistic constraints (e.g. `~> 2.4`) in `required_providers` to prevent breaking upstream API updates.
2. **Favor Implicit Dependencies**: Reference attributes directly (e.g. `local_file.config.id`) rather than relying on explicit `depends_on`, allowing the graph engine to optimize parallelization.
3. **Use Indented Heredocs**: Use `<<-EOT` for templating and multi-line strings to avoid unwanted indentation in output files.
4. **Leverage `terraform_data`**: For arbitrary lifecycle triggering without external plugins, prefer built-in `terraform_data` over deprecated `null_resource`.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `primitives01` | Terraform Configuration Block | `validate` | Define `required_version` and `required_providers` with semantic constraints. |
| `primitives02` | First Resource Declaration | `plan` | Declare a `local_file` resource with required arguments and file permissions. |
| `primitives03` | Resource Dependencies | `plan` | Establish implicit dependency graphs through cross-resource attribute references. |
| `primitives04` | String Interpolation & Heredoc | `plan` | Master `${...}` interpolation and `<<-EOT` indented multi-line heredocs. |
| `primitives05` | Syntax & Formatting | `validate` | Fix syntax discrepancies and adhere to canonical `tofu fmt` standards. |
| `primitives06` | Lifecycle Mechanics | `plan` | Implement `terraform_data` replacement triggers with `triggers_replace`. |
