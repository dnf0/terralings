# Chapter 04: Built-in Functions & Collections

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; String Manipulation, Collection Transforms, Filesystem IO, and Safe Evaluation (`try`/`can`)
-   :material-play-circle: **Interactive Challenges** &bull; 5 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Pure Functional Engine

Terraform and OpenTofu provide a rich library of **pure, deterministic built-in functions**. Functions cannot perform side effects or mutate remote resources; they transform in-memory data structures during plan evaluation.

```text
    ┌──────────────────────────────┐
    │     Raw Input Data           │
    └──────────────┬───────────────┘
                   │
                   ▼ (Pure Functional Transformation)
    ┌──────────────────────────────┐
    │  Built-in Functions          │
    │  • lower(), format(), trim() │
    │  • merge(), flatten(), slice│
    │  • jsondecode(), file()      │
    │  • try(), can() safe guards  │
    └──────────────┬───────────────┘
                   │
                   ▼
    ┌──────────────────────────────┐
    │ Transformed & Sanitized State│
    └──────────────────────────────┘
```

---

## 2. Annotated HCL Anatomy & Common Patterns

```hcl
locals {
  raw_config = file("${path.module}/config.json")
  parsed     = jsondecode(local.raw_config)

  # Collection manipulation and merging
  merged_tags = merge(
    { Tier = "standard" },
    local.parsed.custom_tags
  )

  # Safe property lookup with fallback
  database_port = try(local.parsed.database.port, 5432)
  is_valid_cidr = can(cidrnetmask("10.0.0.0/16"))
}
```

---

## 3. Production Best Practices

1. **Use `try()` for Graceful Fallbacks**: Avoid runtime plan crashes on missing optional dictionary keys by supplying default fallback values via `try(expr, fallback)`.
2. **Use `can()` for Validation Rules**: Guard variable validations with `can(regex(...))` or `can(cidrnetmask(...))` to verify format validity without raising unhandled errors.
3. **Prefer `jsondecode()` Over Complex Regex**: When parsing structured metadata, use native deserializers rather than brittle custom string slicing.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `functions01` | String Manipulation | `plan` | Master `lower()`, `upper()`, `format()`, and `trimspace()` formatting. |
| `functions02` | Collection Operations | `plan` | Use `merge()`, `flatten()`, `concat()`, and `distinct()`. |
| `functions03` | Encoding & Filesystem | `plan` | Read local templates with `file()` and parse structured JSON with `jsondecode()`. |
| `functions04` | Error Handling with `try` & `can` | `plan` | Implement defensive fallback evaluation using `try()` and predicate testing with `can()`. |
| `functions05` | Numeric & Math Functions | `plan` | Apply `min()`, `max()`, `ceil()`, and integer arithmetic expressions. |
