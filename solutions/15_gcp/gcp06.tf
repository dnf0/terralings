# ==============================================================================
# Solution: gcp06
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

resource "google_storage_bucket" "lake" {
  name                        = "corp-gcp-data-lake-99"
  location                    = "US"
  uniform_bucket_level_access = true
  force_destroy               = false

  versioning {
    enabled = true
  }
}

output "bucket_url" {
  value = google_storage_bucket.lake.url
}
