# ==============================================================================
# Solution: pattern03
# Chapter: 11_patterns (Production Patterns & Anti-Patterns)
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

  standard_tags = merge(local.base_tags, local.env_tags)
}

resource "terraform_data" "payment_api" {
  input = {
    service_name = "payment-api"
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
  value = terraform_data.payment_api.input.tags["ManagedBy"]
}

output "effective_tags" {
  value = terraform_data.payment_api.output.tags
}
