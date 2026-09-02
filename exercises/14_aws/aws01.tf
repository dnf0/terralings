# ==============================================================================
# Exercise: aws01
# Chapter: 14_aws (AWS Infrastructure & Production Blueprints)
#
# Task:
# Production AWS architectures start with a robust Virtual Private Cloud (VPC)
# spanning multiple Availability Zones with dedicated public and private subnets:
#
#   resource "aws_vpc" "main" {
#     cidr_block           = "10.0.0.0/16"
#     enable_dns_hostnames = true
#   }
#
#   resource "aws_subnet" "public_a" {
#     vpc_id                  = aws_vpc.main.id
#     cidr_block              = "10.0.1.0/24"
#     availability_zone       = "us-east-1a"
#     map_public_ip_on_launch = true
#   }
#
# In this exercise:
# 1. Complete `aws_vpc.main` with `cidr_block = "10.0.0.0/16"` and `enable_dns_hostnames = true`.
# 2. Complete `aws_subnet.public_a` with `cidr_block = "10.0.1.0/24"` and `map_public_ip_on_launch = true`.
# 3. Complete `aws_subnet.private_a` with `cidr_block = "10.0.10.0/24"` and `map_public_ip_on_launch = false`.
# 4. Attach `aws_route_table_association.public_a` associating `aws_subnet.public_a.id` with `aws_route_table.public.id`.
# ==============================================================================

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

resource "aws_vpc" "main" {
  # TODO (What): Set cidr_block = "10.0.0.0/16" and enable_dns_hostnames = true.
  # TODO (Why): The VPC CIDR defines the top-level IP space for the virtual cloud network.
  cidr_block           = ""
  enable_dns_hostnames = false
}

resource "aws_subnet" "public_a" {
  vpc_id            = aws_vpc.main.id
  # TODO (What): Set cidr_block = "10.0.1.0/24", availability_zone = "us-east-1a", and map_public_ip_on_launch = true.
  # TODO (Why): Public subnets auto-assign public IPs to internet-facing load balancers and NAT gateways.
  cidr_block        = ""
  availability_zone = "us-east-1a"
}

resource "aws_subnet" "private_a" {
  vpc_id                  = aws_vpc.main.id
  # TODO (What): Set cidr_block = "10.0.10.0/24", availability_zone = "us-east-1a", and map_public_ip_on_launch = false.
  # TODO (Why): Private subnets protect sensitive backend compute instances and databases from direct internet access.
  cidr_block              = ""
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }
}

resource "aws_route_table_association" "public_a" {
  # TODO (What): Link subnet_id = aws_subnet.public_a.id and route_table_id = aws_route_table.public.id.
  # TODO (Why): Route table associations route outbound traffic through the Internet Gateway.
  subnet_id      = ""
  route_table_id = ""
}

output "vpc_id" {
  value = aws_vpc.main.id
}
