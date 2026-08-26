# ==============================================================================
# Solution: locals01
# Chapter: 03_outputs_locals (Outputs, Locals & Expressions)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "project" {
  type    = string
  default = "terralings"
}

variable "environment" {
  type    = string
  default = "prod"
}

locals {
  service_prefix = "${var.project}-${var.environment}"
  common_tags = {
    Project     = var.project
    Environment = var.environment
    ManagedBy   = "terralings"
  }
}

resource "terraform_data" "service_instance" {
  input = {
    name = "${local.service_prefix}-backend"
    tags = local.common_tags
  }
}

output "instance_name" {
  value = terraform_data.service_instance.input.name
}
