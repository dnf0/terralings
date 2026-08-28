# ==============================================================================
# Exercise: data03
# Chapter: 07_data_sources (Data Sources & State Querying)
#
# Task:
# The `external` data source allows an external program acting as a data source
# to be integrated into Terraform. The program must produce a single valid JSON
# object on stdout containing string keys and string values (map of strings).
#
# Complete the configuration below:
# 1. Declare `data "external" "system_info"` with:
#    `program = ["echo", "{\"status\": \"ready\", \"engine\": \"terralings\"}"]`
# 2. In output "system_status", output `data.external.system_info.result["status"]`.
#
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

# TODO (What): Declare data "external" "system_info" with program = ["echo", "{\"status\": \"ready\", \"engine\": \"terralings\"}"].
# TODO (Why): The external data source bridges external CLI scripts into Terraform, expecting a JSON object with string values on stdout.
# data "external" "system_info" {
#   program = ["echo", "{\"status\": \"ready\", \"engine\": \"terralings\"}"]
# }

output "system_status" {
  # TODO (What): Set value = data.external.system_info.result["status"].
  # TODO (Why): The external data source stores parsed script JSON output in its .result map attribute.
  value = data.external.system_info.result["status"]
}
