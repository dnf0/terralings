# I AM NOT DONE
# ==============================================================================
# Exercise: variables03
# Chapter: 02_variables (Input Variables, Types & Validations)
#
# Task:
# Structural types allow grouping values of different types into complex schemas:
# - `object({ ... })`: Complex structure with named attributes and optional defaults.
# - `tuple([...])`: Fixed-length sequence with specific types per position.
#
# Complete the object type constraint with `optional(bool, true)` for the `ssl` attribute
# and complete the tuple constraint.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "database_config" {
  # TODO: Define object type with:
  # host = string, port = number, username = string, ssl = optional(bool, true)
  type = object({
    host     = string
    port     = number
    username = string
  })
  description = "Database connection configuration"
  default = {
    host     = "localhost"
    port     = 5432
    username = "postgres"
  }
}

variable "endpoint_pair" {
  # TODO: Define tuple([string, number])
  type        = any
  description = "Tuple of hostname and port"
  default     = ["api.internal", 443]
}
