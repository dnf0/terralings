resource "terraform_data" "service" {
  input = {
    name = var.service_name
    port = var.port
  }
}
