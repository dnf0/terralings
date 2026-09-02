# ==============================================================================
# Solution: gcp05
# Chapter: 15_gcp (Google Cloud (GCP) Architecture Blueprints)
# ==============================================================================

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

variable "project_id" {
  type    = string
  default = "my-gcp-project-99"
}

resource "google_service_account" "runner" {
  account_id   = "runner-sa"
  display_name = "Cloud Run Invoker Identity"
}

resource "google_project_iam_member" "invoker" {
  project = var.project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.runner.email}"
}

output "service_account_email" {
  value = google_service_account.runner.email
}
