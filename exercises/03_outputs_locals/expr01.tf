# ==============================================================================
# Exercise: expr01
# Chapter: 03_outputs_locals (Outputs, Locals & Expressions)
#
# Task:
# Ternary conditional expressions have the form:
#   `condition ? true_val : false_val`
#
# Both branches must produce compatible types.
#
# Complete the locals below using ternary expressions based on `var.is_production`:
# 1. `instance_type`: if `is_production` is true, use "m5.large", else "t3.micro"
# 2. `replica_count`: if `is_production` is true, use 3, else 1
#
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
  # TODO (What): Define instance_type = var.is_production ? "m5.large" : "t3.micro".
  # TODO (Why): Ternary expressions allow dynamic value selection based on boolean flags while keeping type schemas uniform.

  # TODO (What): Define replica_count = var.is_production ? 3 : 1.
  # TODO (Why): Dynamic sizing via conditionals lets a single configuration scale gracefully between dev and prod tiers.
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
