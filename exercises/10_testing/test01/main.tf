# ==============================================================================
# Exercise: test01
# Chapter: 10_testing (Native Unit & Integration Testing)
#
# Task:
# The `.tftest.hcl` testing framework provides native unit testing for Terraform
# and OpenTofu modules.
#
# A `run` block specifies:
# - `command = plan` (for fast unit testing without creating resources) or `command = apply`
# - `variables { ... }` input overrides for the test stage
# - `assert { condition = ... error_message = ... }` assertion checks
#
# In this exercise:
# 1. Review `naming.tftest.hcl`.
# 2. Fix the assertion condition in `run "verify_staging_naming"` so it checks
#    that `output.service_id == "staging-order-service"`.
# 3. Add a second `run "verify_prod_naming"` testing `environment = "prod"` and
#    asserting `output.service_id == "prod-order-service"`.
#
# ==============================================================================

variable "environment" {
  type    = string
  default = "staging"
}

variable "service_name" {
  type    = string
  default = "order-service"
}

output "service_id" {
  value = "${var.environment}-${var.service_name}"
}
