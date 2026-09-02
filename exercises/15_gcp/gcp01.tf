# ==============================================================================
# Exercise: gcp01
# Chapter: 15_gcp (Google Cloud (GCP) Architecture Blueprints)
#
# Task:
# Google Cloud VPCs are global software-defined networks. Enterprise architectures
# require custom-mode networks with regional subnets and tag-filtered firewall rules:
#
#   resource "google_compute_network" "custom_vpc" {
#     name                    = "vpc-production"
#     auto_create_subnetworks = false
#   }
#
#   resource "google_compute_subnetwork" "app_subnet" {
#     name                     = "subnet-us-central1-app"
#     ip_cidr_range            = "10.10.0.0/20"
#     region                   = "us-central1"
#     network                  = google_compute_network.custom_vpc.id
#     private_ip_google_access = true
#   }
#
# In this exercise:
# 1. Complete `google_compute_network.vpc` with `name = "vpc-production"` and `auto_create_subnetworks = false`.
# 2. Complete `google_compute_subnetwork.app_subnet` with `ip_cidr_range = "10.10.0.0/20"`, `region = "us-central1"`,
#    and `private_ip_google_access = true`.
# 3. Complete `google_compute_firewall.allow_internal` allowing TCP port 443 with `target_tags = ["app-backend"]`.
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
  name = "vpc-production"
  # TODO (What): Set auto_create_subnetworks = false.
  # TODO (Why): Custom VPCs prevent default allocation of subnets in unused cloud regions.
  auto_create_subnetworks = true
}

resource "google_compute_subnetwork" "app_subnet" {
  name    = "subnet-us-central1-app"
  network = google_compute_network.vpc.id
  region  = "us-central1"

  # TODO (What): Set ip_cidr_range = "10.10.0.0/20" and private_ip_google_access = true.
  # TODO (Why): Private Google Access enables VMs without external IPs to reach Google APIs over internal network routing.
  ip_cidr_range            = ""
  private_ip_google_access = false
}

resource "google_compute_firewall" "allow_internal" {
  name    = "allow-internal-app"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    # TODO (What): Set ports = ["443"].
    # TODO (Why): Scopes inbound access to secure encrypted endpoints.
    ports    = []
  }

  source_ranges = [google_compute_subnetwork.app_subnet.ip_cidr_range]
  # TODO (What): Set target_tags = ["app-backend"].
  # TODO (Why): GCP firewall rules use network tags to dynamically target specific compute instance pools.
  target_tags   = []
}

output "network_id" {
  value = google_compute_network.vpc.id
}
