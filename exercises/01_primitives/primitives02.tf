# ==============================================================================
# Exercise: primitives02
# Chapter: 01_primitives (HCL Foundations & Core Primitives)
#
# Task:
# Resources are the most important element in HCL configurations.
# A resource block declares a resource of a given type ("local_file") with a
# given local name ("welcome").
#
# Fix the resource block below so that:
# 1. It declares resource "local_file" "welcome"
# 2. Sets filename = "${path.module}/welcome.txt"
# 3. Sets content  = "Welcome to Terralings!"
#
# ==============================================================================

terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.0"
    }
  }
}

# TODO (What): Add the 'content' attribute to resource "local_file" "welcome" set to "Welcome to Terralings!".
# TODO (Why): The local_file provider manages files on disk; without content (or content_base64),
#             provider schema validation fails because it cannot determine what data to write.
resource "local_file" "welcome" {
  filename = "${path.module}/welcome.txt"
  # content is missing!
}
