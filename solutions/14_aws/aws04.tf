# ==============================================================================
# Solution: aws04
# Chapter: 14_aws (AWS Infrastructure & Production Blueprints)
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
  fifo_queue = true
}

resource "aws_sqs_queue" "orders" {
  name                        = "orders.fifo"
  fifo_queue                  = true
  content_based_deduplication = true

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.orders_dlq.arn
    maxReceiveCount     = 5
  })
}

resource "aws_sns_topic" "orders_topic" {
  name       = "orders-topic.fifo"
  fifo_topic = true
}

output "queue_url" {
  value = aws_sqs_queue.orders.url
}
