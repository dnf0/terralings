# ==============================================================================
# Solution: outputs01
# Chapter: 03_outputs_locals (Outputs, Locals & Expressions)
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
  value       = terraform_data.credentials.input.username
}

output "api_key" {
  description = "Sensitive API access key"
  value       = terraform_data.credentials.input.api_key
  sensitive   = true
}
