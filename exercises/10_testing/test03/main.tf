# I AM NOT DONE
# ==============================================================================
# Exercise: test03
# Chapter: 10_testing (Native Unit & Integration Testing)
#
# Task:
# Testing cloud infrastructure often involves external APIs or resources you want
# to mock in unit tests without making real API calls or needing credentials.
#
# You can use:
# 1. `mock_provider "local" { mock_data "..." { defaults = { ... } } }`
# 2. `override_resource { target = ... values = { ... } }`
#
# In this exercise:
# 1. In `mock.tftest.hcl`, use `override_resource` on `terraform_data.external_gateway`
#    to override `output = { status = "healthy", latency_ms = 12 }`.
# 2. Run assertions validating `output.gateway_status == "healthy"` and `output.latency_ms == 12`.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

resource "terraform_data" "external_gateway" {
  input = {
    status     = "connecting"
    latency_ms = 999
  }
}

output "gateway_status" {
  value = terraform_data.external_gateway.output.status
}

output "latency_ms" {
  value = terraform_data.external_gateway.output.latency_ms
}
