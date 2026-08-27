# ==============================================================================
# Exercise: gov03
# Chapter: 13_governance (Architecture Governance & Enterprise Standards)
#
# Task:
# Ephemeral compute workloads and tooling (e.g., Ray cluster provisioners, ephemeral
# model training nodes, CI runner fleets) must be encapsulated in dedicated modules
# to namespace resources and prevent root module pollution.
#
# Key Requirements for Ephemeral Workloads:
#   - Namespaced resource identifiers preventing collision with persistent infra.
#   - Explicit auto-termination TTL / max runtime constraints.
#   - Replacement triggers bound to job/trigger runs (`replace_triggered_by`).
#
# In this exercise:
# 1. In `locals`, define `ephemeral_workload` containing:
#    - `namespace` = "ephemeral-${var.job_id}"
#    - `workers`   = var.worker_nodes
#    - `ttl_mins`  = var.max_runtime_minutes
#    - `isolated`  = true
# 2. Complete `terraform_data.ephemeral_cluster` passing `local.ephemeral_workload`
#    and bind its lifecycle with `replace_triggered_by = [terraform_data.job_trigger]`.
# 3. Output `cluster_namespace` referencing the namespace from the resource output.
#
# ==============================================================================

variable "job_id" {
  type        = string
  description = "Unique identifier for the ephemeral compute run"
  default     = "train-job-404"
}

variable "worker_nodes" {
  type        = number
  description = "Number of ephemeral compute workers"
  default     = 8
}

variable "max_runtime_minutes" {
  type        = number
  description = "Maximum auto-termination threshold in minutes"
  default     = 120
}

resource "terraform_data" "job_trigger" {
  input = var.job_id
}

locals {
  # TODO: Define ephemeral_workload namespace, workers, ttl_mins, and isolated flag
  ephemeral_workload = {}
}

resource "terraform_data" "ephemeral_cluster" {
  # TODO: Set input to local.ephemeral_workload
  input = local.ephemeral_workload

  lifecycle {
    replace_triggered_by = [terraform_data.job_trigger]
    postcondition {
      condition     = lookup(self.output, "namespace", "") == "ephemeral-${var.job_id}" && lookup(self.output, "workers", 0) == var.worker_nodes && lookup(self.output, "isolated", false) == true
      error_message = "Ephemeral compute workload must encapsulate namespace, workers, ttl_mins, and isolated flag."
    }
  }
}

output "cluster_namespace" {
  # TODO: Reference namespace from terraform_data.ephemeral_cluster.output
  value = terraform_data.ephemeral_cluster.output.namespace
}
