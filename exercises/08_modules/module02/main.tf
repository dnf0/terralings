# ==============================================================================
# Exercise: module02
# Chapter: 08_modules (Modular Infrastructure Architecture)
#
# Task:
# When calling a local child module, you specify the relative filesystem path in
# the `source` argument. The root module can consume outputs from the child
# module using `module.<MODULE_NAME>.<OUTPUT_NAME>`.
#
# In this exercise:
# 1. Call the child module at `./modules/service` with:
#    - `service_name = "auth-service"`
#    - `port         = 8080`
# 2. In the root `output "endpoint"`, expose the child module's `service_url` output.
#
# ==============================================================================

# TODO: Complete the module call and root output
module "auth" {
  source = "./modules/service"

  # service_name = ...
  # port         = ...
}

output "endpoint" {
  # value = module.auth.service_url
  value = ""
}
