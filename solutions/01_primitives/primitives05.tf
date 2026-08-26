# ==============================================================================
# Solution: primitives05
# Chapter: 01_primitives (HCL Foundations & Core Primitives)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

resource "terraform_data" "formatted" {
  input = "clean canonical hcl"
}
