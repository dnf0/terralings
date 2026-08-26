# ==============================================================================
# Solution: func04
# Chapter: 04_functions (Built-in Functions & Collections)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

locals {
  target_file  = "${path.module}/func04.tf"
  file_present = fileexists(local.target_file)
  tf_files     = fileset(path.module, "*.tf")
  file_name    = basename(local.target_file)
  file_length  = length(file(local.target_file))
}

resource "terraform_data" "fs_info" {
  input = {
    exists = local.file_present
    files  = local.tf_files
    name   = local.file_name
    length = local.file_length
  }
}
