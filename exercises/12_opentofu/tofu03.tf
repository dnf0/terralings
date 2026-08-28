# ==============================================================================
# Exercise: tofu03
# Chapter: 12_opentofu (OpenTofu Innovations & Enterprise Features)
#
# Task:
# OpenTofu operates on an open, community-driven provider ecosystem backed by the
# OpenTofu registry. Provider declarations in `required_providers` ensure predictable
# constraint resolution and provider authenticity across CI and local environments.
#
# In this exercise:
# 1. Complete the `terraform` block with `required_version = ">= 1.6.0"`.
# 2. In `required_providers`, specify provider `local` with `source = "hashicorp/local"`
#    and `version = "~> 2.0"`.
# 3. Create a `terraform_data.registry_status` resource with input status `"verified"`.
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"

  required_providers {
    # TODO (What): Configure local provider with source = "hashicorp/local" and version = "~> 2.0".
    # TODO (Why): Explicit source and pessimistic version constraints (~> 2.0) ensure reproducible provider downloads across environments.
    local = {
      source = ""
    }
  }
}

resource "terraform_data" "registry_status" {
  # TODO (What): Set input = { status = "verified" }.
  # TODO (Why): Resource inputs record state values that verify successful provider and registry evaluation.
  input = {}
}

output "engine_status" {
  value = terraform_data.registry_status.output
}
