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
    # TODO (What): Change var.missing_env_var to var.environment in the condition.
    # TODO (Why): Validation conditions must refer strictly to the variable being declared (var.environment)
    #             so Terraform can validate input values early before downstream resource evaluation.
    condition     = contains(["dev", "staging", "prod"], var.missing_env_var)
    error_message = "The environment variable must be one of: dev, staging, prod."
  }
}

variable "db_password" {
  type        = string
  description = "Database administrator password"
  default     = "SuperSecretPass123"
  sensitive   = true

  # TODO (What): Add a validation block with condition = length(var.db_password) >= 12 and a helpful error_message.
  # TODO (Why): Enforcing password length constraints declaratively at variable ingestion prevents insecure database credentials
  #             from propagating into provisioning runs.
}

resource "terraform_data" "app" {
  input = {
    env = var.environment
  }
}
