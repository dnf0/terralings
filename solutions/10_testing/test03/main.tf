# ==============================================================================
# Solution: test03
# Chapter: 10_testing (Native Unit & Integration Testing)
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
