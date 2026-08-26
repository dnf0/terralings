override_resource {
  target = terraform_data.external_gateway
  values = {
    output = {
      status     = "healthy"
      latency_ms = 12
    }
  }
}

run "verify_mocked_gateway" {
  command = apply

  assert {
    condition     = output.gateway_status == "healthy" && output.latency_ms == 12
    error_message = "Gateway status and latency were not properly overridden"
  }
}
