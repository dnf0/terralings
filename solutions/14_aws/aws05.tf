# ==============================================================================
# Solution: aws05
# Chapter: 14_aws (AWS Infrastructure & Production Blueprints)
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

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["s3:GetObject"]
      Resource = "arn:aws:s3:::corp-audit-vault/*"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "s3_attach" {
  role       = aws_iam_role.app_role.name
  policy_arn = aws_iam_policy.s3_reader.arn
}

resource "aws_security_group" "app_sg" {
  name        = "app-ingress-sg"
  description = "Filtered ingress security group"
  vpc_id      = "vpc-12345678"

  ingress {
    description = "Allow HTTPS internal traffic"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }
}

output "role_arn" {
  value = aws_iam_role.app_role.arn
}
