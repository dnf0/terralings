# ==============================================================================
# Exercise: primitives01
# Chapter: 01_primitives (HCL Foundations & Core Primitives)
#
# Task:
# Every Terraform / OpenTofu project begins with a `terraform` configuration
# block. This block specifies the minimum required engine version and declares
# the providers needed by your configuration.
#
# In this exercise, complete the `terraform` block to require version ">= 1.6.0"
# and declare the `local` provider from "hashicorp/local" with version "~> 2.0".
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"

  required_providers {
    # TODO: Specify provider source "hashicorp/local" and version "~> 2.0"
    local = {
      source = ""
    }
  }
}
