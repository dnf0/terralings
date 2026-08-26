# I AM NOT DONE
# ==============================================================================
# Exercise: meta04
# Chapter: 05_meta_arguments (Meta-Arguments & Resource Scaling)
#
# Task:
# The `lifecycle` block customizes how Terraform creates, updates, and destroys
# resources. Common lifecycle rules include:
# - `create_before_destroy = true`: Creates replacement resources before destroying
#   the old ones, preventing downtime on immutable resources.
# - `prevent_destroy = true`: Causes Terraform to reject with an error any plan that
#   would destroy the resource (useful for production databases).
#
# Complete the configuration below:
# 1. In resource "terraform_data" "zero_downtime_app", configure a `lifecycle` block
#    with `create_before_destroy = true`.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "app_version" {
  type    = string
  default = "v2.1.0"
}

resource "terraform_data" "zero_downtime_app" {
  input = var.app_version

  # TODO: Configure lifecycle block with create_before_destroy = true
  # lifecycle {
  #   create_before_destroy = ...
  # }
}

output "deployed_version" {
  value = terraform_data.zero_downtime_app.output
}
