# I AM NOT DONE
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
# When done, remove the '# I AM NOT DONE' line at the top.
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
  # TODO: Set for_each = var.microservices

  # TODO: Set input using each.key and each.value
  input = {}
}

output "service_keys" {
  # TODO: Output keys(terraform_data.service)
  value = []
}
