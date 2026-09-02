# ==============================================================================
# Solution: aws06
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

resource "aws_s3_bucket" "lake" {
  bucket        = "corp-data-lake-primary-9901"
  force_destroy = false
}

resource "aws_s3_bucket_public_access_block" "guard" {
  bucket                  = aws_s3_bucket.lake.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_dynamodb_table" "users" {
  name         = "Users"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "UserId"

  attribute {
    name = "UserId"
    type = "S"
  }
}

output "bucket_arn" {
  value = aws_s3_bucket.lake.arn
}

output "table_arn" {
  value = aws_dynamodb_table.users.arn
}
