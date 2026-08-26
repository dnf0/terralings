# ==============================================================================
# Solution: expr01
# Chapter: 03_outputs_locals (Outputs, Locals & Expressions)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "is_production" {
  type        = bool
  description = "Production environment toggle"
  default     = true
}

locals {
  instance_type = var.is_production ? "m5.large" : "t3.micro"
  replica_count = var.is_production ? 3 : 1
}

resource "terraform_data" "cluster_config" {
  input = {
    type     = local.instance_type
    replicas = local.replica_count
  }
}

output "cluster_type" {
  value = terraform_data.cluster_config.input.type
}
