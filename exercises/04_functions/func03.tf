# I AM NOT DONE
# ==============================================================================
# Exercise: func03
# Chapter: 04_functions (Built-in Functions & Collections)
#
# Task:
# Infrastructure definitions frequently need to serialize and deserialize data
# for cloud metadata, IAM policies, and configuration files:
# - jsonencode(val) converts an HCL value into a JSON string
# - yamlencode(val) converts an HCL value into a YAML string
# - base64encode(str) encodes a string as Base64 (useful for user_data)
# - jsondecode(str) parses a JSON string back into HCL structures
# - base64decode(str) decodes a Base64 string back into raw text
#
# Complete the locals block below:
# 1. json_config: serialize local.app_config using jsonencode()
# 2. yaml_config: serialize local.app_config using yamlencode()
# 3. b64_json: encode local.json_config using base64encode()
# 4. parsed_json: deserialize local.json_config using jsondecode()
# 5. decoded_b64: decode local.b64_json using base64decode()
#
# When done, remove the '# I AM NOT DONE' line at the top.
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

  # TODO: jsonencode local.app_config
  json_config = ""

  # TODO: yamlencode local.app_config
  yaml_config = ""

  # TODO: base64encode local.json_config
  b64_json = ""

  # TODO: jsondecode local.json_config
  parsed_json = {}

  # TODO: base64decode local.b64_json
  decoded_b64 = ""
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
