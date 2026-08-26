# ==============================================================================
# Solution: test02
# Chapter: 10_testing (Native Unit & Integration Testing)
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
