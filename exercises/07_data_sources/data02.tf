# I AM NOT DONE
# ==============================================================================
# Exercise: data02
# Chapter: 07_data_sources (Data Sources & State Querying)
#
# Task:
# The `archive_file` data source from the `hashicorp/archive` provider generates
# archive files (zip, tar.gz) directly in memory or on disk during plan/apply.
# This is widely used for packaging serverless Lambda functions and static assets.
#
# Complete the configuration below:
# 1. Define `data "archive_file" "package"` with:
#    - `type        = "zip"`
#    - `output_path = "${path.module}/package.zip"`
#    - A `source` block with:
#      - `content  = "console.log('App ready');"`
#      - `filename = "index.js"`
# 2. In output "archive_info", expose `output_sha` and `output_size`.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
  required_providers {
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.0"
    }
  }
}

# TODO: Declare data "archive_file" "package"
# data "archive_file" "package" {
#   type        = "zip"
#   output_path = "${path.module}/package.zip"
#   source {
#     content  = "console.log('App ready');"
#     filename = "index.js"
#   }
# }

output "archive_info" {
  # TODO: Output object with sha = data.archive_file.package.output_sha and size = data.archive_file.package.output_size
  value = {}
}
