# I AM NOT DONE
# ==============================================================================
# Exercise: module03
# Chapter: 08_modules (Modular Infrastructure Architecture)
#
# Task:
# OpenTofu and Terraform support `for_each` and `count` on `module` blocks to
# instantiate multiple distinct instances of a child module.
#
# In this exercise:
# 1. Instantiate `./modules/worker` using `for_each` over `var.workers` map.
# 2. Pass `worker_name = each.key` and `concurrency = each.value` to each module.
# 3. Expose `output "worker_ids"` as a map of `{ worker_name => module.workers[worker_name].worker_id }`
#    using a `for` expression: `{ for k, w in module.workers : k => w.worker_id }`.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

variable "workers" {
  type = map(number)
  default = {
    "crawler"   = 4
    "indexer"   = 8
    "processor" = 2
  }
}

module "workers" {
  source = "./modules/worker"

  # TODO: Add for_each and pass arguments
  # for_each    = var.workers
  # worker_name = each.key
  # concurrency = each.value
}

output "worker_ids" {
  # TODO: Build map of worker_id outputs
  # value = { for k, w in module.workers : k => w.worker_id }
  value = {}
}
