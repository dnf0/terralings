# ==============================================================================
# Solution: module05
# Chapter: 08_modules (Modular Infrastructure Architecture)
# ==============================================================================

module "database" {
  source  = "./modules/database"
  db_name = "orders"
}

module "app" {
  source   = "./modules/app"
  app_name = "orders-api"
  db_host  = module.database.endpoint
}

output "app_connection" {
  value = module.app.connection_string
}
