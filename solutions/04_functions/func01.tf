# ==============================================================================
# Solution: func01
# Chapter: 04_functions (Built-in Functions & Collections)
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
  normalized_env = lower(var.environment)
  clean_service  = lower(trimspace(var.service_name))
  hostname       = format("%s-%s.internal.net", local.clean_service, local.normalized_env)
  tag_csv        = join(",", var.app_tags)
  slug           = replace(local.hostname, ".", "-")
}

resource "terraform_data" "strings" {
  input = {
    env      = local.normalized_env
    service  = local.clean_service
    hostname = local.hostname
    tags     = local.tag_csv
    slug     = local.slug
  }
}
