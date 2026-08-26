# ==============================================================================
# Solution: module01
# Chapter: 08_modules (Modular Infrastructure Architecture)
# ==============================================================================

module "app_config" {
  source   = "./modules/app_config"
  app_name = "terralings"
}
