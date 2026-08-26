run "verify_staging_naming" {
  command = plan

  variables {
    environment  = "staging"
    service_name = "order-service"
  }

  # TODO: Fix the condition
  assert {
    condition     = output.service_id == "WRONG"
    error_message = "Staging service ID did not match expected pattern"
  }
}

# TODO: Add run "verify_prod_naming"
