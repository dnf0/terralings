# ==============================================================================
# Exercise: pattern02
# Chapter: 11_patterns (Production Patterns & Anti-Patterns)
#
# Task:
# Feature flags allow turning optional infrastructure components on or off
# dynamically using boolean input variables and `count = var.enabled ? 1 : 0`.
#
# When referencing conditionally created resources in outputs, use safe accessor
# patterns such as `try(resource[0].id, null)` or `one(resource[*].id)`:
#
#   resource "terraform_data" "dr_vault" {
#     count = var.enable_disaster_recovery ? 1 : 0
#     input = { ... }
#   }
#
#   output "vault_id" {
#     value = try(terraform_data.dr_vault[0].output.region, "none")
#   }
#
# In this exercise:
# 1. Update `terraform_data.dr_vault` to set `count = var.enable_disaster_recovery ? 1 : 0`.
# 2. In `terraform_data.dr_vault.input`, pass `target_region = var.dr_region` and
#    `primary_cluster = terraform_data.primary_cluster.output.cluster_id`.
# 3. Update output `dr_vault_region` using `try()` or `one()` to return the DR
#    vault's target region, or `"disabled"` when the flag is false.
#
# ==============================================================================

variable "enable_disaster_recovery" {
  type        = bool
  description = "Feature flag enabling asynchronous cross-region disaster recovery replication"
  default     = true
}

variable "dr_region" {
  type        = string
  description = "Target disaster recovery region"
  default     = "us-west-2"
}

resource "terraform_data" "primary_cluster" {
  input = {
    cluster_id = "cluster-main-01"
    region     = "us-east-1"
  }
}

resource "terraform_data" "dr_vault" {
  # TODO: Use count with var.enable_disaster_recovery to conditionally provision this resource
  count = 0

  # TODO: Set target_region and primary_cluster inputs
  input = {}
}

output "primary_id" {
  value = terraform_data.primary_cluster.output.cluster_id
}

output "dr_vault_region" {
  # TODO: Safely access dr_vault target_region, fallback to "disabled" if count is 0
  value = terraform_data.dr_vault[0].output.target_region
}
