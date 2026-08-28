terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = ">= 2.0.0"
      # TODO (What): Add configuration_aliases = [local.target].
      # TODO (Why): Declaring configuration_aliases informs the engine that this child module expects an explicit provider alias passed by parent modules.
    }
  }
}

resource "local_file" "log" {
  provider = local.target
  filename = "${path.module}/${var.file_name}"
  content  = "backup initialised\n"
}
