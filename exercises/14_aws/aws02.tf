# ==============================================================================
# Exercise: aws02
# Chapter: 14_aws (AWS Infrastructure & Production Blueprints)
#
# Task:
# Production EC2 workloads utilize Launch Templates paired with Auto Scaling Groups (ASG)
# and Application Load Balancer Target Groups to ensure high availability and elastic scaling:
#
#   resource "aws_launch_template" "app" {
#     name_prefix   = "app-server-"
#     image_id      = "ami-0123456789abcdef0"
#     instance_type = "t3.medium"
#   }
#
#   resource "aws_autoscaling_group" "app" {
#     min_size         = 2
#     max_size         = 10
#     desired_capacity = 4
#   }
#
# In this exercise:
# 1. Define `aws_launch_template.app` with `instance_type = "t3.medium"` and `image_id = "ami-0123456789abcdef0"`.
# 2. Configure `aws_lb_target_group.app` on port 80 with protocol "HTTP".
# 3. Configure `aws_autoscaling_group.app` with `min_size = 2`, `max_size = 10`, `desired_capacity = 4`,
#    and attach the target group ARN and launch template.
# ==============================================================================

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

resource "aws_launch_template" "app" {
  name_prefix   = "app-server-"
  # TODO (What): Set image_id = "ami-0123456789abcdef0" and instance_type = "t3.medium".
  # TODO (Why): Launch templates standardize machine configuration, AMI selection, and instance sizing.
  image_id      = ""
  instance_type = ""
}

resource "aws_lb_target_group" "app" {
  name     = "app-tg"
  # TODO (What): Set port = 80 and protocol = "HTTP".
  # TODO (Why): Target groups route traffic to healthy backend instances on the application service port.
  port     = 0
  protocol = ""
  vpc_id   = "vpc-12345678"

  health_check {
    path     = "/healthz"
    matcher  = "200"
    interval = 30
  }
}

resource "aws_autoscaling_group" "app" {
  name_prefix         = "app-asg-"
  vpc_zone_identifier = ["subnet-1a", "subnet-1b"]
  # TODO (What): Set target_group_arns = [aws_lb_target_group.app.arn], min_size = 2, max_size = 10, desired_capacity = 4.
  # TODO (Why): Auto Scaling Groups balance capacity across availability zones and register targets with the load balancer.
  target_group_arns   = []
  min_size            = 0
  max_size            = 0
  desired_capacity    = 0

  launch_template {
    id      = aws_launch_template.app.id
    version = "$Latest"
  }
}

output "asg_name" {
  value = aws_autoscaling_group.app.name
}
