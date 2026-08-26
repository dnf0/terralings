# I AM NOT DONE
# ==============================================================================
# Exercise: pattern03
# Chapter: 11_patterns (Production Patterns & Anti-Patterns)
#
# Task:
# Organizations require consistent metadata tagging for ownership, cost attribution,
# compliance, and automated lifecycle management.
#
# The Tagging Factory Pattern creates a single standard tags map in `locals` by
# merging baseline governance tags, environment tags, and component-specific tags:
#
#   locals {
#     base_tags = {
#       ManagedBy  = "terralings"
#       CostCenter = var.cost_center
#     }
#     env_tags = {
#       Environment = var.environment
#     }
#     standard_tags = merge(local.base_tags, local.env_tags)
#   }
#
# In this exercise:
# 1. In `locals`, compute `standard_tags` by merging `local.base_tags` and `local.env_tags`.
# 2. In `terraform_data.payment_api.input.tags`, merge `local.standard_tags` with
#    `var.custom_tags` and resource-specific tags `{ Component = "payment-gateway" }`.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

variable "environment" {
  type        = string
  description = "Target deployment environment"
  default     = "production"
}

variable "cost_center" {
  type        = string
  description = "Cost allocation center code"
  default     = "fintech-core-802"
}

variable "custom_tags" {
  type        = map(string)
  description = "Caller-supplied custom tags"
  default = {
    Team = "Billing"
  }
}

locals {
  base_tags = {
    ManagedBy  = "terralings"
    CostCenter = var.cost_center
  }

  env_tags = {
    Environment = var.environment
  }

  # TODO: Merge local.base_tags and local.env_tags
  standard_tags = {}
}

resource "terraform_data" "payment_api" {
  input = {
    service_name = "payment-api"
    # TODO: Merge local.standard_tags, var.custom_tags, and { Component = "payment-gateway" }
    tags = {}
  }
}

output "effective_tags" {
  value = terraform_data.payment_api.output.tags
}
