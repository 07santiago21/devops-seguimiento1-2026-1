output "project_name" {
  description = "Project name used as the base for resource naming."
  value       = var.project_name
}

output "environment" {
  description = "Deployment environment."
  value       = var.environment
}

output "aws_region" {
  description = "AWS region configured in the root provider."
  value       = var.aws_region
}

output "name_prefix" {
  description = "Computed name prefix shared by all modules."
  value       = local.name_prefix
}

output "lambda_function_name" {
  description = "Lambda function name reserved for the compute layer."
  value       = var.lambda_function_name
}

output "artifact_path" {
  description = "Path to the Lambda deployment package."
  value       = var.artifact_path
}

output "database_name" {
  description = "Database name reserved for the PostgreSQL layer."
  value       = var.database_name
}

output "database_username" {
  description = "Database username reserved for the PostgreSQL layer."
  value       = var.database_username
}

output "api_invoke_url" {
  description = "API Gateway endpoint for the deployed Lambda (prod stage)."
  value       = module.compute.api_invoke_url
}

output "lambda_arn" {
  description = "ARN of the deployed Lambda function."
  value       = module.compute.lambda_function_arn
}

output "lambda_name" {
  description = "Name of the deployed Lambda function."
  value       = module.compute.lambda_function_name
}
