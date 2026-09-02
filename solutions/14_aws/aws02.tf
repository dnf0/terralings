# ==============================================================================
# Solution: aws02
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

resource "aws_launch_template" "app" {
  name_prefix   = "app-server-"
  image_id      = "ami-0123456789abcdef0"
  instance_type = "t3.medium"
}

resource "aws_lb_target_group" "app" {
  name     = "app-tg"
  port     = 80
  protocol = "HTTP"
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
  target_group_arns   = [aws_lb_target_group.app.arn]
  min_size            = 2
  max_size            = 10
  desired_capacity    = 4

  launch_template {
    id      = aws_launch_template.app.id
    version = "$Latest"
  }
}

output "asg_name" {
  value = aws_autoscaling_group.app.name
}
