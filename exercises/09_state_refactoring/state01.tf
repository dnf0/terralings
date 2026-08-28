# ==============================================================================
# Exercise: state01
# Chapter: 09_state_refactoring (State Management & Refactoring Surgery)
#
# Task:
# When you rename a resource block in OpenTofu or Terraform, the engine would
# naturally destroy the old resource address and create a brand-new one, causing
# data loss or downtime.
#
# The `moved` block allows you to record refactoring operations declaratively in
# code without manually running CLI state commands:
#
#   moved {
#     from = terraform_data.legacy_cache
#     to   = terraform_data.redis_cache
#   }
#
# In this exercise:
# 1. Add a `moved` block recording that `terraform_data.legacy_cache` moved to
#    `terraform_data.redis_cache`.
# 2. Rename the resource block to `terraform_data.redis_cache`.
#
# ==============================================================================

moved {
  # TODO (What): Update 'to' address to terraform_data.redis_cache.
  # TODO (Why): Declarative moved blocks migrate state addresses automatically during plan/apply, preventing accidental destruction of live infrastructure.
  from = terraform_data.legacy_cache
  to   = terraform_data.missing_cache
}

# TODO (What): Rename this resource from legacy_cache to redis_cache.
# TODO (Why): Renaming the resource definition completes the refactoring target matched by the moved block.
resource "terraform_data" "legacy_cache" {
  input = {
    host = "redis-cluster.internal"
    port = 6379
  }
}

output "cache_host" {
  # TODO (What): Update reference to terraform_data.redis_cache.output.host.
  # TODO (Why): Downstream references must target the updated resource identifier.
  value = terraform_data.legacy_cache.output.host
}
