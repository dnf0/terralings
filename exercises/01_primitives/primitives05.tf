# ==============================================================================
# Exercise: primitives05
# Chapter: 01_primitives (HCL Foundations & Core Primitives)
#
# Task:
# HCL is picky about syntax: block labels must not use commas, attribute
# assignments use `=`, and strings require double quotes (not single quotes).
#
# Fix the syntax errors in this configuration so that `validate` succeeds.
#
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

# TODO: Fix the syntax errors below (block declaration and string quote style)
resource "terraform_data", "formatted" {
  input = 'clean canonical hcl'
}
