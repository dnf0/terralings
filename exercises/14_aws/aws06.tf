# ==============================================================================
# Exercise: aws06
# Chapter: 14_aws (AWS Infrastructure & Production Blueprints)
#
# Task:
# Production data architectures require encrypted S3 buckets with Public Access Blocks
# and DynamoDB on-demand tables with Global Secondary Indexes:
#
#   resource "aws_s3_bucket_public_access_block" "guard" {
#     bucket                  = aws_s3_bucket.lake.id
#     block_public_acls       = true
#     block_public_policy     = true
#     ignore_public_acls      = true
#     restrict_public_buckets = true
#   }
#
# In this exercise:
# 1. Create `aws_s3_bucket.lake` with `bucket = "corp-data-lake-primary-9901"` and `force_destroy = false`.
# 2. Attach `aws_s3_bucket_public_access_block.guard` setting all four blocking arguments to `true`.
# 3. Create `aws_dynamodb_table.users` with `name = "Users"`, `billing_mode = "PAY_PER_REQUEST"`,
#    `hash_key = "UserId"`, and attribute definition for `UserId` of type `"S"`.
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
  # TODO (What): Set force_destroy = false.
  # TODO (Why): Production buckets protect business data by preventing accidental deletion when non-empty.
  force_destroy = true
}

resource "aws_s3_bucket_public_access_block" "guard" {
  bucket = aws_s3_bucket.lake.id
  # TODO (What): Set block_public_acls = true, block_public_policy = true, ignore_public_acls = true, restrict_public_buckets = true.
  # TODO (Why): The S3 Public Access Block prevents accidental data leakage across all bucket policies and object ACLs.
  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_dynamodb_table" "users" {
  name         = "Users"
  # TODO (What): Set billing_mode = "PAY_PER_REQUEST" and hash_key = "UserId".
  # TODO (Why): On-demand billing (PAY_PER_REQUEST) automatically scales capacity with zero manual capacity planning.
  billing_mode = "PROVISIONED"
  hash_key     = ""

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
