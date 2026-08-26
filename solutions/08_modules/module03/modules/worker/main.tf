resource "terraform_data" "worker" {
  input = {
    id          = "worker-${var.worker_name}"
    concurrency = var.concurrency
  }
}
