# ==============================================================================
# Solution: func05
# Chapter: 04_functions (Built-in Functions & Collections)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "raw_data" {
  type = any
  default = {
    network = {
      host = "localhost"
    }
    timeout_str  = "invalid_number"
    payload_json = "{\"status\": \"ok\"}"
  }
}

locals {
  safe_port      = try(var.raw_data.network.port, 8080)
  parsed_timeout = try(tonumber(var.raw_data.timeout_str), 30)
  has_database   = can(var.raw_data.database)
  is_valid_json  = can(jsondecode(var.raw_data.payload_json))
}

resource "terraform_data" "safety" {
  input = {
    port       = local.safe_port
    timeout    = local.parsed_timeout
    has_db     = local.has_database
    valid_json = local.is_valid_json
  }
}
