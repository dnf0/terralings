# ==============================================================================
# Solution: gov02
# Chapter: 13_governance (Architecture Governance & Enterprise Standards)
# ==============================================================================

variable "service_name" {
  type        = string
  description = "Target service requiring secure IAM credentials"
  default     = "auth-service"
}

locals {
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
    policy_arns = {
      secret_read = local.secret_module.read_only_policy_arn
      kms_decrypt = local.kms_module.decrypt_policy_arn
    }
  }
}

output "attached_policies" {
  value = terraform_data.iam_role.output.policy_arns
}
