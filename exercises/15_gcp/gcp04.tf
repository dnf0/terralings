# ==============================================================================
# Exercise: gcp04
# Chapter: 15_gcp (Google Cloud (GCP) Architecture Blueprints)
#
# Task:
# Google Cloud Pub/Sub provides scalable asynchronous messaging. Production pipelines
# attach Dead-Letter Topics to subscriptions to catch poisoned messages:
#
#   resource "google_pubsub_subscription" "events_sub" {
#     name  = "events-worker-sub"
#     topic = google_pubsub_topic.events.name
#     dead_letter_policy {
#       dead_letter_topic     = google_pubsub_topic.events_dlq.id
#       max_delivery_attempts = 5
#     }
#   }
#
# In this exercise:
# 1. Create `google_pubsub_topic.events` with `name = "telemetry-events"`.
# 2. Create `google_pubsub_topic.events_dlq` with `name = "telemetry-events-dlq"`.
# 3. Create `google_pubsub_subscription.events_sub` with `ack_deadline_seconds = 20`,
#    and attach the `dead_letter_policy` pointing to the DLQ topic with `max_delivery_attempts = 5`.
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
  name  = "telemetry-worker-sub"
  topic = google_pubsub_topic.events.name

  # TODO (What): Set ack_deadline_seconds = 20.
  # TODO (Why): The acknowledgement deadline defines the window workers have to process messages.
  ack_deadline_seconds = 0

  dead_letter_policy {
    # TODO (What): Set dead_letter_topic = google_pubsub_topic.events_dlq.id and max_delivery_attempts = 5.
    # TODO (Why): Dead letter policies isolate unprocessable messages after threshold retries.
    dead_letter_topic     = ""
    max_delivery_attempts = 0
  }
}

output "subscription_id" {
  value = google_pubsub_subscription.events_sub.id
}
