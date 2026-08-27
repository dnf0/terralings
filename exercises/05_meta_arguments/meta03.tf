# ==============================================================================
# Exercise: meta03
# Chapter: 05_meta_arguments (Meta-Arguments & Resource Scaling)
#
# Task:
# Terraform automatically creates dependency graphs based on attribute references
# (implicit dependencies). When two resources have an ordering requirement that
# cannot be expressed via direct reference, use the `depends_on` meta-argument.
#
# Complete the configuration below:
# 1. Add `depends_on = [terraform_data.database_migration]` to resource "terraform_data" "web_api".
#
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

resource "terraform_data" "database_migration" {
  input = "db_schema_v2_applied"
}

resource "terraform_data" "web_api" {
  input = "web_api_v2_ready"

  # TODO: Declare explicit dependency on terraform_data.database_migration
  depends_on = [terraform_data.missing_migration]
}

output "status" {
  value = {
    db  = terraform_data.database_migration.output
    api = terraform_data.web_api.output
  }
}
