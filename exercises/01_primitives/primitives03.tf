# I AM NOT DONE
# ==============================================================================
# Exercise: primitives03
# Chapter: 01_primitives (HCL Foundations & Core Primitives)
#
# Task:
# Terraform and OpenTofu automatically build a Directed Acyclic Graph (DAG) of
# dependencies between resources based on attribute references.
#
# In this exercise, make the `second` resource depend on the `first` resource by
# referencing the `content` attribute of `local_file.first` in `local_file.second`.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
  }
}

resource "local_file" "first" {
  filename = "${path.module}/first.txt"
  content  = "Hello from the first file!"
}

resource "local_file" "second" {
  filename = "${path.module}/second.txt"
  # TODO: Reference local_file.first.content inside the string below
  content = "Second file contains: "
}
