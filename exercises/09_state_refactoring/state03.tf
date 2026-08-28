# ==============================================================================
# Exercise: state03
# Chapter: 09_state_refactoring (State Management & Refactoring Surgery)
#
# Task:
# Declarative `import` blocks allow bringing existing infrastructure into
# state management as part of the code, reviewed via plans:
#
#   import {
#     to = terraform_data.existing_storage
#     id = "corp-storage-vault-01"
#   }
#
# In this exercise:
# 1. Write an `import` block targeting `terraform_data.existing_storage` with
#    `id = "corp-storage-vault-01"`.
# 2. Define the matching `resource "terraform_data" "existing_storage"` block.
#
# ==============================================================================

# TODO (What): Add an import block targeting 'to = terraform_data.existing_storage' with 'id = "corp-storage-vault-01"'.
# TODO (Why): Declarative import blocks allow engineers to review import operations via standard pull requests and plan diffs.

# TODO (What): Define resource "terraform_data" "existing_storage" with input = "corp-storage-vault-01".
# TODO (Why): The import block requires a corresponding resource configuration block in HCL to receive the imported remote resource.
# resource "terraform_data" "existing_storage" {
#   input = "corp-storage-vault-01"
# }

output "imported_id" {
  # value = terraform_data.existing_storage.input
  value = terraform_data.existing_storage.input
}
