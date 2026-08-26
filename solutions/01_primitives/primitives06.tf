# ==============================================================================
# Solution: primitives06
# Chapter: 01_primitives (HCL Foundations & Core Primitives)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

resource "terraform_data" "deploy_token" {
  input = "v1.0.0"

  triggers_replace = [
    "v1.0.0"
  ]
}
