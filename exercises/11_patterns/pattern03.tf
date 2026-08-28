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

  # TODO (What): Merge local.base_tags and local.env_tags using merge().
  # TODO (Why): Centralizing standardized baseline tags in locals enforces consistent cloud governance across all resources.
  standard_tags = {}
}

resource "terraform_data" "payment_api" {
  input = {
    service_name = "payment-api"
    # TODO (What): Merge local.standard_tags, var.custom_tags, and { Component = "payment-gateway" }.
    # TODO (Why): Multi-layer merge allows module callers to supply custom tags while preserving organizational governance defaults.
    tags = merge(
      local.standard_tags,
      var.custom_tags,
      {
        Component = "payment-gateway"
      }
    )
  }
}

output "managed_by_tag" {
  # TODO (What): Reference ManagedBy tag from terraform_data.payment_api.input.tags["ManagedBy"].
  # TODO (Why): Validates that merged governance tags are properly populated on the resource.
  value = terraform_data.payment_api.input.tags["ManagedBy"]
}

output "effective_tags" {
  value = terraform_data.payment_api.output.tags
}
