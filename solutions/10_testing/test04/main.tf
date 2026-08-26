# ==============================================================================
# Solution: test04
# Chapter: 10_testing (Native Unit & Integration Testing)
# ==============================================================================

variable "port" {
  type        = number
  description = "Port number between 1 and 65535"

  validation {
    condition     = var.port >= 1 && var.port <= 65535
    error_message = "port must be between 1 and 65535"
  }
}

output "port_number" {
  value = var.port
}
