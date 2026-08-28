# ==============================================================================
# Exercise: dynamic01
# Chapter: 06_dynamic_blocks (Dynamic Blocks & Advanced HCL)
#
# Task:
# A `dynamic` block acts like a `for` expression, but produces nested blocks
# inside top-level blocks (such as resources or data sources) instead of values.
# Inside a dynamic block, the iterator object (defaulting to the block name)
# exposes `.key` and `.value`.
#
# Complete the data "archive_file" block below:
# 1. Declare a `dynamic "source"` block iterating over `var.files`.
# 2. Inside the `content` block, set:
#    - `filename = source.key`
#    - `content  = source.value`
#
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

variable "files" {
  type = map(string)
  default = {
    "index.html" = "<html><body>Hello Terralings</body></html>"
    "style.css"  = "body { background-color: #1a1a1a; color: white; }"
    "app.js"     = "console.log('Terralings loaded');"
  }
}

data "archive_file" "bundle" {
  type        = "zip"
  output_path = "${path.module}/bundle.zip"

  # TODO (What): Declare a dynamic "source" block iterating over var.files with filename = source.key and content = source.value.
  # TODO (Why): Dynamic blocks generate repeating nested HCL blocks from collections programmatically without hardcoding each block.
  # dynamic "source" {
  #   for_each = var.files
  #   content {
  #     filename = ...
  #     content  = ...
  #   }
  # }
}

output "archive_sha" {
  value = data.archive_file.bundle.output_sha
}
