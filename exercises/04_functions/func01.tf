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
  # TODO (What): Compute normalized_env = lower(var.environment).
  # TODO (Why): lower() standardizes uppercase environment inputs to prevent casing mismatches in naming conventions.
  normalized_env = ""

  # TODO (What): Compute clean_service = lower(trimspace(var.service_name)).
  # TODO (Why): trimspace() strips leading/trailing whitespace before lowercasing, preventing invalid DNS hostnames.
  clean_service = ""

  # TODO (What): Compute hostname = format("%s-%s.internal.net", local.clean_service, local.normalized_env).
  # TODO (Why): format() provides printf-style string construction for predictable domain names and endpoint URIs.
  hostname = ""

  # TODO (What): Compute tag_csv = join(",", var.app_tags).
  # TODO (Why): join() serializes string lists into delimited string representations required by legacy environments and headers.
  tag_csv = ""

  # TODO (What): Compute slug = replace(local.hostname, ".", "-").
  # TODO (Why): replace() substitutes target characters (such as dots to hyphens) to produce valid filesystem or cloud resource identifier slugs.
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
