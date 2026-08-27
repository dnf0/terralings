# TODO: Add override_resource block for terraform_data.external_gateway
# override_resource {
#   target = terraform_data.external_gateway
#   values = {
#     output = {
#       status     = "healthy"
#       latency_ms = 12
#     }
#   }
# }

run "verify_mocked_gateway" {
  command = apply

  assert {
    # TODO: Assert gateway_status is healthy and latency_ms is 12
    condition     = output.gateway_status == "healthy" && output.latency_ms == 12
    error_message = "Gateway status was not overridden by mock to healthy with 12ms latency"
  }
}
