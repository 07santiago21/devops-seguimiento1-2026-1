variable "aws_region" {
  description = "AWS region where all resources will be created."
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Logical project name used as the base for resource names."
  type        = string
  default     = "devops-seguimiento3"
}

variable "environment" {
  description = "Deployment environment name."
  type        = string
  default     = "dev"
}


variable "database_name" {
  description = "PostgreSQL database name that will be used later by the RDS module."
  type        = string
  default     = "academia"
}

variable "database_username" {
  description = "Master username for PostgreSQL."
  type        = string
  default     = "postgres"
}

variable "database_port" {
  description = "PostgreSQL port exposed by the database module."
  type        = number
  default     = 5432
}

variable "lambda_function_name" {
  description = "Lambda function name to be created in the compute module."
  type        = string
  default     = "api-seguimiento1"
}

variable "artifact_path" {
  description = "Path to the packaged Lambda ZIP artifact."
  type        = string
  default     = "../dist/function.zip"
}

variable "allowed_ssh_cidr" {
  description = "CIDR allowed to reach the bastion host via SSH."
  type        = string
  default     = "0.0.0.0/0"
}
