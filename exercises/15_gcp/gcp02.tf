# ==============================================================================
# Exercise: gcp02
# Chapter: 15_gcp (Google Cloud (GCP) Architecture Blueprints)
#
# Task:
# Scalable VM architectures in GCP rely on Regional Managed Instance Groups (MIGs)
# that deploy across multiple zones using an Instance Template and Health Checks:
#
#   resource "google_compute_instance_template" "worker" {
#     name_prefix  = "worker-template-"
#     machine_type = "e2-standard-4"
#   }
#
#   resource "google_compute_region_instance_group_manager" "mig" {
#     name        = "app-worker-mig"
#     region      = "us-central1"
#     target_size = 3
#   }
#
# In this exercise:
# 1. Complete `google_compute_instance_template.worker` with `machine_type = "e2-standard-4"`.
# 2. Configure `google_compute_region_instance_group_manager.mig` with `target_size = 3`
#    and link the instance template ID.
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
  # TODO (What): Set machine_type = "e2-standard-4".
  # TODO (Why): Instance templates specify machine types, boot disks, and metadata for all managed instances.
  machine_type = ""

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
  # TODO (What): Set target_size = 3 and version { instance_template = google_compute_instance_template.worker.id }.
  # TODO (Why): Regional MIGs distribute instances automatically across availability zones in the region.
  target_size        = 0

  version {
    instance_template = ""
  }
}

output "mig_name" {
  value = google_compute_region_instance_group_manager.mig.name
}
