# ==============================================================================
# Solution: module03
# Chapter: 08_modules (Modular Infrastructure Architecture)
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
  source      = "./modules/worker"
  for_each    = var.workers
  worker_name = each.key
  concurrency = each.value
}

output "worker_ids" {
  value = { for k, w in module.workers : k => w.worker_id }
}
