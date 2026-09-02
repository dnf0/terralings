# ==============================================================================
# Exercise: aws05
# Chapter: 14_aws (AWS Infrastructure & Production Blueprints)
#
# Task:
# Zero-trust cloud security requires least-privilege IAM policies scoped to exact
# resource ARNs and strict security groups with stateful egress filtering:
#
#   resource "aws_security_group" "web_sg" {
#     vpc_id = "vpc-12345678"
#     ingress {
#       from_port   = 443
#       to_port     = 443
#       protocol    = "tcp"
#       cidr_blocks = ["10.0.0.0/16"]
#     }
#   }
#
# In this exercise:
# 1. Create `aws_iam_policy.s3_reader` granting `s3:GetObject` strictly to bucket ARN `arn:aws:s3:::corp-audit-vault/*`.
# 2. Attach policy to `aws_iam_role.app_role` using `aws_iam_role_policy_attachment.s3_attach`.
# 3. Configure `aws_security_group.app_sg` with HTTPS ingress on port 443 from CIDR `10.0.0.0/16`.
# ==============================================================================

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

resource "aws_iam_role" "app_role" {
  name = "enterprise-app-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
  })
}

resource "aws_iam_policy" "s3_reader" {
  name        = "AuditVaultReaderPolicy"
  description = "Scoped read-only access to audit bucket"

  # TODO (What): Define policy document granting Action = ["s3:GetObject"] on Resource = "arn:aws:s3:::corp-audit-vault/*".
  # TODO (Why): Principle of least privilege prohibits wildcard (*) actions on global resources.
  policy = ""
}

resource "aws_iam_role_policy_attachment" "s3_attach" {
  # TODO (What): Attach role = aws_iam_role.app_role.name and policy_arn = aws_iam_policy.s3_reader.arn.
  # TODO (Why): Policy attachments bind permissions declaratively to compute identities.
  role       = ""
  policy_arn = ""
}

resource "aws_security_group" "app_sg" {
  name        = "app-ingress-sg"
  description = "Filtered ingress security group"
  vpc_id      = "vpc-12345678"

  ingress {
    description = "Allow HTTPS internal traffic"
    # TODO (What): Set from_port = 443, to_port = 443, protocol = "tcp", cidr_blocks = ["10.0.0.0/16"].
    # TODO (Why): Restricting ingress to internal VPC CIDRs protects compute tiers from public probing.
    from_port   = 0
    to_port     = 0
    protocol    = ""
    cidr_blocks = []
  }
}

output "role_arn" {
  value = aws_iam_role.app_role.arn
}
