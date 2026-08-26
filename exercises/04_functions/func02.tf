# I AM NOT DONE
# ==============================================================================
# Exercise: func02
# Chapter: 04_functions (Built-in Functions & Collections)
#
# Task:
# HCL provides powerful functions for manipulating collections (lists, maps):
# - merge(m1, m2, ...) combines multiple maps (later arguments override earlier)
# - lookup(map, key, default) retrieves a key or returns fallback if absent
# - distinct(list) removes duplicate elements while preserving order
# - slice(list, start, end) extracts a subset of a list
# - flatten(nested_list) recursively flattens nested lists into a single list
# - zipmap(keys, values) builds a map from paired lists of equal length
#
# Complete the locals block below:
# 1. merged_tags: merge var.default_tags with var.custom_tags
# 2. owner: lookup "Owner" in local.merged_tags, defaulting to "Platform"
# 3. unique_zones: remove duplicates from var.raw_zones with distinct()
# 4. primary_zones: slice the first 2 elements (index 0 to 2) from unique_zones
# 5. flat_subnets: flatten var.nested_subnets into a 1D list
# 6. config_map: zipmap var.keys_list and var.values_list
#
# When done, remove the '# I AM NOT DONE' line at the top.
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
  # TODO: Compute merged_tags using merge()
  merged_tags = {}

  # TODO: Compute owner using lookup(local.merged_tags, "Owner", "Platform")
  owner = ""

  # TODO: Compute unique_zones using distinct()
  unique_zones = []

  # TODO: Compute primary_zones using slice(local.unique_zones, 0, 2)
  primary_zones = []

  # TODO: Compute flat_subnets using flatten()
  flat_subnets = []

  # TODO: Compute config_map using zipmap()
  config_map = {}
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
