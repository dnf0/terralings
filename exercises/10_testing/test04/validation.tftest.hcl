run "test_invalid_port" {
  command = plan

  variables {
    port = -1
  }

  # TODO: Add expect_failures = [var.port]
}

# TODO: Add run "test_valid_port" with port = 443
