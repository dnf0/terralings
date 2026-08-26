# I AM NOT DONE
# ==============================================================================
# Exercise: state04
# Chapter: 09_state_refactoring (State Management & Refactoring Surgery)
#
# Task:
# 1. OpenTofu and Terraform allow controlling replacement triggers using
#    `lifecycle { replace_triggered_by = [...] }`. This forces a resource to be
#    replaced whenever another resource or attribute changes.
# 2. Modern OpenTofu and Terraform also support `removed` blocks with
#    `lifecycle { destroy = false }` to safely unmanage resources from state
#    without destroying the real-world infrastructure.
#
# In this exercise:
# 1. Add `replace_triggered_by = [terraform_data.app_version]` to `terraform_data.worker_fleet`.
# 2. Declare a `removed` block for `terraform_data.legacy_queue` with `destroy = false`.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

resource "terraform_data" "app_version" {
  input = "v2.4.0"
}

resource "terraform_data" "worker_fleet" {
  input = "fleet-runner"

  # TODO: Add lifecycle with replace_triggered_by app_version
  # lifecycle {
  #   replace_triggered_by = [
  #     terraform_data.app_version
  #   ]
  # }
}

# TODO: Add removed block for legacy_queue with destroy = false
# removed {
#   from = terraform_data.legacy_queue
#   lifecycle {
#     destroy = false
#   }
# }
