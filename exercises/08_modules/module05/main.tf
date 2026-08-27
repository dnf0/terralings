# ==============================================================================
# Exercise: module05
# Chapter: 08_modules (Modular Infrastructure Architecture)
#
# Task:
# In clean infrastructure architecture, sibling modules remain isolated and
# unaware of each other. The root module coordinates communication by passing
# outputs from one module as inputs to another.
#
# In this exercise:
# 1. Instantiate `module "database"` from `./modules/database` with `db_name = "orders"`.
# 2. Instantiate `module "app"` from `./modules/app` with:
#    - `app_name = "orders-api"`
#    - `db_host  = module.database.endpoint`
# 3. Output `app_connection` displaying the application's connection string.
#
# ==============================================================================

# TODO: Coordinate module.database and module.app
module "database" {
  source = "./modules/database"
  # db_name = ...
}

module "app" {
  source = "./modules/app"
  # app_name = ...
  # db_host  = ...
}

output "app_connection" {
  # value = module.app.connection_string
  value = ""
}
