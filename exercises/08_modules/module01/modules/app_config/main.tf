resource "terraform_data" "config" {
  input = "${var.environment}-${var.app_name}"
}
