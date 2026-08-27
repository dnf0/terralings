# ==============================================================================
# Exercise: data01
# Chapter: 07_data_sources (Data Sources & State Querying)
#
# Task:
# Data sources allow Terraform to use information defined outside of Terraform,
# defined by another separate Terraform configuration, or exported by functions.
# A `data` block requests data from a provider (such as the `local` provider's
# `local_file` data source).
#
# Complete the configuration below:
# 1. Define a `data "local_file" "manifest"` reading filename = "${path.module}/data01.tf".
# 2. In output "manifest_content_length", output `length(data.local_file.manifest.content)`.
#
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
  }
}

# TODO: Declare data "local_file" "manifest" reading "${path.module}/data01.tf"
# data "local_file" "manifest" {
#   filename = "${path.module}/data01.tf"
# }

output "manifest_content_length" {
  # TODO: Output length of data.local_file.manifest.content
  value = length(data.local_file.manifest.content)
}
