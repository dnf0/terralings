# ==============================================================================
# Solution: gcp01
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

resource "google_compute_network" "vpc" {
  name                    = "vpc-production"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "app_subnet" {
  name                     = "subnet-us-central1-app"
  network                  = google_compute_network.vpc.id
  region                   = "us-central1"
  ip_cidr_range            = "10.10.0.0/20"
  private_ip_google_access = true
}

resource "google_compute_firewall" "allow_internal" {
  name    = "allow-internal-app"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["443"]
  }

  source_ranges = [google_compute_subnetwork.app_subnet.ip_cidr_range]
  target_tags   = ["app-backend"]
}

output "network_id" {
  value = google_compute_network.vpc.id
}
