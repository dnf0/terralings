# ==============================================================================
# Exercise: gcp03
# Chapter: 15_gcp (Google Cloud (GCP) Architecture Blueprints)
#
# Task:
# Cloud Run v2 deploys containerized serverless microservices with automated
# request concurrency, traffic splitting, and scale-to-zero capabilities:
#
#   resource "google_cloud_run_v2_service" "api" {
#     name     = "orders-api"
#     location = "us-central1"
#     template {
#       scaling {
#         max_instance_count = 20
#       }
#     }
#   }
#
# In this exercise:
# 1. Complete `google_cloud_run_v2_service.orders_api` with `location = "us-central1"`.
# 2. In `template.scaling`, configure `min_instance_count = 1` and `max_instance_count = 20`.
# 3. Configure `traffic` with `type = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"` and `percent = 100`.
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
  # TODO (What): Set location = "us-central1".
  # TODO (Why): Cloud Run services run in regional serverless compute regions.
  location = ""

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
    }

    scaling {
      # TODO (What): Set min_instance_count = 1 and max_instance_count = 20.
      # TODO (Why): Autoscaling boundaries prevent unbounded cloud spend while keeping warm instances available.
      min_instance_count = 0
      max_instance_count = 0
    }
  }

  traffic {
    # TODO (What): Set type = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST" and percent = 100.
    # TODO (Why): Routes 100% of incoming production HTTP traffic to the latest deployed revision.
    type    = ""
    percent = 0
  }
}

output "service_uri" {
  value = google_cloud_run_v2_service.orders_api.uri
}
