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
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

locals {
  target_file = "${path.module}/func04.tf"

  # TODO (What): Compute file_present = fileexists(local.target_file).
  # TODO (Why): fileexists() tests disk presence without failing if the file is absent, enabling conditional logic.
  file_present = false

  # TODO (What): Compute tf_files = fileset(path.module, "*.tf").
  # TODO (Why): fileset() discovers matching files dynamically across directories for bulk template rendering or file deployment.
  tf_files = []

  # TODO (What): Compute file_name = basename(local.target_file).
  # TODO (Why): basename() isolates the filename from directory paths for logging and output keys.
  file_name = ""

  # TODO (What): Compute file_length = length(file(local.target_file)).
  # TODO (Why): file() reads raw content into memory, allowing length() to calculate size or verify non-empty file contents.
  file_length = 0
}

resource "terraform_data" "fs_info" {
  input = {
    exists = local.file_present
    files  = local.tf_files
    name   = local.file_name
    length = local.file_length
  }

  lifecycle {
    postcondition {
      condition     = local.file_present && length(local.tf_files) > 0 && local.file_name == "func04.tf" && local.file_length > 0
      error_message = "Filesystem functions must return valid metadata for func04.tf."
    }
  }
}
