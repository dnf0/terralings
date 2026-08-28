# ==============================================================================
# Exercise: meta02
# Chapter: 05_meta_arguments (Meta-Arguments & Resource Scaling)
#
# Task:
# While `count` indexes resources by integer position (0, 1, 2), `for_each`
# creates resources keyed by strings (from maps or sets). This prevents
# unwanted cascade destructions when removing an element from the middle of a list.
#
# Complete the configuration below:
# 1. In resource "terraform_data" "service", add `for_each = var.microservices`.
# 2. Inside the resource block, use `each.key` and `each.value` to construct:
#    input = {
#      name = each.key
#      port = each.value.port
#      tier = each.value.tier
#    }
# 3. In the output "service_keys", return `keys(terraform_data.service)`.
#
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
  # TODO (What): Add for_each = var.microservices.
  # TODO (Why): for_each identifies instances by stable string keys (e.g. "auth", "payment"), preventing instance recreation when array items are reordered.

  # TODO (What): Populate input using each.key for the name, and each.value.port and each.value.tier for attributes.
  # TODO (Why): Inside a for_each block, each.key references the map key and each.value references the corresponding map element.
  input = {
    name = each.key
    port = each.value.port
    tier = each.value.tier
  }
}

output "service_keys" {
  # TODO (What): Set value = keys(terraform_data.service).
  # TODO (Why): Resources declared with for_each evaluate to a map of instances keyed by the for_each map keys.
  value = keys(terraform_data.service)
}
