# Chapter 10: Native Unit & Integration Testing (.tftest.hcl)

<div class="grid cards" markdown>

-   :material-school: **Topic Focus** &bull; Native Testing Framework, `run` Blocks, Mock Providers, Assertions, and Negative Testing
-   :material-play-circle: **Interactive Challenges** &bull; 4 Hands-on Exercises
-   :material-rocket-launch: [**Launch Playground in Wasm →**](../playground/index.html){ .md-button .md-button--primary }

</div>

---

## 1. Architectural Overview & Test Runner Mechanics

Terraform &ge; 1.6.0 and OpenTofu feature a native, declarative testing framework (`tofu test` / `terraform test`). Tests are written in `.tftest.hcl` files, executing sequential `run` blocks against isolated ephemeral states with assertions and provider mocking.

```text
    ┌──────────────────────────────┐
    │     main.tf & variables.tf   │ (Module Under Test)
    └──────────────┬───────────────┘
                   │
                   ▼ (tofu test / terraform test)
    ┌──────────────────────────────┐
    │  Test Suite (tests/*.tftest) │
    │  • run "unit_validation"     │ (command = plan)
    │  • run "e2e_integration"     │ (command = apply)
    └──────────────┬───────────────┘
                   │
                   ▼ (Assertion Engine)
    ┌──────────────────────────────┐
    │  assert { condition = ... }  │ ──► [ Pass / Fail Report ]
    └──────────────────────────────┘
```

---

## 2. Annotated HCL Anatomy & Schema Reference

```hcl
# tests/validation.tftest.hcl

# Mock provider data/resources for fast unit tests without cloud credentials
mock_provider "local" {}

run "verify_port_assignment" {
  command = plan

  variables {
    environment = "staging"
  }

  assert {
    condition     = local_file.app_config.file_permission == "0644"
    error_message = "File permissions must be set to 0644."
  }
}

run "expect_invalid_environment_failure" {
  command = plan

  variables {
    environment = "invalid-tier"
  }

  expect_failures = [
    var.environment
  ]
}
```

---

## 3. Production Best Practices

1. **Test Plans Before Applies**: Run fast `command = plan` unit tests for syntax, variable validations, and output calculations before running slow infrastructure provisioning tests.
2. **Negative Testing with `expect_failures`**: Explicitly verify that invalid inputs correctly trigger expected validation errors.
3. **Use Mock Providers**: Isolate module logic from external cloud APIs using `mock_provider` for fast, reproducible CI/CD execution.

---

## 4. Hands-on Exercises in this Chapter

| Exercise ID | Name | Mode | Key Learning Objective |
|---|---|:---:|---|
| `testing01` | First Native Unit Test | `test` | Author a basic `.tftest.hcl` test suite with `run` blocks and assertions. |
| `testing02` | Plan Mode Assertions | `test` | Validate plan-time attribute values and output calculations without applying. |
| `testing03` | Mock Providers & Stubs | `test` | Mock provider interactions for instant offline test execution. |
| `testing04` | Negative Testing with `expect_failures` | `test` | Verify that constraint violations trigger intended custom error conditions. |
