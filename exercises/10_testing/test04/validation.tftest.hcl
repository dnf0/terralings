run "test_invalid_port" {
  command = plan

  variables {
    port = -1
  }

  # TODO (What): Add expect_failures = [var.port].
  # TODO (Why): expect_failures asserts that a negative test case properly triggers the declared custom validation condition on var.port.
}

# TODO (What): Add run "test_valid_port" with command = plan and port = 443.
# TODO (Why): Testing both the failure case and the success case guarantees complete coverage of variable validation boundaries.
