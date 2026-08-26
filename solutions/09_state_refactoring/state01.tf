# ==============================================================================
# Solution: state01
# Chapter: 09_state_refactoring (State Management & Refactoring Surgery)
# ==============================================================================

moved {
  from = terraform_data.legacy_cache
  to   = terraform_data.redis_cache
}

resource "terraform_data" "redis_cache" {
  input = {
    host = "redis-cluster.internal"
    port = 6379
  }
}

output "cache_host" {
  value = terraform_data.redis_cache.output.host
}
