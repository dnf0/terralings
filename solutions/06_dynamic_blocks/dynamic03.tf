# ==============================================================================
# Solution: dynamic03
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

variable "service_configs" {
  type = map(list(object({
    name    = string
    content = string
  })))
  default = {
    "auth" = [
      { name = "config.json", content = "{\"service\": \"auth\"}" },
      { name = "policy.hcl", content = "path \"secret/*\" { capabilities = [\"read\"] }" }
    ]
    "billing" = [
      { name = "config.json", content = "{\"service\": \"billing\"}" }
    ]
  }
}

locals {
  all_sources = flatten([
    for s_name, files in var.service_configs : [
      for f in files : {
        path = "${s_name}/${f.name}"
        body = f.content
      }
    ]
  ])
}

data "archive_file" "multi_service" {
  type        = "zip"
  output_path = "${path.module}/services.zip"

  dynamic "source" {
    for_each = local.all_sources
    content {
      filename = source.value.path
      content  = source.value.body
    }
  }
}

output "package_sha" {
  value = data.archive_file.multi_service.output_sha
}
