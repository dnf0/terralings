# ==============================================================================
# Solution: gov01
# Chapter: 13_governance (Architecture Governance & Enterprise Standards)
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
  workload_module = {
    cluster            = var.cluster_name
    service_name       = var.service_name
    log_retention_days = var.log_retention_days
    security_groups    = ["sg-workload-internal"]
  }
}

resource "terraform_data" "workload_module" {
  input = local.workload_module
}

output "workload_service" {
  value = terraform_data.workload_module.output.service_name
}
