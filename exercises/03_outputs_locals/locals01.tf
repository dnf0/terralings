# ==============================================================================
# Exercise: locals01
# Chapter: 03_outputs_locals (Outputs, Locals & Expressions)
#
# Task:
# Local values (`locals { ... }`) assign names to expressions so they can be
# computed once and referenced multiple times throughout the configuration without
# repeating logic (DRY principle).
#
# Define:
# 1. `local.service_prefix` = "${var.project}-${var.environment}"
# 2. `local.common_tags` map containing Project, Environment, and ManagedBy ("terralings")
# 3. Reference `local.service_prefix` and `local.common_tags` in the resource below.
#
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
  # TODO (What): Define service_prefix = "${var.project}-${var.environment}".
  # TODO (Why): Centralizing naming prefixes in locals prevents string duplication and simplifies naming policy updates.

  # TODO (What): Define common_tags map with Project = var.project, Environment = var.environment, and ManagedBy = "terralings".
  # TODO (Why): Standardizing metadata tags in locals guarantees consistent resource tracking and cost allocation tags across infrastructure.
}

resource "terraform_data" "service_instance" {
  input = {
    # TODO (What): Reference local.service_prefix and local.common_tags in this resource.
    # TODO (Why): Referencing local values keeps resource declarations declarative and maintainable.
    name = "${local.service_prefix}-backend"
    tags = local.common_tags
  }
}

output "instance_name" {
  value = terraform_data.service_instance.input.name
}
