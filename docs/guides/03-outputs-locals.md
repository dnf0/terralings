# Chapter 03: Outputs, Locals & Expressions

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; Outputs, Sensitive Masking, Local Expressions, Ternaries, and Splat Projections
-   :material-play-circle: **Interactive Challenges** &bull; 4 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Data Propagation

Outputs expose module state to root callers, remote state consumers, or CLI automation. Local values (`locals`) compute reusable intermediate expressions, keeping code DRY and declarative.

```text
    ┌──────────────────────────────┐
    │      Input Variables /       │
    │      Resource Attributes     │
    └──────────────┬───────────────┘
                   │
                   ▼
    ┌──────────────────────────────┐
    │        Local Values          │ ◄─── DRY Intermediate Calculations
    │      (locals { ... })        │      (Ternaries, String Ops, Splats)
    └──────────────┬───────────────┘
                   │
                   ▼
    ┌──────────────────────────────┐
    │     Module Outputs           │ ───► [ CLI / Remote State Consumers ]
    │  (sensitive = true masked)   │
    └──────────────────────────────┘
```

---

## 2. Annotated HCL Anatomy & Schema Reference

```hcl
locals {
  is_production = var.environment == "prod"
  instance_type = local.is_production ? "c6i.2xlarge" : "t3.micro"

  # Splat projection across resource lists
  all_server_ips = local_file.servers[*].id

  common_tags = {
    Environment = var.environment
    ManagedBy   = "Terralings"
  }
}

output "api_endpoint" {
  description = "Public API endpoint URL."
  value       = "https://${local_file.dns_record.content}:8443"
}

output "database_password" {
  description = "Root database administrative credential."
  value       = local_file.db_secret.content
  sensitive   = true # Masked from console stdout
}
```

---

## 3. Production Best Practices

1. **Protect Secrets with `sensitive = true`**: Any output referencing tokens, private keys, or passwords must have `sensitive = true` to prevent accidental logging in CI/CD stdout.
2. **Use Splat Expressions (`[*]`)**: Simplify list extractions (`resource[*].attribute`) instead of writing manual `for` loops when projecting single attributes.
3. **Keep Locals Declarative**: Avoid deeply nested, multi-tier chained locals that obscure dependency tracking.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `outputs01` | Basic Output Values | `plan` | Define root outputs with descriptive metadata and explicit values. |
| `outputs02` | Sensitive Output Redaction | `plan` | Secure confidential outputs using the `sensitive = true` flag. |
| `locals01` | DRY Expressions with Locals | `plan` | Consolidate repeated expressions and common tagging structures. |
| `expressions01`| Ternary & Splat Projections | `plan` | Write ternary conditional logic (`c ? t : f`) and splat operators (`[*]`). |
