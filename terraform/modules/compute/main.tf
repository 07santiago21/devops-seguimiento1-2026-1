data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

resource "aws_iam_role" "lambda_exec" {
  name = "${var.lambda_function_name}-exec-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "basic" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "vpc_access" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}

resource "aws_lambda_function" "this" {
  function_name = var.lambda_function_name
  filename      = var.artifact_path
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  role          = aws_iam_role.lambda_exec.arn

  source_code_hash = filebase64sha256(var.artifact_path)

  memory_size = var.memory_size
  timeout     = var.timeout

  vpc_config {
    subnet_ids         = var.subnet_ids
    security_group_ids = var.vpc_security_group_ids
  }

  environment {
    variables = {
      DATABASE_HOST     = var.database_host
      DATABASE_USER     = var.database_user
      DATABASE_PASSWORD = var.database_password
      DATABASE_NAME     = var.database_name
      DATABASE_PORT     = tostring(var.database_port)
    }
  }
}

resource "aws_api_gateway_rest_api" "this" {
  name        = "${var.lambda_function_name}-api"
  description = "API Gateway proxy for ${var.lambda_function_name}"
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_rest_api.this.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "root_any" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  resource_id   = aws_api_gateway_rest_api.this.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_method" "any" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  resource_id   = aws_api_gateway_resource.proxy.id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "proxy" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_method.any.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${data.aws_region.current.name}:lambda:path/2015-03-31/functions/${aws_lambda_function.this.arn}/invocations"
}

resource "aws_api_gateway_integration" "root_proxy" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_rest_api.this.root_resource_id
  http_method = aws_api_gateway_method.root_any.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_api_gateway_integration.proxy.uri
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.this.id}/*/*"
}

resource "aws_api_gateway_deployment" "this" {
  depends_on  = [aws_api_gateway_integration.proxy, aws_api_gateway_integration.root_proxy]
  rest_api_id = aws_api_gateway_rest_api.this.id

  triggers = {
    redeployment = sha1(join("-", [
      aws_api_gateway_integration.proxy.id,
      aws_api_gateway_integration.root_proxy.id,
      aws_api_gateway_method.any.id,
      aws_api_gateway_method.root_any.id,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "this" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  deployment_id = aws_api_gateway_deployment.this.id
  stage_name    = "prod"
}

output "lambda_function_arn" {
  value = aws_lambda_function.this.arn
}

output "lambda_function_name" {
  value = aws_lambda_function.this.function_name
}

output "api_invoke_url" {
  value = aws_api_gateway_stage.this.invoke_url
}
