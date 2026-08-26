# ==============================================================================
# Solution: state04
# Chapter: 09_state_refactoring (State Management & Refactoring Surgery)
# ==============================================================================

resource "terraform_data" "app_version" {
  input = "v2.4.0"
}

resource "terraform_data" "worker_fleet" {
  input = "fleet-runner"

  lifecycle {
    replace_triggered_by = [
      terraform_data.app_version
    ]
  }
}

removed {
  from = terraform_data.legacy_queue
}
