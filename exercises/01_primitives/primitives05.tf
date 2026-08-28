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

# TODO (What): Remove the comma separating block labels ("terraform_data" "formatted") and wrap the string in double quotes ("clean canonical hcl").
# TODO (Why): HCL syntax uses whitespace to separate block type and name labels, and standard HCL strings only accept double quotes.
resource "terraform_data", "formatted" {
  input = 'clean canonical hcl'
}
