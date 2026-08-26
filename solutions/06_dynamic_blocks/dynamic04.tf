# ==============================================================================
# Solution: dynamic04
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

variable "base_files" {
  type = map(string)
  default = {
    "server.js" = "require('http').createServer().listen(8080);"
  }
}

variable "include_debug" {
  type        = bool
  description = "Flag to conditionally include debug log artifact"
  default     = true
}

variable "debug_content" {
  type    = string
  default = "VERBOSE=true; DEBUG_LEVEL=trace"
}

data "archive_file" "package" {
  type        = "zip"
  output_path = "${path.module}/package.zip"

  dynamic "source" {
    for_each = var.base_files
    content {
      filename = source.key
      content  = source.value
    }
  }

  dynamic "source" {
    for_each = var.include_debug ? [{ filename = "debug.log", content = var.debug_content }] : []
    content {
      filename = source.value.filename
      content  = source.value.content
    }
  }
}

output "archive_sha" {
  value = data.archive_file.package.output_sha
}
