# ==============================================================================
# Exercise: tofu02
# Chapter: 12_opentofu (OpenTofu Innovations & Enterprise Features)
#
# Task:
# OpenTofu evaluates input variables early in the dependency graph. This enables
# robust dynamic interpolation across configurations and unified variable evaluation.
#
# In this exercise:
# 1. Define local value `cluster_fqdn` by combining `var.subdomain`, `var.environment`,
#    and `var.domain_suffix` in the format `<subdomain>.<environment>.<domain_suffix>`.
# 2. Complete `terraform_data.gateway` resource inputs with `fqdn = local.cluster_fqdn`
#    and `port = var.port`.
# 3. Output `gateway_fqdn` referencing the planned FQDN.
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
  # TODO (What): Construct cluster_fqdn as "${var.subdomain}.${var.environment}.${var.domain_suffix}".
  # TODO (Why): Standardizing domain name strings in locals guarantees consistent internal addressing across workload gateways.
  cluster_fqdn = ""
}

resource "terraform_data" "gateway" {
  # TODO (What): Set fqdn = local.cluster_fqdn and port = var.port in the input map.
  # TODO (Why): Clean local interpolation passes validated endpoint strings directly into infrastructure inputs.
  input = {
    fqdn = split(".", local.cluster_fqdn)[2]
    port = var.port
  }
}

output "gateway_fqdn" {
  # TODO (What): Reference terraform_data.gateway.output.fqdn.
  # TODO (Why): Exposes the generated gateway address to consuming client services.
  value = terraform_data.gateway.output.fqdn
}
