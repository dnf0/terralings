# ==============================================================================
# Solution: expr02
# Chapter: 03_outputs_locals (Outputs, Locals & Expressions)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

locals {
  servers = [
    { hostname = "web-01", ip_address = "10.0.1.10", role = "web" },
    { hostname = "web-02", ip_address = "10.0.1.11", role = "web" },
    { hostname = "db-01", ip_address = "10.0.2.20", role = "db" },
  ]

  all_ips       = local.servers[*].ip_address
  all_hostnames = local.servers[*].hostname
}

resource "terraform_data" "network_directory" {
  input = {
    ips   = local.all_ips
    hosts = local.all_hostnames
  }
}

output "ip_addresses" {
  value = local.all_ips
}
