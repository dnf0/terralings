variable "environment" {
  type        = string
  description = "Deployment environment"
  default     = "dev"
}

# TODO (What): Declare the app_name variable of type string with a description.
# TODO (Why): Declaring explicit variable definitions in child modules defines the module's public input contract.
