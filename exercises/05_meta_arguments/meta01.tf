# I AM NOT DONE
# ==============================================================================
# Exercise: meta01
# Chapter: 05_meta_arguments (Meta-Arguments & Resource Scaling)
#
# Task:
# The `count` meta-argument creates multiple instances of a resource or module.
# Inside a block where `count` is set, the `count.index` object is available,
# representing the zero-based index of each instance.
#
# Complete the configuration below:
# 1. In resource "terraform_data" "worker", add `count = var.replica_count`.
# 2. Set input to `format("worker-%d.internal", count.index + 1)`.
# 3. In the output "worker_hostnames", reference all worker inputs using splat [*].
#
# When done, remove the '# I AM NOT DONE' line at the top.
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
  # TODO: Set count = var.replica_count

  # TODO: Set input = format("worker-%d.internal", count.index + 1)
  input = "worker"
}

output "worker_hostnames" {
  # TODO: Extract the inputs of all worker instances using splat expression: terraform_data.worker[*].input
  value = []
}
