# ==============================================================================
# Exercise: gov02
# Chapter: 13_governance (Architecture Governance & Enterprise Standards)
#
# Task:
# Under Architecture Decision Record ADR-0005:
# "The module that owns the resource owns the policies that talk to it."
#
# Anti-Patterns:
#   - Writing inline IAM JSON policies (`aws_iam_role_policy`).
#   - Wildcard grants (`secret:*`, `kms:*`, or string ARN wildcards like `arn:aws:...:secret:foo*`).
#
# Correct Architectural Pattern:
#   - Resources (secrets, KMS keys) are instantiated via dedicated catalog modules.
#   - The resource module creates and exports least-privilege managed policy ARNs
#     (e.g., `read_only_policy_arn`, `decrypt_policy_arn`).
#   - Consuming IAM roles attach those managed policy ARNs directly via `policy_arns = map(string)`.
#
# In this exercise:
# 1. In `terraform_data.iam_role.input.policy_arns`, attach the managed policy ARNs
#    exported by the simulated secret and KMS modules:
#    - `secret_read` = local.secret_module.read_only_policy_arn
#    - `kms_decrypt` = local.kms_module.decrypt_policy_arn
# 2. Output `attached_policies` referencing the IAM role's attached policy map.
#
# ==============================================================================

variable "service_name" {
  type        = string
  description = "Target service requiring secure IAM credentials"
  default     = "auth-service"
}

locals {
  # Simulated outputs from resource modules adhering to ADR-0005
  secret_module = {
    secret_arn           = "arn:aws:secretsmanager:us-east-1:123456789012:secret:db-credentials-7x8a"
    read_only_policy_arn = "arn:aws:iam::123456789012:policy/secret-db-credentials-read"
  }

  kms_module = {
    key_arn            = "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
    decrypt_policy_arn = "arn:aws:iam::123456789012:policy/kms-app-key-decrypt"
  }
}

resource "terraform_data" "iam_role" {
  input = {
    role_name = "${var.service_name}-execution-role"
    # TODO (What): Set policy_arns = { secret_read = local.secret_module.read_only_policy_arn, kms_decrypt = local.kms_module.decrypt_policy_arn }.
    # TODO (Why): ADR-0005: The module that owns the resource owns the policies that talk to it. Attaching exported managed policy ARNs avoids inline wildcard permissions.
    policy_arns = {}
  }

  lifecycle {
    postcondition {
      condition     = try(self.input.policy_arns.secret_read, "") == local.secret_module.read_only_policy_arn && try(self.input.policy_arns.kms_decrypt, "") == local.kms_module.decrypt_policy_arn
      error_message = "IAM role must attach managed policy ARNs exported by the secret and KMS modules."
    }
  }
}

output "attached_policies" {
  # TODO (What): Reference terraform_data.iam_role.output.policy_arns.
  # TODO (Why): Exposes the attached managed policies map for security auditing and compliance verification.
  value = terraform_data.iam_role.output.policy_arns
}
