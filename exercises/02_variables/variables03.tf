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
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "database_config" {
  # TODO (What): Add `ssl = optional(bool, true)` to the object type definition.
  # TODO (Why): The `optional(...)` modifier allows callers to omit the attribute, automatically falling back
  #             to the specified default value (true) without requiring boilerplate ternary expressions.
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
  type        = tuple([string, number])
  description = "Tuple of hostname and port"
  # TODO (What): Change the second element of the tuple from "not_a_number" to a number literal like 443.
  # TODO (Why): Tuples enforce positional type constraints; index 0 must be string and index 1 must be number.
  default     = ["api.internal", "not_a_number"]
}
