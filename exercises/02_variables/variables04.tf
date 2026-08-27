# ==============================================================================
# Exercise: variables04
# Chapter: 02_variables (Input Variables, Types & Validations)
#
# Task:
# By default, a variable without a default value is required.
# In addition, variables in modern Terraform / OpenTofu can specify `nullable = false`
# to reject explicit `null` values and guarantee non-null behavior.
#
# Fix the variables below by:
# 1. Setting default = "terralings-app" and nullable = false on `app_name`
# 2. Setting default = 2 and nullable = false on `replica_count`
#
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "app_name" {
  type        = string
  description = "Name of the application"
  # TODO: Add default = "terralings-app" and nullable = false
}

variable "replica_count" {
  type        = number
  description = "Number of worker replicas"
  # TODO: Add default = 2 and nullable = false
}

resource "terraform_data" "service" {
  input = {
    name     = var.app_name
    replicas = var.replica_count
  }
}
