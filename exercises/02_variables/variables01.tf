# ==============================================================================
# Exercise: variables01
# Chapter: 02_variables (Input Variables, Types & Validations)
#
# Task:
# Input variables serve as parameters for Terraform / OpenTofu modules.
# Primitive types include `string`, `number`, and `bool`.
#
# Declare the three input variables below with proper types, descriptions, and defaults:
# 1. `environment` (type: string, default: "development")
# 2. `port` (type: number, default: 8080)
# 3. `debug_mode` (type: bool, default: false)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "environment" {
  type        = string
  # TODO (What): Change default to a string value like "development".
  # TODO (Why): The type constraint is 'string'; assigning a numeric literal (12345) violates HCL type safety.
  default     = 12345
  description = "Deployment target environment"
}

variable "port" {
  type        = number
  # TODO (What): Change default to a number value like 8080.
  # TODO (Why): The type constraint is 'number'; assigning an unparsable string violates HCL type checking.
  default     = "not-a-number"
  description = "Application server listening port"
}

variable "debug_mode" {
  type        = bool
  # TODO (What): Change default to a boolean literal (false or true).
  # TODO (Why): The type constraint is 'bool'; booleans in HCL must evaluate directly to true or false.
  default     = "not-a-bool"
  description = "Flag to toggle verbose logging"
}
