# I AM NOT DONE
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
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

terraform {
  # TODO: Set required_version = ">= 1.6.0"
  required_version = ""

  required_providers {
    # TODO: Add local provider requirement
  }
}

resource "terraform_data" "registry_status" {
  # TODO: Set input map with status = "verified"
  input = {}
}

output "engine_status" {
  value = terraform_data.registry_status.output
}
