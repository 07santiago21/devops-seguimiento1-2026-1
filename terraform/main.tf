terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }

    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.4"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = local.common_tags
  }
}

locals {
  name_prefix = "${var.project_name}-${var.environment}"

  common_tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
  }
}

module "network" {
  source           = "./modules/network"
  name_prefix      = local.name_prefix
  allowed_ssh_cidr = var.allowed_ssh_cidr
  tags             = local.common_tags
}

# Generate a random password for the RDS master user
resource "random_password" "rds" {
  length  = 20
  special = true
}

module "database" {
  source = "./modules/database"

  name_prefix            = local.name_prefix
  subnet_ids             = [module.network.private_subnet_id]
  vpc_security_group_ids = [module.network.rds_sg_id]

  db_name  = var.database_name
  username = var.database_username
  password = random_password.rds.result

  engine_version    = "15"
  instance_class    = "db.t3.micro"
  allocated_storage = 20

  tags = local.common_tags
}

module "compute" {
  source = "./modules/compute"

  lambda_function_name = var.lambda_function_name
  artifact_path        = var.artifact_path

  subnet_ids             = [module.network.private_subnet_id]
  vpc_security_group_ids = [module.network.lambda_sg_id]

  database_host     = module.database.endpoint
  database_user     = module.database.username
  database_password = module.database.password
  database_name     = module.database.db_name
  database_port     = module.database.port

  memory_size = 128
  timeout     = 10
}
