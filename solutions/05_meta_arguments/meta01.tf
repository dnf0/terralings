# ==============================================================================
# Solution: meta01
# Chapter: 05_meta_arguments (Meta-Arguments & Resource Scaling)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "replica_count" {
  type        = number
  description = "Number of worker instances to scale"
  default     = 3
}

resource "terraform_data" "worker" {
  count = var.replica_count
  input = format("worker-%d.internal", count.index + 1)
}

output "worker_hostnames" {
  value = terraform_data.worker[*].input
}
