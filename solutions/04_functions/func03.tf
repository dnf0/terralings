# ==============================================================================
# Solution: func03
# Chapter: 04_functions (Built-in Functions & Collections)
# ==============================================================================

terraform {
  required_version = ">= 1.6.0"
}

locals {
  app_config = {
    service     = "terralings"
    port        = 8080
    debug_mode  = false
    maintainers = ["devops", "platform"]
  }

  json_config = jsonencode(local.app_config)
  yaml_config = yamlencode(local.app_config)
  b64_json    = base64encode(local.json_config)
  parsed_json = jsondecode(local.json_config)
  decoded_b64 = base64decode(local.b64_json)
}

resource "terraform_data" "encodings" {
  input = {
    json        = local.json_config
    yaml        = local.yaml_config
    b64         = local.b64_json
    parsed_port = local.parsed_json["port"]
    decoded     = local.decoded_b64
  }
}
