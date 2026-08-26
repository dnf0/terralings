# ==============================================================================
# Solution: meta04
# Chapter: 05_meta_arguments (Meta-Arguments & Resource Scaling)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "app_version" {
  type    = string
  default = "v2.1.0"
}

resource "terraform_data" "zero_downtime_app" {
  input = var.app_version

  lifecycle {
    create_before_destroy = true
  }
}

output "deployed_version" {
  value = terraform_data.zero_downtime_app.output
}
