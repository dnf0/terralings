# ==============================================================================
# Solution: data02
# Chapter: 07_data_sources (Data Sources & State Querying)
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

data "archive_file" "package" {
  type        = "zip"
  output_path = "${path.module}/package.zip"

  source {
    content  = "console.log('App ready');"
    filename = "index.js"
  }
}

output "archive_info" {
  value = {
    sha  = data.archive_file.package.output_sha
    size = data.archive_file.package.output_size
  }
}
