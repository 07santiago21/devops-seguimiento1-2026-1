variable "name_prefix" {
  type        = string
  description = "Prefix for resource names"
}

variable "subnet_ids" {
  type        = list(string)
  description = "List of private subnet IDs for the RDS subnet group"
}

variable "vpc_security_group_ids" {
  type        = list(string)
  description = "Security group IDs to attach to the RDS instance (should include Bastion and Lambda SGs)"
  default     = []
}

variable "allocated_storage" {
  type    = number
  default = 20
}

variable "engine" {
  type    = string
  default = "postgres"
}

variable "engine_version" {
  type    = string
  default = "15"
}

variable "instance_class" {
  type    = string
  default = "db.t3.micro"
}

variable "db_name" {
  type        = string
  description = "Database name"
}

variable "username" {
  type        = string
  description = "Master username for the database"
}

variable "password" {
  type        = string
  description = "Master password for the database"
  sensitive   = true
}

variable "skip_final_snapshot" {
  type    = bool
  default = true
}

variable "tags" {
  type    = map(string)
  default = {}
}
