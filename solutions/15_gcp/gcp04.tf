# ==============================================================================
# Solution: gcp04
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

resource "google_pubsub_topic" "events" {
  name = "telemetry-events"
}

resource "google_pubsub_topic" "events_dlq" {
  name = "telemetry-events-dlq"
}

resource "google_pubsub_subscription" "events_sub" {
  name                 = "telemetry-worker-sub"
  topic                = google_pubsub_topic.events.name
  ack_deadline_seconds = 20

  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.events_dlq.id
    max_delivery_attempts = 5
  }
}

output "subscription_id" {
  value = google_pubsub_subscription.events_sub.id
}
