# ==============================================================================
# Solution: pattern04
# Chapter: 11_patterns (Production Patterns & Anti-Patterns)
# ==============================================================================

variable "services" {
  type = map(object({
    port    = number
    enabled = optional(bool, true)
    tier    = optional(string, "standard")
  }))
  description = "Self-service service registration contract"
  default = {
    auth = {
      port    = 8081
      enabled = true
      tier    = "premium"
    }
    legacy = {
      port    = 8080
      enabled = false
    }
    billing = {
      port    = 8082
      enabled = true
    }
  }
}

locals {
  active_services = {
    for name, svc in var.services : name => svc if svc.enabled
  }

  routing_table = {
    for name, svc in local.active_services :
    name => format("https://%s.internal:%d", name, svc.port)
  }
}

resource "terraform_data" "service_route" {
  for_each = local.active_services

  input = {
    name     = each.key
    port     = each.value.port
    endpoint = local.routing_table[each.key]
  }
}

output "routes" {
  value = { for k, v in terraform_data.service_route : k => v.output.endpoint }
}
