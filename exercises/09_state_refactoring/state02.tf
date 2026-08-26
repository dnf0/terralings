# I AM NOT DONE
# ==============================================================================
# Exercise: state02
# Chapter: 09_state_refactoring (State Management & Refactoring Surgery)
#
# Task:
# When refactoring resources from `count` (indexed `[0]`, `[1]`) to `for_each`
# (keyed `["primary"]`, `["replica"]`), `moved` blocks prevent destructive
# recreation:
#
#   moved {
#     from = terraform_data.db_cluster[0]
#     to   = terraform_data.db_cluster["primary"]
#   }
#   moved {
#     from = terraform_data.db_cluster[1]
#     to   = terraform_data.db_cluster["replica"]
#   }
#
# In this exercise:
# 1. Add `moved` blocks migrating index `0` to `"primary"` and index `1` to `"replica"`.
# 2. Convert the resource from `count = 2` to `for_each = toset(["primary", "replica"])`.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

# TODO: Add moved blocks for [0] -> ["primary"] and [1] -> ["replica"]

resource "terraform_data" "db_cluster" {
  # TODO: Change count = 2 to for_each = toset(["primary", "replica"])
  count = 2

  input = {
    role = "node-${count.index}"
  }
}

output "cluster_roles" {
  # TODO: Update output expression to extract roles from for_each map
  value = [for node in terraform_data.db_cluster : node.output.role]
}
