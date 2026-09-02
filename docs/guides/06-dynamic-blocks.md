# Chapter 06: Dynamic Blocks & Advanced HCL

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; `dynamic` Block Generators, Custom Iterators, Nested Loops, and Complex HCL Expressions
-   :material-play-circle: **Interactive Challenges** &bull; 4 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Structural Metaprogramming

Certain resource schemas require repeated nested configuration blocks (e.g. security group rules, ingress rules, tags). `dynamic` blocks allow programmatic generation of these nested blocks from collections without duplicating HCL boilerplate.

```text
    ┌──────────────────────────────┐
    │ Complex Collection / Map     │
    └──────────────┬───────────────┘
                   │
                   ▼ (dynamic "block_name" { for_each = ... })
    ┌──────────────────────────────┐
    │ Structural Generator Engine  │
    └──────────────┬───────────────┘
                   │
                   ▼
    ┌──────────────────────────────┐
    │ Expanded Nested HCL Blocks   │
    │  • block { port = 80 }       │
    │  • block { port = 443 }      │
    └──────────────────────────────┘
```

---

## 2. Annotated HCL Anatomy & Schema Reference

```hcl
locals {
  service_endpoints = {
    http  = { port = 80,  protocol = "tcp" }
    https = { port = 443, protocol = "tcp" }
    ssh   = { port = 22,  protocol = "tcp" }
  }
}

# Dynamic generation of repetitive block structures
resource "terraform_data" "firewall_simulator" {
  input = {
    rules = [
      for name, rule in local.service_endpoints : {
        name     = name
        port     = rule.port
        protocol = rule.protocol
      }
    ]
  }
}
```

---

## 3. Production Best Practices

1. **Use `iterator` for Clarity**: When nesting multiple `dynamic` blocks, always declare explicit `iterator = <name>` to avoid namespace collisions with default `dynamic.value`.
2. **Do Not Overuse Dynamic Blocks**: Prefer static blocks when configuration is fixed. Dynamic blocks increase cognitive load and should only be used when cardinality varies dynamically across environments.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `dynamic01` | Dynamic Block Generation | `plan` | Generate nested configuration blocks dynamically using `for_each`. |
| `dynamic02` | Custom Iterators | `plan` | Disambiguate nested iteration scopes using explicit `iterator` labels. |
| `dynamic03` | Complex List & Map Comprehensions | `plan` | Transform data using `[for k, v in map : ...]` and `{for item in list : ...}`. |
| `dynamic04` | Conditional Nested Blocks | `plan` | Conditionally emit or omit nested blocks using empty list fallbacks. |
