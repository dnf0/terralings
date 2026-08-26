# I AM NOT DONE
# ==============================================================================
# Exercise: func04
# Chapter: 04_functions (Built-in Functions & Collections)
#
# Task:
# Terraform provides filesystem functions to inspect paths and load files:
# - fileexists(path) returns true if a file exists on disk
# - fileset(path, pattern) finds matching files within a directory
# - basename(path) returns the last element of a path (e.g. file name)
# - dirname(path) returns all but the last element of a path
# - file(path) reads the entire contents of a file as a string
#
# Complete the locals block below:
# 1. file_present: check if "${path.module}/func04.tf" exists with fileexists()
# 2. tf_files: find all "*.tf" files in path.module using fileset()
# 3. file_name: get the base name of "${path.module}/func04.tf" with basename()
# 4. file_length: calculate the character length of the content of "${path.module}/func04.tf" using length(file(...))
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

locals {
  target_file = "${path.module}/func04.tf"

  # TODO: Check if local.target_file exists using fileexists()
  file_present = false

  # TODO: List all *.tf files using fileset(path.module, "*.tf")
  tf_files = []

  # TODO: Extract file name using basename(local.target_file)
  file_name = ""

  # TODO: Calculate content length using length(file(local.target_file))
  file_length = 0
}

resource "terraform_data" "fs_info" {
  input = {
    exists = local.file_present
    files  = local.tf_files
    name   = local.file_name
    length = local.file_length
  }
}
