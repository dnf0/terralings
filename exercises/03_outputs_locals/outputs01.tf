# ==============================================================================
# Exercise: outputs01
# Chapter: 03_outputs_locals (Outputs, Locals & Expressions)
#
# Task:
# Output values export information about resources to the CLI or parent modules.
# When an output contains sensitive information (like passwords or API keys),
# mark it with `sensitive = true` to prevent it from being printed in cleartext in logs and plans.
#
# Complete the two output blocks below:
# 1. `username`: exports terraform_data.credentials.input.username (non-sensitive)
# 2. `api_key`: exports terraform_data.credentials.input.api_key with sensitive = true
#
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

resource "terraform_data" "credentials" {
  input = {
    username = "admin"
    api_key  = "secret-api-token-998877"
  }
}

output "username" {
  description = "Public admin username"
  # TODO: Set value = terraform_data.credentials.input.username
}

output "api_key" {
  description = "Sensitive API access key"
  # TODO: Set value = terraform_data.credentials.input.api_key and sensitive = true
}
