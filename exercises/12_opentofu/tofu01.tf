# ==============================================================================
# Exercise: tofu01
# Chapter: 12_opentofu (OpenTofu Innovations & Enterprise Features)
#
# Task:
# OpenTofu 1.7+ introduces native end-to-end client-side state and plan encryption
# at rest. You configure encryption within the top-level `terraform` block:
#
#   terraform {
#     encryption {
#       key_provider "pbkdf2" "passphrase" {
#         passphrase = "my-secure-passphrase"
#       }
#
#       method "aes_gcm" "main" {
#         keys = key_provider.pbkdf2.passphrase
#       }
#
#       state {
#         method   = method.aes_gcm.main
#         enforced = true
#       }
#     }
#   }
#
# In this exercise:
# 1. Complete the `encryption` block inside `terraform {}`:
#    - Define `key_provider "pbkdf2" "passphrase"` with a passphrase.
#    - Define `method "aes_gcm" "main"` referencing `key_provider.pbkdf2.passphrase`.
#    - Define `state` block with `method = method.aes_gcm.main` and `enforced = true`.
# ==============================================================================

terraform {
  # TODO: Configure OpenTofu encryption block
  encryption {
    # TODO: Define key_provider "pbkdf2" "passphrase" with passphrase = "correct-horse-battery-staple"
    # TODO: Define method "aes_gcm" "main" referencing key_provider.pbkdf2.passphrase

    state {
      # TODO: Set method and enforced
      method   = method.aes_gcm.main
      enforced = true
    }
  }
}

resource "terraform_data" "encrypted_payload" {
  input = {
    status = "encrypted_at_rest"
  }
}
