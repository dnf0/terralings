# ==============================================================================
# Solution: primitives03
# Chapter: 01_primitives (HCL Foundations & Core Primitives)
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
  content  = "Second file contains: ${local_file.first.content}"
}
