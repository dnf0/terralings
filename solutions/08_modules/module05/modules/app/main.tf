resource "terraform_data" "app" {
  input = {
    service = var.app_name
    connect = "postgres://${var.db_host}/${var.app_name}"
  }
}
