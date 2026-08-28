# ==============================================================================
# Exercise: func05
# Chapter: 04_functions (Built-in Functions & Collections)
#
# Task:
# Dynamic and loosely-typed data structures may have missing attributes or
# invalid conversions. HCL provides try() and can() to handle these safely:
# - try(expr1, expr2, ... fallback) evaluates expressions in order and returns
#   the result of the first one that does not produce an error.
# - can(expr) evaluates an expression and returns true if no error occurs, false otherwise.
#
# Complete the locals block below:
# 1. safe_port: try to read var.raw_data.network.port, fallback to 8080
# 2. parsed_timeout: try to parse tonumber(var.raw_data.timeout_str), fallback to 30
# 3. has_database: check if var.raw_data.database exists and is accessible using can()
# 4. is_valid_json: check if jsondecode(var.raw_data.payload_json) succeeds using can()
#
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "raw_data" {
  type = any
  default = {
    network = {
      host = "localhost"
      # port is missing!
    }
    timeout_str  = "invalid_number"
    payload_json = "{\"status\": \"ok\"}"
    # database is missing!
  }
}

locals {
  # TODO (What): Compute safe_port = try(var.raw_data.network.port, 8080).
  # TODO (Why): try() evaluates expressions sequentially and catches missing attribute errors, safely applying a default value.
  safe_port = 0

  # TODO (What): Compute parsed_timeout = try(tonumber(var.raw_data.timeout_str), 30).
  # TODO (Why): tonumber() fails when parsing non-numeric strings; wrapping it in try() guarantees a valid numeric fallback.
  parsed_timeout = 0

  # TODO (What): Compute has_database = can(var.raw_data.database).
  # TODO (Why): can() tests if an attribute access expression evaluates without error, returning a clean boolean for feature gating.
  has_database = false

  # TODO (What): Compute is_valid_json = can(jsondecode(var.raw_data.payload_json)).
  # TODO (Why): can(jsondecode(...)) performs schema/syntax validation on JSON payloads without terminating plan execution if invalid.
  is_valid_json = false
}

resource "terraform_data" "safety" {
  input = {
    port       = local.safe_port
    timeout    = local.parsed_timeout
    has_db     = local.has_database
    valid_json = local.is_valid_json
  }

  lifecycle {
    postcondition {
      condition     = local.safe_port == 8080 && local.parsed_timeout == 30 && !local.has_database && local.is_valid_json
      error_message = "try() and can() safety evaluations must produce expected fallbacks and boolean results."
    }
  }
}
