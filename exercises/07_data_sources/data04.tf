# ==============================================================================
# Exercise: data04
# Chapter: 07_data_sources (Data Sources & State Querying)
#
# Task:
# Terraform 1.2+ introduced custom conditions on resources and data sources:
# - `precondition`: Evaluated before the resource or data source is planned/evaluated.
#   Used to validate upstream inputs, variables, or environment constraints.
# - `postcondition`: Evaluated after the resource or data source is planned/evaluated.
#   Used to validate that output attributes meet required security/compliance criteria
#   (using `self.<attr>`).
#
# Complete the configuration below:
# 1. In resource "terraform_data" "cluster", add a `lifecycle` block.
# 2. Add a `precondition` ensuring `var.cluster_size >= 3` with an error message:
#    "Cluster size must be at least 3 for quorum."
# 3. Add a `postcondition` ensuring `self.output != null` with an error message:
#    "Cluster output must not be null."
#
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "cluster_size" {
  type        = number
  description = "Target number of cluster nodes"
  default     = 3
}

resource "terraform_data" "cluster" {
  input = {
    nodes = var.cluster_size
    role  = "consensus"
  }

  lifecycle {
    precondition {
      # TODO: Fix condition to var.cluster_size >= 3
      condition     = var.cluster_size >= 10
      error_message = "Cluster size must be at least 3 for quorum."
    }
    postcondition {
      condition     = self.output != null
      error_message = "Cluster output must not be null."
    }
  }
}

output "cluster_state" {
  value = terraform_data.cluster.output
}
