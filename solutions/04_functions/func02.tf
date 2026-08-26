# ==============================================================================
# Solution: func02
# Chapter: 04_functions (Built-in Functions & Collections)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "default_tags" {
  type = map(string)
  default = {
    "Environment" = "dev"
    "ManagedBy"   = "terralings"
  }
}

variable "custom_tags" {
  type = map(string)
  default = {
    "Environment" = "prod"
    "CostCenter"  = "cc-100"
  }
}

variable "raw_zones" {
  type    = list(string)
  default = ["us-east-1a", "us-east-1b", "us-east-1a", "us-east-1c", "us-east-1b"]
}

variable "nested_subnets" {
  type    = list(list(string))
  default = [["10.0.1.0/24", "10.0.2.0/24"], ["10.0.3.0/24"]]
}

variable "keys_list" {
  type    = list(string)
  default = ["db_host", "db_port"]
}

variable "values_list" {
  type    = list(string)
  default = ["10.0.0.5", "5432"]
}

locals {
  merged_tags   = merge(var.default_tags, var.custom_tags)
  owner         = lookup(local.merged_tags, "Owner", "Platform")
  unique_zones  = distinct(var.raw_zones)
  primary_zones = slice(local.unique_zones, 0, 2)
  flat_subnets  = flatten(var.nested_subnets)
  config_map    = zipmap(var.keys_list, var.values_list)
}

resource "terraform_data" "collections" {
  input = {
    tags    = local.merged_tags
    owner   = local.owner
    zones   = local.primary_zones
    subnets = local.flat_subnets
    config  = local.config_map
  }
}
