# ==============================================================================
# Exercise: test04
# Chapter: 10_testing (Native Unit & Integration Testing)
#
# Task:
# Testing should not only cover happy paths, but also verify that invalid inputs
# and failed preconditions are properly rejected.
#
# The `expect_failures` argument inside a `run` block asserts that an execution
# is expected to fail with errors on specific variables, resources, or check blocks:
#
#   run "test_invalid_port" {
#     command = plan
#     variables {
#       port = 99999
#     }
#     expect_failures = [
#       var.port
#     ]
#   }
#
# In this exercise:
# 1. In `validation.tftest.hcl`, configure `run "test_invalid_port"` with
#    `port = -1` and `expect_failures = [var.port]`.
# 2. Add a passing `run "test_valid_port"` with `port = 443` asserting
#    `output.port_number == 443`.
#
# ==============================================================================

variable "port" {
  type        = number
  description = "Port number between 1 and 65535"

  validation {
    condition     = var.port >= 1 && var.port <= 65535
    error_message = "port must be between 1 and 65535"
  }
}

output "port_number" {
  value = var.port
}
