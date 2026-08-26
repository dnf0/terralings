# I AM NOT DONE
# ==============================================================================
# Exercise: module04
# Chapter: 08_modules (Modular Infrastructure Architecture)
#
# Task:
# By default, child modules inherit default (un-aliased) provider configurations
# from the caller. When a child module needs an aliased provider (e.g. for a
# secondary region, account, or custom config), you must:
# 1. Declare `configuration_aliases = [local.target]` in the child module's
#    `terraform.required_providers.local` block.
# 2. Pass `providers = { local.target = local.backup }` in the parent module block.
#
# In this exercise:
# 1. Update `modules/storage/main.tf` to declare `configuration_aliases = [local.target]`.
# 2. In this root `main.tf`, pass `providers = { local.target = local.backup }` to `module.backup_storage`.
#
# When done, remove the '# I AM NOT DONE' line at the top.
# ==============================================================================

terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = ">= 2.0.0"
    }
  }
}

provider "local" {
  alias = "backup"
}

module "backup_storage" {
  source    = "./modules/storage"
  file_name = "backup.log"

  # TODO: Pass the aliased provider to the child module
  # providers = {
  #   local.target = local.backup
  # }
}

output "backup_path" {
  value = module.backup_storage.path
}
