# ==============================================================================
# Exercise: gov01
# Chapter: 13_governance (Architecture Governance & Enterprise Standards)
#
# Task:
# A foundational architectural standard in enterprise OpenTofu & Terraform is
# Root Module Encapsulation:
#   - Zero plain/loose workload compute resources in root environments.
#   - Workload compute primitives (task definitions, compute instances, log groups
#     with explicit retention, security groups, and IAM roles) belong inside
#     dedicated workload modules, never declared as loose resources in root.
#   - Root configurations (`environments/prod/main.tf`) must serve exclusively
#     as orchestrators wiring modules together.
#
# In this exercise:
# 1. In `locals`, structure the cohesive `workload_module` definition bundling:
#    - `cluster`            = var.cluster_name
#    - `service_name`       = var.service_name
#    - `log_retention_days` = var.log_retention_days
#    - `security_groups`    = ["sg-workload-internal"]
# 2. Instantiate `terraform_data.workload_module` passing the encapsulated configuration.
# 3. Output `workload_service` referencing the service name from the module.
#
# ==============================================================================

variable "cluster_name" {
  type        = string
  description = "Target compute cluster"
  default     = "prod-ecs-cluster"
}

variable "service_name" {
  type        = string
  description = "Workload service name"
  default     = "orders-api"
}

variable "log_retention_days" {
  type        = number
  description = "Explicit log group retention in days"
  default     = 30
}

locals {
  # TODO: Encapsulate workload compute primitives into a structured module payload
  workload_module = {}
}

resource "terraform_data" "workload_module" {
  # TODO: Set input to local.workload_module
  input = local.workload_module

  lifecycle {
    postcondition {
      condition     = lookup(self.output, "service_name", "") == var.service_name && lookup(self.output, "cluster", "") == var.cluster_name && length(lookup(self.output, "security_groups", [])) == 1
      error_message = "Workload compute module must encapsulate cluster, service_name, log retention, and security groups."
    }
  }
}

output "workload_service" {
  # TODO: Reference service_name from terraform_data.workload_module.output
  value = terraform_data.workload_module.output.service_name
}
