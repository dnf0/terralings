run "verify_staging_naming" {
  command = plan

  variables {
    environment  = "staging"
    service_name = "order-service"
  }

  # TODO (What): Fix the assertion condition to check output.service_id == "stg-order-service".
  # TODO (Why): The module formats staging service IDs with the "stg-" prefix according to company naming standards.
  assert {
    condition     = output.service_id == "WRONG"
    error_message = "Staging service ID did not match expected pattern"
  }
}

# TODO (What): Add run "verify_prod_naming" testing environment = "production" and asserting output.service_id == "prd-order-service".
# TODO (Why): Running separate run blocks validates multiple input permutations within the native test framework.
