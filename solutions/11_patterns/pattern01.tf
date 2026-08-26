# ==============================================================================
# Solution: pattern01
# Chapter: 11_patterns (Production Patterns & Anti-Patterns)
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

  current_env = local.environments[var.env]
}

resource "terraform_data" "app_cluster" {
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
