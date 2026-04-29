variable "lambda_function_name" {
  description = "Name for the Lambda function"
  type        = string
}

variable "artifact_path" {
  description = "Path to the Lambda deployment ZIP artifact"
  type        = string
}

variable "subnet_ids" {
  description = "List of private subnet IDs where the Lambda will run"
  type        = list(string)
}

variable "vpc_security_group_ids" {
  description = "Security groups to attach to the Lambda function"
  type        = list(string)
}

variable "database_host" {
  description = "RDS endpoint address to inject into Lambda"
  type        = string
}

variable "database_user" {
  description = "DB user for Lambda to connect"
  type        = string
}

variable "database_password" {
  description = "DB password (sensitive)"
  type        = string
  sensitive   = true
}

variable "database_name" {
  description = "DB name"
  type        = string
}

variable "database_port" {
  description = "DB port"
  type        = number
  default     = 5432
}

variable "memory_size" {
  description = "Lambda memory size (MB)"
  type        = number
  default     = 128
}

variable "timeout" {
  description = "Lambda timeout in seconds"
  type        = number
  default     = 10
}
