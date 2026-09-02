# ==============================================================================
# Exercise: aws04
# Chapter: 14_aws (AWS Infrastructure & Production Blueprints)
#
# Task:
# Event-driven architectures decouple microservices using Amazon SQS FIFO queues
# and Amazon SNS topics with Dead-Letter Queues (DLQ) for resilient error handling:
#
#   resource "aws_sqs_queue" "dlq" {
#     name       = "orders-dlq.fifo"
#     fifo_queue = true
#   }
#
#   resource "aws_sqs_queue" "orders" {
#     name           = "orders.fifo"
#     fifo_queue     = true
#     redrive_policy = jsonencode({
#       deadLetterTargetArn = aws_sqs_queue.dlq.arn
#       maxReceiveCount     = 5
#     })
#   }
#
# In this exercise:
# 1. Create `aws_sqs_queue.orders_dlq` with `name = "orders-dlq.fifo"` and `fifo_queue = true`.
# 2. Configure `aws_sqs_queue.orders` with `fifo_queue = true`, `content_based_deduplication = true`,
#    and attach the `redrive_policy` pointing to the DLQ ARN.
# 3. Create `aws_sns_topic.orders_topic` with `fifo_topic = true` and `name = "orders-topic.fifo"`.
# ==============================================================================

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

resource "aws_sqs_queue" "orders_dlq" {
  name       = "orders-dlq.fifo"
  # TODO (What): Set fifo_queue = true.
  # TODO (Why): Dead letter queues for FIFO topics must also be configured with FIFO semantics.
  fifo_queue = false
}

resource "aws_sqs_queue" "orders" {
  name                        = "orders.fifo"
  # TODO (What): Set fifo_queue = true and content_based_deduplication = true.
  # TODO (Why): FIFO queues guarantee exactly-once message delivery and strict ordering.
  fifo_queue                  = false
  content_based_deduplication = false

  # TODO (What): Set redrive_policy = jsonencode({ deadLetterTargetArn = aws_sqs_queue.orders_dlq.arn, maxReceiveCount = 5 }).
  # TODO (Why): Redrive policies automatically redirect poison-pill messages after repeated failures.
  redrive_policy = ""
}

resource "aws_sns_topic" "orders_topic" {
  name       = "orders-topic.fifo"
  # TODO (What): Set fifo_topic = true.
  # TODO (Why): SNS FIFO topics publish ordered events across distributed fan-out subscriber queues.
  fifo_topic = false
}

output "queue_url" {
  value = aws_sqs_queue.orders.url
}
