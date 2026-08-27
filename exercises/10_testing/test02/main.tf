# ==============================================================================
# Exercise: test02
# Chapter: 10_testing (Native Unit & Integration Testing)
#
# Task:
# Tests can execute `command = apply` to test actual resource creation and output
# values. Tests run sequentially, allowing multi-stage lifecycle testing.
#
# In this exercise:
# 1. In `lifecycle.tftest.hcl`, configure `command = apply` in the run block.
# 2. Pass variable `cluster_tier = "premium"`.
# 3. Assert that `terraform_data.cluster.output.tier == "premium"` and
#    `terraform_data.cluster.output.nodes == 5`.
#
# ==============================================================================

variable "cluster_tier" {
  type    = string
  default = "standard"
}

locals {
  node_count = var.cluster_tier == "premium" ? 5 : 2
}

resource "terraform_data" "cluster" {
  input = {
    tier  = var.cluster_tier
    nodes = local.node_count
  }
}

output "cluster_info" {
  value = terraform_data.cluster.output
}
