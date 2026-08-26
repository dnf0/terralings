# I AM NOT DONE
# ==============================================================================
# Exercise: pattern01
# Chapter: 11_patterns (Production Patterns & Anti-Patterns)
#
# Task:
# In production codebases, managing multi-environment configurations (dev, staging,
# prod) using brittle ternary chains or copy-pasted files is an anti-pattern.
#
# A clean, maintainable pattern is the Multi-Environment Map Lookup:
#
#   locals {
#     environments = {
#       dev = {
#         instance_count    = 1
#         enable_monitoring = false
#         tier              = "standard"
#       }
#       prod = {
#         instance_count    = 5
#         enable_monitoring = true
#         tier              = "premium"
#       }
#     }
#
#     current_env = local.environments[var.env]
#   }
#
# In this exercise:
# 1. Define `current_env` in `locals` by indexing `local.environments` with `var.env`.
# 2. Complete the `terraform_data.app_cluster` resource `input` to pass:
#    - `environment` = var.env
#    - `instances`   = local.current_env.instance_count
#    - `monitoring`  = local.current_env.enable_monitoring
#    - `tier`        = local.current_env.tier
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

variable "env" {
  type        = string
  description = "Target deployment environment (dev, staging, prod)"
  default     = "dev"

  validation {
    condition     = contains(["dev", "staging", "prod"], var.env)
    error_message = "The environment must be one of: dev, staging, prod."
  }
}

locals {
  environments = {
    dev = {
      instance_count    = 1
      enable_monitoring = false
      tier              = "standard"
    }
    staging = {
      instance_count    = 2
      enable_monitoring = true
      tier              = "standard"
    }
    prod = {
      instance_count    = 6
      enable_monitoring = true
      tier              = "premium"
    }
  }

  # TODO: Look up the active environment configuration using var.env as the key
  current_env = null
}

resource "terraform_data" "app_cluster" {
  # TODO: Populate the input map using var.env and local.current_env attributes
  input = {
    environment = var.env
  }
}

output "cluster_config" {
  value = terraform_data.app_cluster.output
}
