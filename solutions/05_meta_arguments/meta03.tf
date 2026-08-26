# ==============================================================================
# Solution: meta03
# Chapter: 05_meta_arguments (Meta-Arguments & Resource Scaling)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

resource "terraform_data" "database_migration" {
  input = "db_schema_v2_applied"
}

resource "terraform_data" "web_api" {
  input = "web_api_v2_ready"

  depends_on = [terraform_data.database_migration]
}

output "status" {
  value = {
    db  = terraform_data.database_migration.output
    api = terraform_data.web_api.output
  }
}
