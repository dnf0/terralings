# I AM NOT DONE
# ==============================================================================
# Exercise: dynamic03
# Chapter: 06_dynamic_blocks (Dynamic Blocks & Advanced HCL)
#
# Task:
# When working with hierarchical or multi-layer configurations (such as a map of
# services each containing multiple configuration files), you often combine
# `flatten()` and `for` expressions to feed a `dynamic` block.
#
# Complete the configuration below:
# 1. In `locals`, compute `all_sources` by flattening the nested `var.service_configs`:
#    For each service (key `s_name`) and each file (`f`), create an object with:
#    - `path`: "${s_name}/${f.name}"
#    - `body`: f.content
# 2. In `data "archive_file" "multi_service"`, iterate over `local.all_sources`
#    using `dynamic "source"`, mapping `filename = source.value.path` and
#    `content = source.value.body`.
#
# When done, remove the '# I AM NOT DONE' line at the top.
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
  # TODO: Flatten nested var.service_configs into a 1D list of { path = "...", body = "..." }
  all_sources = []
}

data "archive_file" "multi_service" {
  type        = "zip"
  output_path = "${path.module}/services.zip"

  # TODO: Declare dynamic "source" over local.all_sources
  # dynamic "source" {
  #   for_each = local.all_sources
  #   content {
  #     filename = ...
  #     content  = ...
  #   }
  # }
}

output "package_sha" {
  value = data.archive_file.multi_service.output_sha
}
