terraform {
  required_providers {
    local = {
      source                = "hashicorp/local"
      version               = ">= 2.0.0"
      configuration_aliases = [local.target]
    }
  }
}

resource "local_file" "log" {
  provider = local.target
  filename = "${path.module}/${var.file_name}"
  content  = "backup initialised\n"
}
