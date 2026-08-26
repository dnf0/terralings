# I AM NOT DONE
# ==============================================================================
# Exercise: primitives04
# Chapter: 01_primitives (HCL Foundations & Core Primitives)
#
# Task:
# HCL supports multi-line strings using heredoc syntax (`<<EOT` or `<<-EOT`).
# The indented heredoc (`<<-EOT`) strips leading whitespace matching the closing
# delimiter line.
#
# Fix the broken heredoc syntax below so that the resource validates and plans cleanly.
#
# When done, remove the '# I AM NOT DONE' line at the top.
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
  # TODO: Fix the heredoc opening and closing delimiters below
  content  = <<-EOT
    ==============================
    Terralings System Status
    Environment: Local Learning
    Status: Operational
    ==============================
  # Missing proper EOT delimiter!
}
