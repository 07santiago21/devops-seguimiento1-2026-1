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
  default     = "prod"
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
  default     = "api-seguimiento3"
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

variable "bastion_key_name" {
  description = "Optional EC2 key pair name for bastion SSH access."
  type        = string
  default     = ""
}

variable "bastion_public_key_path" {
  description = "Path to a local SSH public key to import as the bastion EC2 key pair. Defaults to ~/.ssh/id_rsa.pub."
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}

variable "database_password" {
  description = "Master password for PostgreSQL RDS (provide your own, not auto-generated)."
  type        = string
  sensitive   = true
}

variable "aws_access_key" {
  description = "Optional AWS access key; prefer env/profile/role instead of using this variable."
  type        = string
  sensitive   = true
  default     = null
}

variable "aws_secret_key" {
  description = "Optional AWS secret key; prefer env/profile/role instead of using this variable."
  type        = string
  sensitive   = true
  default     = null
}

variable "eks_node_instance_type" {
  description = "EC2 instance type for EKS worker nodes."
  type        = string
  default     = "t3.small"
}

variable "eks_desired_nodes" {
  description = "Desired number of EKS worker nodes."
  type        = number
  default     = 2
}
