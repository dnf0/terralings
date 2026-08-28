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
  # TODO (What): Add count = var.replica_count to create multiple resource instances.
  # TODO (Why): The count meta-argument replicates resources based on an integer value, exposing count.index inside each instance.

  # TODO (What): Set input = format("worker-%d.internal", count.index + 1).
  # TODO (Why): Using count.index enables distinct instance naming and numbering in sequential clusters.
  input = format("worker-%d.internal", count.index + 1)
}

output "worker_hostnames" {
  # TODO (What): Set value = terraform_data.worker[*].input.
  # TODO (Why): When count is used on a resource, referencing it requires splat [*] or indexing [i] because the resource address resolves to a list of instances.
  value = terraform_data.worker[*].input
}
