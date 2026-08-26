# ==============================================================================
# Solution: state02
# Chapter: 09_state_refactoring (State Management & Refactoring Surgery)
# ==============================================================================

moved {
  from = terraform_data.db_cluster[0]
  to   = terraform_data.db_cluster["primary"]
}

moved {
  from = terraform_data.db_cluster[1]
  to   = terraform_data.db_cluster["replica"]
}

resource "terraform_data" "db_cluster" {
  for_each = toset(["primary", "replica"])

  input = {
    role = "node-${each.key}"
  }
}

output "cluster_roles" {
  value = [for k, node in terraform_data.db_cluster : node.output.role]
}
