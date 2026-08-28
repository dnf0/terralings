# TODO (What): Add override_resource block targeting terraform_data.external_gateway with values = { output = { status = "healthy", latency_ms = 12 } }.
# TODO (Why): The override_resource block intercepts real resource synthesis, providing deterministic mock data without requiring live infrastructure or API credentials.
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
    # TODO (What): Assert output.gateway_status == "healthy" && output.latency_ms == 12.
    # TODO (Why): The assertion validates that the mock override values propagated accurately into the module's outputs.
    condition     = output.gateway_status == "healthy" && output.latency_ms == 12
    error_message = "Gateway status was not overridden by mock to healthy with 12ms latency"
  }
}
