# ==============================================================================
# Solution: module02
# Chapter: 08_modules (Modular Infrastructure Architecture)
# ==============================================================================

module "auth" {
  source       = "./modules/service"
  service_name = "auth-service"
  port         = 8080
}

output "endpoint" {
  value = module.auth.service_url
}
