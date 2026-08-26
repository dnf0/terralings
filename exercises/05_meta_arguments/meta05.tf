# I AM NOT DONE
# ==============================================================================
# Exercise: meta05
# Chapter: 05_meta_arguments (Meta-Arguments & Resource Scaling)
#
# Task:
# Sometimes resources are modified outside of Terraform (e.g. autoscaling policies,
# external tagging systems, or dynamically updated metadata). The `ignore_changes`
# lifecycle argument instructs Terraform to ignore differences in specified attributes
# during plan and apply.
#
# Complete the configuration below:
# 1. In resource "terraform_data" "managed_resource", add a `lifecycle` block.
# 2. Add `ignore_changes = [triggers_replace]` inside the lifecycle block so external
#    drift on triggers_replace does not plan updates.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

resource "terraform_data" "managed_resource" {
  input = "static-cluster-name"

  triggers_replace = {
    external_revision = "rev-1"
  }

  # TODO: Configure lifecycle block with ignore_changes = [triggers_replace]
  # lifecycle {
  #   ignore_changes = [ ... ]
  # }
}

output "resource_name" {
  value = terraform_data.managed_resource.output
}
