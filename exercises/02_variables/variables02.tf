# ==============================================================================
# Exercise: variables02
# Chapter: 02_variables (Input Variables, Types & Validations)
#
# Task:
# Collection types group multiple values of the same type:
# - `list(<TYPE>)`: Ordered sequential elements indexed by position: `var.list[0]`
# - `map(<TYPE>)`: Key-value pairs indexed by key: `var.map["key"]`
# - `set(<TYPE>)`: Unique, unordered elements
#
# Complete the variable definitions and fix the attribute references in `terraform_data`.
#
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "availability_zones" {
  type        = list(string)
  description = "Ordered list of availability zones"
  default     = ["zone-a", "zone-b", "zone-c"]
}

variable "tags" {
  type        = map(string)
  description = "Resource metadata tags"
  default = {
    project = "terralings"
    owner   = "student"
  }
}

variable "allowed_ips" {
  type        = set(string)
  description = "Unique set of whitelisted IP addresses"
  default     = ["192.168.1.1", "10.0.0.1"]
}

resource "terraform_data" "config_summary" {
  input = {
    # TODO: Index the first availability zone using [0]
    primary_zone = var.availability_zones[99]
    # TODO: Index the "project" key from tags using ["project"]
    project_tag = var.tags
    # TODO: Compute count of allowed_ips using length()
    ip_count = 0
  }
}
