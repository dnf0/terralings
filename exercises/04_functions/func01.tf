# ==============================================================================
# Exercise: func01
# Chapter: 04_functions (Built-in Functions & Collections)
#
# Task:
# HCL provides rich built-in functions for string transformations.
# In this exercise, use the following string functions:
# - lower() and trimspace() to normalize strings
# - format() to construct formatted hostnames
# - join() to concatenate list elements with a delimiter
# - replace() to substitute characters
#
# Complete the locals block below:
# 1. normalized_env: convert var.environment to lowercase
# 2. clean_service: trim whitespace and lowercase var.service_name
# 3. hostname: format as "${clean_service}-${normalized_env}.internal.net"
# 4. tag_csv: join var.app_tags with ","
# 5. slug: replace all "." in hostname with "-"
#
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "environment" {
  type    = string
  default = "PRODUCTION"
}

variable "service_name" {
  type    = string
  default = "  Auth-Gateway  "
}

variable "app_tags" {
  type    = list(string)
  default = ["auth", "security", "v2"]
}

locals {
  # TODO: Compute normalized_env using lower()
  normalized_env = ""

  # TODO: Compute clean_service using lower() and trimspace()
  clean_service = ""

  # TODO: Compute hostname using format("%s-%s.internal.net", ...)
  hostname = ""

  # TODO: Compute tag_csv using join(",", ...)
  tag_csv = ""

  # TODO: Compute slug using replace(..., ".", "-")
  slug = ""
}

resource "terraform_data" "strings" {
  input = {
    env      = local.normalized_env
    service  = local.clean_service
    hostname = local.hostname
    tags     = local.tag_csv
    slug     = local.slug
  }

  lifecycle {
    postcondition {
      condition     = local.normalized_env == "production" && local.clean_service == "auth-gateway" && local.hostname == "auth-gateway-production.internal.net" && local.tag_csv == "auth,security,v2" && local.slug == "auth-gateway-production-internal-net"
      error_message = "All string transformations must match the expected format."
    }
  }
}
