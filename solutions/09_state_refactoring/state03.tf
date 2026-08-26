# ==============================================================================
# Solution: state03
# Chapter: 09_state_refactoring (State Management & Refactoring Surgery)
# ==============================================================================

import {
  to = terraform_data.existing_storage
  id = "corp-storage-vault-01"
}

resource "terraform_data" "existing_storage" {
  input = "corp-storage-vault-01"
}

output "imported_id" {
  value = terraform_data.existing_storage.input
}
