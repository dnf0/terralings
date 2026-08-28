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
    # TODO (What): Set provider source to "hashicorp/local" and version to "~> 2.0".
    # TODO (Why): Terraform/OpenTofu requires explicit source addresses to resolve provider plugins
    #             from the registry, and version constraints guarantee reproducible behavior across environments.
    local = {
      source = ""
    }
  }
}
