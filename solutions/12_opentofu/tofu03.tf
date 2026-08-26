# ==============================================================================
# Solution: tofu03
# Chapter: 12_opentofu (OpenTofu Innovations & Enterprise Features)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"

  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
  }
}

resource "terraform_data" "registry_status" {
  input = {
    status = "verified"
  }
}

output "engine_status" {
  value = terraform_data.registry_status.output
}
