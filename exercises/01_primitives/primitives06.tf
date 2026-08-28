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
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

resource "terraform_data" "deploy_token" {
  input = "v1.0.0"

  # TODO (What): Set triggers_replace = ["v1.0.0"] to track the deployment version string.
  # TODO (Why): terraform_data uses triggers_replace to force recreation of the resource whenever
  #             any tracked value in the collection changes during terraform plan/apply.
  triggers_replace = [var.missing_trigger_token]
}
