# Chapter 09: State Management & Refactoring Surgery

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; Declarative `moved` Blocks, Config-Driven `import` Blocks, State Address Renaming, and Zero-Downtime Refactoring
-   :material-play-circle: **Interactive Challenges** &bull; 4 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & State Graph Evolution

Terraform/OpenTofu state acts as the single source of truth mapping HCL resource addresses to real-world cloud IDs. When refactoring code (e.g. moving a resource into a module or renaming it), declarative `moved` blocks instruct the engine to update state pointers without destroying and recreating resources.

```text
    ┌──────────────────────────────┐
    │ Existing State Address       │ (local_file.old_name)
    └──────────────┬───────────────┘
                   │
                   ▼ (Declarative moved { from = ... to = ... })
    ┌──────────────────────────────┐
    │ State Pointer Mutation       │ (No remote create/destroy API calls)
    └──────────────┬───────────────┘
                   │
                   ▼
    ┌──────────────────────────────┐
    │ New State Address            │ (module.app.local_file.new_name)
    └──────────────────────────────┘
```

---

## 2. Annotated HCL Anatomy & Schema Reference

```hcl
# Declarative resource rename without destroy/recreate
moved {
  from = local_file.legacy_config
  to   = local_file.modern_config
}

# Declarative resource move into a child module
moved {
  from = local_file.database
  to   = module.db.local_file.instance
}

# Config-driven resource import (Terraform 1.5+ / OpenTofu)
import {
  to = local_file.external_resource
  id = "/etc/existing-config.json"
}
```

---

## 3. Production Best Practices

1. **Never Edit State Manually**: Always use declarative `moved` blocks committed to Git for repeatable, code-reviewed state refactors.
2. **Commit `moved` Blocks Permanently in Modules**: Preserve `moved` blocks across multiple releases so consumers upgrading across version gaps transition state safely.
3. **Use `import` Blocks for Drift Remediation**: Adopt unmanaged infrastructure into code declaratively with `import` blocks rather than imperative CLI commands.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `state01` | Declarative Resource Renaming | `plan` | Refactor resource identifiers with zero downtime using `moved` blocks. |
| `state02` | Module Extraction Refactoring | `plan` | Move existing root resources into child modules with `moved` blocks. |
| `state03` | Config-Driven `import` Blocks | `plan` | Adopt existing infrastructure into state using declarative `import` syntax. |
| `state04` | Multi-Resource Collection Migration | `plan` | Migrate resources from single instances to `for_each` sets safely. |
