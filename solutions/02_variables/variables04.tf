# ==============================================================================
# Solution: variables04
# Chapter: 02_variables (Input Variables, Types & Validations)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "app_name" {
  type        = string
  description = "Name of the application"
  default     = "terralings-app"
  nullable    = false
}

variable "replica_count" {
  type        = number
  description = "Number of worker replicas"
  default     = 2
  nullable    = false
}

resource "terraform_data" "service" {
  input = {
    name     = var.app_name
    replicas = var.replica_count
  }
}
