resource "terraform_data" "db" {
  input = {
    endpoint = "${var.db_name}.internal.db:5432"
  }
}
