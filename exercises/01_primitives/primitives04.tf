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
  # TODO (What): Close the multi-line heredoc string with the 'EOT' delimiter on a dedicated line.
  # TODO (Why): Indented heredocs (<<-EOT) strip leading whitespace for clean formatting, but HCL's parser
  #             requires the matching terminator (EOT) to know where the string payload ends.
  content  = <<-EOT
    ==============================
    Terralings System Status
    Environment: Local Learning
    Status: Operational
    ==============================
  # Missing proper EOT delimiter!
}
