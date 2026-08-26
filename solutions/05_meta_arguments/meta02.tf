# ==============================================================================
# Solution: meta02
# Chapter: 05_meta_arguments (Meta-Arguments & Resource Scaling)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

variable "microservices" {
  type = map(object({
    port = number
    tier = string
  }))
  default = {
    "auth" = {
      port = 8081
      tier = "backend"
    }
    "payment" = {
      port = 8082
      tier = "backend"
    }
    "web" = {
      port = 80
      tier = "frontend"
    }
  }
}

resource "terraform_data" "service" {
  for_each = var.microservices

  input = {
    name = each.key
    port = each.value.port
    tier = each.value.tier
  }
}

output "service_keys" {
  value = keys(terraform_data.service)
}
