# ==============================================================================
# Exercise: aws03
# Chapter: 14_aws (AWS Infrastructure & Production Blueprints)
#
# Task:
# Serverless microservices on AWS combine AWS Lambda functions with API Gateway v2 HTTP APIs
# and least-privilege IAM execution roles:
#
#   resource "aws_lambda_function" "api" {
#     function_name = "order-processor"
#     role          = aws_iam_role.lambda_exec.arn
#     runtime       = "python3.12"
#     handler       = "index.handler"
#   }
#
# In this exercise:
# 1. Define `aws_iam_role.lambda_exec` with an AssumeRole policy for `lambda.amazonaws.com`.
# 2. Configure `aws_lambda_function.order_api` with `runtime = "python3.12"`, `handler = "index.handler"`,
#    and reference the IAM role ARN.
# 3. Create `aws_apigatewayv2_api.http_api` with `protocol_type = "HTTP"`.
# 4. Create `aws_apigatewayv2_integration.lambda_link` connecting the HTTP API to the Lambda invoke ARN.
# ==============================================================================

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

resource "aws_iam_role" "lambda_exec" {
  name = "lambda-exec-role"

  # TODO (What): Provide jsonencode trust policy granting sts:AssumeRole to lambda.amazonaws.com.
  # TODO (Why): Lambda requires execution role trust to assume identity and write CloudWatch logs.
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
    }]
  })
}

resource "aws_lambda_function" "order_api" {
  function_name = "order-processor"
  role          = aws_iam_role.lambda_exec.arn
  # TODO (What): Set runtime = "python3.12", handler = "index.handler", and timeout = 15.
  # TODO (Why): Defines the execution environment and entrypoint for the serverless function.
  runtime       = ""
  handler       = ""
  timeout       = 0
  filename      = "function.zip"
}

resource "aws_apigatewayv2_api" "http_api" {
  name          = "orders-http-api"
  # TODO (What): Set protocol_type = "HTTP".
  # TODO (Why): API Gateway v2 HTTP APIs provide lightweight, low-latency endpoints for serverless backends.
  protocol_type = ""
}

resource "aws_apigatewayv2_integration" "lambda_link" {
  api_id           = aws_apigatewayv2_api.http_api.id
  integration_type = "AWS_PROXY"
  # TODO (What): Set integration_uri = aws_lambda_function.order_api.invoke_arn.
  # TODO (Why): Connects incoming API Gateway requests to the Lambda function execution engine.
  integration_uri  = ""
}

output "api_endpoint" {
  value = aws_apigatewayv2_api.http_api.api_endpoint
}
