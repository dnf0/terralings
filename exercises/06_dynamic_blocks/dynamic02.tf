# ==============================================================================
# Exercise: dynamic02
# Chapter: 06_dynamic_blocks (Dynamic Blocks & Advanced HCL)
#
# Task:
# By default, a dynamic block uses the label of the block (e.g. `source`) as
# the iterator variable name. You can customize this name using the `iterator`
# argument to make code clearer, especially when working with descriptive types.
#
# Complete the configuration below:
# 1. In `data "archive_file" "docs"`, define `dynamic "source"`.
# 2. Add `iterator = doc` inside the dynamic block.
# 3. Set `for_each = var.documentation_pages`.
# 4. In `content`, set `filename = doc.value.filename` and `content = doc.value.markdown`.
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

variable "documentation_pages" {
  type = list(object({
    filename = string
    markdown = string
  }))
  default = [
    {
      filename = "README.md"
      markdown = "# Welcome to Terralings\nAn interactive learning environment."
    },
    {
      filename = "CONTRIBUTING.md"
      markdown = "# Contributing Guidelines\nPlease submit pull requests to main."
    }
  ]
}

data "archive_file" "docs" {
  type        = "zip"
  output_path = "${path.module}/docs.zip"

  # TODO (What): Declare dynamic "source" with iterator = doc, for_each = var.documentation_pages, and content with filename = doc.value.filename and content = doc.value.markdown.
  # TODO (Why): The iterator argument overrides the default block-name iterator variable to provide clean, readable symbol names in nested scopes.
  # dynamic "source" {
  #   iterator = doc
  #   for_each = var.documentation_pages
  #   content {
  #     filename = ...
  #     content  = ...
  #   }
  # }
}

output "docs_sha" {
  value = data.archive_file.docs.output_sha
}
