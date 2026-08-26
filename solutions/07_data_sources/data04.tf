# ==============================================================================
# Solution: data04
# Chapter: 07_data_sources (Data Sources & State Querying)
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
      condition     = var.cluster_size >= 3
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
