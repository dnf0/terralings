# ==============================================================================
# Solution: tofu02
# Chapter: 12_opentofu (OpenTofu Innovations & Enterprise Features)
# ==============================================================================

variable "environment" {
  type        = string
  description = "Target deployment environment"
  default     = "staging"
}

variable "subdomain" {
  type        = string
  description = "Application subdomain"
  default     = "api"
}

variable "domain_suffix" {
  type        = string
  description = "Root domain suffix"
  default     = "internal.cloud"
}

variable "port" {
  type        = number
  description = "Gateway port"
  default     = 8443
}

locals {
  cluster_fqdn = "${var.subdomain}.${var.environment}.${var.domain_suffix}"
}

resource "terraform_data" "gateway" {
  input = {
    fqdn = local.cluster_fqdn
    port = var.port
  }
}

output "gateway_fqdn" {
  value = terraform_data.gateway.output.fqdn
}
