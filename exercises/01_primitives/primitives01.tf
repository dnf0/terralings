# I AM NOT DONE
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
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"

  # TODO: Declare the required_providers block for the local provider
  # required_providers {
  #   local = {
  #     source  = "hashicorp/local"
  #     version = "~> 2.0"
  #   }
  # }
}
