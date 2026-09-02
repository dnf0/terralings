# ==============================================================================
# Solution: gcp02
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

resource "google_compute_instance_template" "worker" {
  name_prefix  = "worker-template-"
  machine_type = "e2-standard-4"

  disk {
    source_image = "debian-cloud/debian-12"
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network = "default"
  }
}

resource "google_compute_region_instance_group_manager" "mig" {
  name               = "app-worker-mig"
  base_instance_name = "worker"
  region             = "us-central1"
  target_size        = 3

  version {
    instance_template = google_compute_instance_template.worker.id
  }
}

output "mig_name" {
  value = google_compute_region_instance_group_manager.mig.name
}
