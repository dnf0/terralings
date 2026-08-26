# ==============================================================================
# Solution: dynamic02
# Chapter: 06_dynamic_blocks (Dynamic Blocks & Advanced HCL)
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

  dynamic "source" {
    iterator = doc
    for_each = var.documentation_pages
    content {
      filename = doc.value.filename
      content  = doc.value.markdown
    }
  }
}

output "docs_sha" {
  value = data.archive_file.docs.output_sha
}
