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
# ==============================================================================

# TODO (What): Add moved blocks mapping terraform_data.db_cluster[0] -> terraform_data.db_cluster["primary"] and terraform_data.db_cluster[1] -> terraform_data.db_cluster["replica"].
# TODO (Why): Moving indexed count items to semantic string keys prevents database teardown and preserves existing provisioned instances.

resource "terraform_data" "db_cluster" {
  # TODO (What): Change count = 2 to for_each = toset(["primary", "replica"]).
  # TODO (Why): for_each indexes resources by domain-meaningful string keys ("primary", "replica") rather than fragile numerical positions.
  count = 2

  input = {
    # TODO (What): Reference each.key for the role attribute.
    # TODO (Why): Inside a for_each iteration, each.key supplies the active set element string value.
    role = each.key
  }
}

output "cluster_roles" {
  # TODO (What): Extract roles using [for node in terraform_data.db_cluster : node.output.role].
  # TODO (Why): Comprehensions iterate over the map of resource instances to construct a clean list output.
  value = [for node in terraform_data.db_cluster : node.output.role]
}
