# ==============================================================================
# Solution: primitives04
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

resource "local_file" "motd" {
  filename = "${path.module}/motd.txt"
  content  = <<-EOT
    ==============================
    Terralings System Status
    Environment: Local Learning
    Status: Operational
    ==============================
  EOT
}
