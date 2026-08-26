# ==============================================================================
# Solution: module04
# Chapter: 08_modules (Modular Infrastructure Architecture)
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

  providers = {
    local.target = local.backup
  }
}

output "backup_path" {
  value = module.backup_storage.path
}
