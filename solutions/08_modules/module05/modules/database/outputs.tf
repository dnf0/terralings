output "endpoint" {
  description = "Database endpoint host and port"
  value       = terraform_data.db.output.endpoint
}
