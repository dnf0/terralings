output "connection_string" {
  description = "Full application database connection string"
  value       = terraform_data.app.output.connect
}
