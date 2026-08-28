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

  # TODO (What): Compute json_config = jsonencode(local.app_config).
  # TODO (Why): jsonencode() transforms native HCL maps and lists into valid JSON strings without manual string formatting.
  json_config = ""

  # TODO (What): Compute yaml_config = yamlencode(local.app_config).
  # TODO (Why): yamlencode() converts structured data into canonical YAML format required by Kubernetes manifests and CI configs.
  yaml_config = ""

  # TODO (What): Compute b64_json = base64encode(local.json_config).
  # TODO (Why): base64encode() encodes string payloads for cloud APIs, userdata initialization scripts, and binary secrets.
  b64_json = ""

  # TODO (What): Compute parsed_json = jsondecode(local.json_config).
  # TODO (Why): jsondecode() parses external JSON payloads back into navigable HCL maps and objects.
  parsed_json = {}

  # TODO (What): Compute decoded_b64 = base64decode(local.b64_json).
  # TODO (Why): base64decode() retrieves original plaintext data from encoded base64 strings.
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
