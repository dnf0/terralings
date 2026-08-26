# ==============================================================================
# Solution: tofu01
# Chapter: 12_opentofu (OpenTofu Innovations & Enterprise Features)
# ==============================================================================

terraform {
  encryption {
    key_provider "pbkdf2" "passphrase" {
      passphrase = "correct-horse-battery-staple"
    }

    method "aes_gcm" "main" {
      keys = key_provider.pbkdf2.passphrase
    }

    state {
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
