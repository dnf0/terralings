# ==============================================================================
# Solution: dynamic01
# Chapter: 06_dynamic_blocks (Dynamic Blocks & Advanced HCL)
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

variable "files" {
  type = map(string)
  default = {
    "index.html" = "<html><body>Hello Terralings</body></html>"
    "style.css"  = "body { background-color: #1a1a1a; color: white; }"
    "app.js"     = "console.log('Terralings loaded');"
  }
}

data "archive_file" "bundle" {
  type        = "zip"
  output_path = "${path.module}/bundle.zip"

  dynamic "source" {
    for_each = var.files
    content {
      filename = source.key
      content  = source.value
    }
  }
}

output "archive_sha" {
  value = data.archive_file.bundle.output_sha
}
