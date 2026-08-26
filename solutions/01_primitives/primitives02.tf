# ==============================================================================
# Solution: primitives02
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

resource "local_file" "welcome" {
  filename = "${path.module}/welcome.txt"
  content  = "Welcome to Terralings!"
}
