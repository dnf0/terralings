output "service_url" {
  description = "Internal service URL"
  value       = "http://${var.service_name}:${var.port}"
}
