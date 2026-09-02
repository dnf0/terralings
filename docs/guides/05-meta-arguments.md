# Chapter 05: Meta-Arguments & Resource Scaling

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; `count`, `for_each`, `depends_on`, Lifecycle Hooks, and Target Control
-   :material-play-circle: **Interactive Challenges** &bull; 5 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Scaling Mechanics

Meta-arguments are special top-level parameters understood directly by the core engine rather than provider plugins. They control resource multiplicity, iteration, execution order, and replacement lifecycle behaviors.

```text
    ┌──────────────────────────────┐
    │     Resource Definition      │
    └──────────────┬───────────────┘
                   │
         ┌─────────┴─────────┐
         ▼                   ▼
    ┌──────────┐       ┌───────────┐
    │  count   │       │ for_each  │
    │ (Indexed)│       │ (Key-Val) │
    └────┬─────┘       └─────┬─────┘
         │                   │
         ▼                   ▼
  [res[0], res[1]]    [res["web"], res["api"]]
   (Fragile Shifts)    (Stable Identifiers)
```

---

## 2. Annotated HCL Anatomy & Schema Reference

```hcl
variable "environments" {
  type    = set(string)
  default = ["dev", "staging", "prod"]
}

# Key-based scaling with for_each
resource "local_file" "env_configs" {
  for_each = var.environments

  filename = "${path.module}/config-${each.key}.json"
  content  = jsonencode({
    tier = each.key
    id   = each.value
  })

  lifecycle {
    prevent_destroy       = false
    create_before_destroy = true
  }
}
```

---

## 3. Production Best Practices

1. **Prefer `for_each` Over `count`**: `count` relies on list index numbers. Removing an item from the middle of a list causes all subsequent resources to be destroyed and recreated. `for_each` keys provide stable, immutable addresses.
2. **Use `lifecycle.prevent_destroy` on Critical Data**: Guard databases, persistent volumes, and root configurations against accidental deletion.
3. **Minimize Explicit `depends_on`**: Use `depends_on` only when hidden ordering constraints exist that cannot be inferred via attribute references.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `meta01` | Resource Multiplicity with `count` | `plan` | Scale resources using `count` and `count.index`. |
| `meta02` | Key-Based Scaling with `for_each` | `plan` | Build stable resource sets using maps and sets with `each.key`/`each.value`. |
| `meta03` | Explicit Dependencies with `depends_on` | `plan` | Enforce sequential execution ordering where implicit dependencies do not exist. |
| `meta04` | Lifecycle Customizations | `plan` | Apply `create_before_destroy`, `prevent_destroy`, and `ignore_changes`. |
| `meta05` | Dynamic Replacement Triggers | `plan` | Force intentional resource recreation using `lifecycle.replace_triggered_by`. |
