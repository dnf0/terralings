# ==============================================================================
# Exercise: gcp05
# Chapter: 15_gcp (Google Cloud (GCP) Architecture Blueprints)
#
# Task:
# Google Cloud IAM assigns least-privilege roles to Service Accounts. Using
# authoritative `google_project_iam_member` grants granular access without overwriting project policies:
#
#   resource "google_service_account" "runner" {
#     account_id   = "runner-sa"
#     display_name = "Cloud Run Invoker Identity"
#   }
#
#   resource "google_project_iam_member" "invoker" {
#     project = "my-project-123"
#     role    = "roles/run.invoker"
#     member  = "serviceAccount:${google_service_account.runner.email}"
#   }
#
# In this exercise:
# 1. Create `google_service_account.runner` with `account_id = "runner-sa"`.
# 2. Attach `google_project_iam_member.invoker` granting `role = "roles/run.invoker"`
#    to `member = "serviceAccount:${google_service_account.runner.email}"`.
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
  # TODO (What): Set account_id = "runner-sa" and display_name = "Cloud Run Invoker Identity".
  # TODO (Why): Service accounts represent non-human identities used by workloads and automated pipelines.
  account_id   = ""
  display_name = ""
}

resource "google_project_iam_member" "invoker" {
  project = var.project_id
  # TODO (What): Set role = "roles/run.invoker" and member = "serviceAccount:${google_service_account.runner.email}".
  # TODO (Why): Granting targeted IAM roles ensures workloads can only invoke permitted APIs.
  role    = ""
  member  = ""
}

output "service_account_email" {
  value = google_service_account.runner.email
}
