# I AM NOT DONE
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
# When done, remove the '# I AM NOT DONE' line at the top.
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
  # TODO: Compute safe_port with fallback 8080 using try()
  safe_port = 0

  # TODO: Compute parsed_timeout with fallback 30 using try(tonumber(...), 30)
  parsed_timeout = 0

  # TODO: Check if var.raw_data.database exists using can()
  has_database = false

  # TODO: Check if var.raw_data.payload_json is valid JSON using can(jsondecode(...))
  is_valid_json = false
}

resource "terraform_data" "safety" {
  input = {
    port         = local.safe_port
    timeout      = local.parsed_timeout
    has_db       = local.has_database
    valid_json   = local.is_valid_json
  }
}
