# ==============================================================================
# Exercise: state04
# Chapter: 09_state_refactoring (State Management & Refactoring Surgery)
#
# Task:
# 1. OpenTofu and Terraform allow controlling replacement triggers using
#    `lifecycle { replace_triggered_by = [...] }`. This forces a resource to be
#    replaced whenever another resource or attribute changes.
# 2. Modern OpenTofu and Terraform also support `removed` blocks to safely
#    unmanage resources from state without configuration errors.
#
# In this exercise:
# 1. Add `replace_triggered_by = [terraform_data.app_version]` to `terraform_data.worker_fleet`.
# 2. Declare a `removed` block for `terraform_data.legacy_queue`.
#
# ==============================================================================

resource "terraform_data" "app_version" {
  input = "v2.4.0"
}

resource "terraform_data" "worker_fleet" {
  input = "fleet-runner"

  lifecycle {
    # TODO (What): Update replace_triggered_by = [terraform_data.app_version].
    # TODO (Why): replace_triggered_by establishes explicit replacement dependencies so instances recreate when dependent assets (like image tags or config hashes) change.
    replace_triggered_by = [
      terraform_data.missing_version
    ]
  }
}

# TODO (What): Add a removed block: removed { from = terraform_data.legacy_queue lifecycle { destroy = false } }.
# TODO (Why): Removed blocks cleanly unmanage retired resources from state without triggering unmanaged resource errors.
# removed {
#   from = terraform_data.legacy_queue
# }
