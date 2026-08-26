run "test_invalid_port" {
  command = plan

  variables {
    port = -1
  }

  expect_failures = [
    var.port
  ]
}

run "test_valid_port" {
  command = plan

  variables {
    port = 443
  }

  assert {
    condition     = output.port_number == 443
    error_message = "Port number output was not set correctly"
  }
}
