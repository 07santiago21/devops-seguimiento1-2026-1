variable "name_prefix" {
  description = "Prefix used to name network resources."
  type        = string
}

variable "vpc_cidr" {
  description = "VPC CIDR block."
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidr" {
  description = "CIDR for the public subnet."
  type        = string
  default     = "10.0.1.0/24"
}

variable "private_subnet_cidr" {
  description = "CIDR for the private subnet."
  type        = string
  default     = "10.0.2.0/24"
}

variable "allowed_ssh_cidr" {
  description = "CIDR allowed to reach the bastion host via SSH."
  type        = string
  default     = "0.0.0.0/0"
}

variable "bastion_instance_type" {
  description = "EC2 instance type for the bastion host."
  type        = string
  default     = "t3.micro"
}

variable "bastion_key_name" {
  description = "Optional EC2 key pair name for bastion SSH access."
  type        = string
  default     = ""
}

variable "tags" {
  description = "Map of tags to apply to resources."
  type        = map(string)
  default     = {}
}
