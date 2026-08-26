# ==============================================================================
# Solution: pattern02
# Chapter: 11_patterns (Production Patterns & Anti-Patterns)
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
  count = var.enable_disaster_recovery ? 1 : 0

  input = {
    target_region   = var.dr_region
    primary_cluster = terraform_data.primary_cluster.output.cluster_id
  }
}

output "primary_id" {
  value = terraform_data.primary_cluster.output.cluster_id
}

output "dr_vault_region" {
  value = try(terraform_data.dr_vault[0].output.target_region, "disabled")
}
