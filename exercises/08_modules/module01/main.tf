# ==============================================================================
# Exercise: module01
# Chapter: 08_modules (Modular Infrastructure Architecture)
#
# Task:
# Modules are the primary mechanism for packaging and reusing infrastructure
# configurations in OpenTofu and Terraform.
#
# A clean child module defines:
# 1. Inputs via variable blocks (with types and descriptions)
# 2. Resources implementing the internal logic
# 3. Outputs exposing values to calling parent modules
#
# In this exercise:
# 1. Inspect the child module in `modules/app_config/`
# 2. Define the missing variable `app_name` in `modules/app_config/variables.tf`
# 3. Define the output `full_name` in `modules/app_config/outputs.tf`
# 4. Instantiate the module in this root `main.tf` passing `app_name = "terralings"`
#
# ==============================================================================

module "app_config" {
  source = "./modules/app_config"

  # TODO (What): Pass app_name = "terralings" into the child module invocation.
  # TODO (Why): Child modules encapsulate functionality behind explicit input arguments, making configurations reusable.
  # app_name = "terralings"
}
