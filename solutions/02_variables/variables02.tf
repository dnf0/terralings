# ==============================================================================
# Solution: variables02
# Chapter: 02_variables (Input Variables, Types & Validations)
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
    primary_zone = var.availability_zones[0]
    project_tag  = var.tags["project"]
    ip_count     = length(var.allowed_ips)
  }
}
