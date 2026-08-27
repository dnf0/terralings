# ==============================================================================
# Exercise: variables05
# Chapter: 02_variables (Input Variables, Types & Validations)
#
# Task:
# Custom validation blocks allow declaring constraints on variable values.
# Inside a `validation` block:
# - `condition`: A boolean expression referring only to the variable itself via `var.<name>`.
# - `error_message`: A message displayed to the user if `condition` evaluates to false.
#
# In this exercise:
# 1. Add validation to `environment` ensuring it is one of: ["dev", "staging", "prod"] using contains().
# 2. Add validation to `db_password` ensuring length(var.db_password) >= 12.
#
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "environment" {
  type        = string
  description = "Target environment: dev, staging, or prod"
  default     = "dev"

  validation {
    # TODO: Fix validation condition to refer to var.environment
    condition     = contains(["dev", "staging", "prod"], var.missing_env_var)
    error_message = "The environment variable must be one of: dev, staging, prod."
  }
}

variable "db_password" {
  type        = string
  description = "Database administrator password"
  default     = "SuperSecretPass123"
  sensitive   = true

  # TODO: Add validation block checking length(var.db_password) >= 12
}

resource "terraform_data" "app" {
  input = {
    env = var.environment
  }
}
