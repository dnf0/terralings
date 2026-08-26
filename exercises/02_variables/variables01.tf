# I AM NOT DONE
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
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "environment" {
  # TODO: Set type = string and default = "development"
  description = "Deployment target environment"
}

variable "port" {
  # TODO: Set type = number and default = 8080
  description = "Application server listening port"
}

variable "debug_mode" {
  # TODO: Set type = bool and default = false
  description = "Flag to toggle verbose logging"
}
