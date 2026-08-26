# ==============================================================================
# Solution: variables01
# Chapter: 02_variables (Input Variables, Types & Validations)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "environment" {
  type        = string
  description = "Deployment target environment"
  default     = "development"
}

variable "port" {
  type        = number
  description = "Application server listening port"
  default     = 8080
}

variable "debug_mode" {
  type        = bool
  description = "Flag to toggle verbose logging"
  default     = false
}
