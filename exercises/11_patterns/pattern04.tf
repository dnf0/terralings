# ==============================================================================
# Exercise: pattern04
# Chapter: 11_patterns (Production Patterns & Anti-Patterns)
#
# Task:
# Platform engineering teams provide self-service infrastructure contracts by
# defining typed input variables with validations and using `for` comprehensions
# to project user requests into concrete resources:
#
#   locals {
#     active_services = {
#       for name, svc in var.services : name => svc if svc.enabled
#     }
#     routing_table = {
#       for name, svc in local.active_services :
#       name => format("https://%s.internal:%d", name, svc.port)
#     }
#   }
#
# In this exercise:
# 1. In `locals`, construct `active_services` filtering only services where `svc.enabled == true`.
# 2. In `locals`, construct `routing_table` mapping service names to `https://<name>.internal:<port>`.
# 3. Create `terraform_data.service_route` iterating with `for_each = local.active_services`
#    passing `name`, `port`, and `endpoint = local.routing_table[each.key]`.
# 4. Output `routes` collecting all registered service endpoints.
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
  # TODO: Filter var.services for enabled == true
  active_services = {
    for name, svc in var.services : name => svc if svc.enabled
  }

  # TODO: Construct map of service name to formatted endpoint "https://<name>.internal:<port>"
  routing_table = {}
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
  # TODO: Map registered service names to their planned endpoints
  value = { for k, v in terraform_data.service_route : k => v.output.endpoint }
}
