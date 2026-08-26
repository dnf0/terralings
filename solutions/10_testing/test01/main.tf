# ==============================================================================
# Solution: test01
# Chapter: 10_testing (Native Unit & Integration Testing)
# ==============================================================================

variable "environment" {
  type    = string
  default = "staging"
}

variable "service_name" {
  type    = string
  default = "order-service"
}

output "service_id" {
  value = "${var.environment}-${var.service_name}"
}
