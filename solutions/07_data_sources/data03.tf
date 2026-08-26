# ==============================================================================
# Solution: data03
# Chapter: 07_data_sources (Data Sources & State Querying)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
  required_providers {
    external = {
      source  = "hashicorp/external"
      version = "~> 2.0"
    }
  }
}

data "external" "system_info" {
  program = ["echo", "{\"status\": \"ready\", \"engine\": \"terralings\"}"]
}

output "system_status" {
  value = data.external.system_info.result["status"]
}
