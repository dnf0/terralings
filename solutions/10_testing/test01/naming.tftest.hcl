run "verify_staging_naming" {
  command = plan

  variables {
    environment  = "staging"
    service_name = "order-service"
  }

  assert {
    condition     = output.service_id == "staging-order-service"
    error_message = "Staging service ID did not match expected pattern"
  }
}

run "verify_prod_naming" {
  command = plan

  variables {
    environment  = "prod"
    service_name = "order-service"
  }

  assert {
    condition     = output.service_id == "prod-order-service"
    error_message = "Production service ID did not match expected pattern"
  }
}
