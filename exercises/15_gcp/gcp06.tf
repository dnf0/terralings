# ==============================================================================
# Exercise: gcp06
# Chapter: 15_gcp (Google Cloud (GCP) Architecture Blueprints)
#
# Task:
# Production Google Cloud Storage buckets enforce Uniform Bucket-Level Access
# and object versioning to protect critical data assets:
#
#   resource "google_storage_bucket" "lake" {
#     name                        = "corp-gcp-data-lake-99"
#     location                    = "US"
#     uniform_bucket_level_access = true
#     versioning {
#       enabled = true
#     }
#   }
#
# In this exercise:
# 1. Create `google_storage_bucket.lake` with `name = "corp-gcp-data-lake-99"`, `location = "US"`,
#    `uniform_bucket_level_access = true`, and `force_destroy = false`.
# 2. Add `versioning { enabled = true }` inside `google_storage_bucket.lake`.
# ==============================================================================

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

resource "google_storage_bucket" "lake" {
  # TODO (What): Set name = "corp-gcp-data-lake-99" and location = "US".
  # TODO (Why): Storage buckets require globally unique names and explicit regional/multi-regional locations.
  name     = ""
  location = ""

  # TODO (What): Set uniform_bucket_level_access = true and force_destroy = false.
  # TODO (Why): Uniform access simplifies security by disabling individual object ACLs and unifying policies under IAM.
  uniform_bucket_level_access = false
  force_destroy               = true

  versioning {
    # TODO (What): Set enabled = true.
    # TODO (Why): Versioning protects objects against accidental overwrites or malicious deletion.
    enabled = false
  }
}

output "bucket_url" {
  value = google_storage_bucket.lake.url
}
