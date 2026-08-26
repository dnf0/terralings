# ==============================================================================
# Solution: gov03
# Chapter: 13_governance (Architecture Governance & Enterprise Standards)
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
  ephemeral_workload = {
    namespace = "ephemeral-${var.job_id}"
    workers   = var.worker_nodes
    ttl_mins  = var.max_runtime_minutes
    isolated  = true
  }
}

resource "terraform_data" "ephemeral_cluster" {
  input = local.ephemeral_workload

  lifecycle {
    replace_triggered_by = [
      terraform_data.job_trigger
    ]
  }
}

output "cluster_namespace" {
  value = terraform_data.ephemeral_cluster.output.namespace
}
