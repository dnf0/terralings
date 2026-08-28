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

  # TODO (What): Look up the active environment configuration using local.environments[var.env].
  # TODO (Why): Map lookups eliminate brittle nested ternaries and centralize per-environment sizing matrices into clean data structures.
  current_env = null
}

resource "terraform_data" "app_cluster" {
  # TODO (What): Populate input using var.env and local.current_env attributes (instances, monitoring, tier).
  # TODO (Why): Passing pre-computed configuration objects simplifies resource definitions and improves testability.
  input = {
    environment = var.env
    instances   = local.current_env.instance_count
    monitoring  = local.current_env.enable_monitoring
    tier        = local.current_env.tier
  }
}

output "cluster_config" {
  value = terraform_data.app_cluster.output
}
