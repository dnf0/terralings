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
  # TODO (What): Configure OpenTofu encryption block with pbkdf2 key provider, aes_gcm method, and enforced state encryption.
  # TODO (Why): OpenTofu state encryption guarantees client-side data protection for secrets in state files prior to remote backend transit.
  encryption {
    # TODO (What): Define key_provider "pbkdf2" "passphrase" with passphrase = "correct-horse-battery-staple".
    # TODO (Why): Key providers supply cryptographic material used to derive data encryption keys.
    # TODO (What): Define method "aes_gcm" "main" with keys = key_provider.pbkdf2.passphrase.
    # TODO (Why): Encryption methods define the cipher (AES-GCM) applied to state snapshots and plan files.

    state {
      # TODO (What): Set method = method.aes_gcm.main and enforced = true.
      # TODO (Why): Enforcing encryption ensures unencrypted fallback states are rejected during apply.
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
