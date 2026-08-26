# ==============================================================================
# Solution: variables05
# Chapter: 02_variables (Input Variables, Types & Validations)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "environment" {
  type        = string
  description = "Target environment: dev, staging, or prod"
  default     = "dev"

  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "The environment variable must be one of: dev, staging, prod."
  }
}

variable "db_password" {
  type        = string
  description = "Database administrator password"
  default     = "SuperSecretPass123"
  sensitive   = true

  validation {
    condition     = length(var.db_password) >= 12
    error_message = "The db_password must be at least 12 characters in length."
  }
}

resource "terraform_data" "app" {
  input = {
    env = var.environment
  }
}
