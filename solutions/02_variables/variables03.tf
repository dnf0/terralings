# ==============================================================================
# Solution: variables03
# Chapter: 02_variables (Input Variables, Types & Validations)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "database_config" {
  type = object({
    host     = string
    port     = number
    username = string
    ssl      = optional(bool, true)
  })
  description = "Database connection configuration"
  default = {
    host     = "localhost"
    port     = 5432
    username = "postgres"
  }
}

variable "endpoint_pair" {
  type        = tuple([string, number])
  description = "Tuple of hostname and port"
  default     = ["api.internal", 443]
}
