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
  # TODO: Fix default value to match type string ("development")
  default     = 12345
  description = "Deployment target environment"
}

variable "port" {
  type        = number
  # TODO: Fix default value to match type number (8080)
  default     = "not-a-number"
  description = "Application server listening port"
}

variable "debug_mode" {
  type        = bool
  # TODO: Fix default value to match type bool (false)
  default     = "not-a-bool"
  description = "Flag to toggle verbose logging"
}
