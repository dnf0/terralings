# ==============================================================================
# Solution: gcp03
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

resource "google_cloud_run_v2_service" "orders_api" {
  name     = "orders-api"
  location = "us-central1"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }

    scaling {
      min_instance_count = 1
      max_instance_count = 20
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
}

output "service_uri" {
  value = google_cloud_run_v2_service.orders_api.uri
}
