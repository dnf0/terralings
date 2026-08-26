# ==============================================================================
# Solution: meta05
# Chapter: 05_meta_arguments (Meta-Arguments & Resource Scaling)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

resource "terraform_data" "managed_resource" {
  input = "static-cluster-name"

  triggers_replace = {
    external_revision = "rev-1"
  }

  lifecycle {
    ignore_changes = [triggers_replace]
  }
}

output "resource_name" {
  value = terraform_data.managed_resource.output
}
