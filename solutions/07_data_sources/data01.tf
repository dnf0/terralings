# ==============================================================================
# Solution: data01
# Chapter: 07_data_sources (Data Sources & State Querying)
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

data "local_file" "manifest" {
  filename = "${path.module}/data01.tf"
}

output "manifest_content_length" {
  value = length(data.local_file.manifest.content)
}
