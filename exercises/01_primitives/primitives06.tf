# I AM NOT DONE
# ==============================================================================
# Exercise: primitives06
# Chapter: 01_primitives (HCL Foundations & Core Primitives)
#
# Task:
# The built-in `terraform_data` resource allows storing arbitrary values and
# triggering resource replacement when specific values change via `triggers_replace`.
#
# Fix the resource declaration below by adding `triggers_replace` to force
# replacement whenever the version token changes.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

resource "terraform_data" "deploy_token" {
  input = "v1.0.0"

  # TODO: Add triggers_replace = ["v1.0.0"]
}
