# ==============================================================================
# Exercise: expr02
# Chapter: 03_outputs_locals (Outputs, Locals & Expressions)
#
# Task:
# Splat expressions (`[*]`) provide a concise way to extract attributes across
# a list of objects or complex structures without writing a full `for` comprehension.
#
# Given the `servers` list of objects below:
# 1. Define `all_ips` extracting `ip_address` across all servers using splat syntax: `local.servers[*].ip_address`
# 2. Define `all_hostnames` extracting `hostname` across all servers using splat syntax: `local.servers[*].hostname`
#
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

  # TODO (What): Define all_ips = local.servers[*].ip_address.
  # TODO (Why): The splat expression [*] projects a single attribute across all elements of a list without verbose for loops.

  # TODO (What): Define all_hostnames = local.servers[*].hostname.
  # TODO (Why): Splat projections provide clean, idiomatic attribute extraction over collections of uniform objects.
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
