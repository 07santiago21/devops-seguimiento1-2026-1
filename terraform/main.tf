terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }

    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.4"
    }
  }
}

provider "aws" {
  region     = var.aws_region
  access_key = var.aws_access_key != null && var.aws_access_key != "" ? var.aws_access_key : null
  secret_key = var.aws_secret_key != null && var.aws_secret_key != "" ? var.aws_secret_key : null

  default_tags {
    tags = local.common_tags
  }
}

locals {
  name_prefix = "${var.project_name}-${var.environment}"

  # Resuelve ~ al directorio home y verifica que el archivo exista
  _resolved_key_path = var.bastion_public_key_path != "" ? pathexpand(var.bastion_public_key_path) : ""
  _key_file_exists   = local._resolved_key_path != "" && fileexists(local._resolved_key_path)

  bastion_key_name = var.bastion_key_name != "" ? var.bastion_key_name : try(aws_key_pair.bastion[0].key_name, "")

  common_tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
  }
}

resource "aws_key_pair" "bastion" {
  count      = local._key_file_exists ? 1 : 0
  key_name   = "${local.name_prefix}-bastion"
  public_key = file(local._resolved_key_path)
}

module "network" {
  source           = "./modules/network"
  name_prefix      = local.name_prefix
  allowed_ssh_cidr = var.allowed_ssh_cidr
  bastion_key_name = local.bastion_key_name
  tags             = local.common_tags
}

module "database" {
  source = "./modules/database"

  name_prefix            = local.name_prefix
  subnet_ids             = module.network.private_subnet_ids
  vpc_security_group_ids = [module.network.rds_sg_id]

  db_name  = var.database_name
  username = var.database_username
  password = var.database_password

  engine_version    = "15"
  instance_class    = "db.t3.micro"
  allocated_storage = 20

  tags = local.common_tags
}

module "compute" {
  source = "./modules/compute"

  lambda_function_name = var.lambda_function_name
  artifact_path        = var.artifact_path

  subnet_ids             = module.network.private_subnet_ids
  vpc_security_group_ids = [module.network.lambda_sg_id]

  database_host     = module.database.endpoint
  database_user     = module.database.username
  database_password = module.database.password
  database_name     = module.database.db_name
  database_port     = module.database.port

  memory_size = 128
  timeout     = 10
}

module "ecr" {
  source          = "./modules/ecr"
  repository_name = "${local.name_prefix}-api"
  tags            = local.common_tags
}

module "eks" {
  source             = "./modules/eks"
  name_prefix        = local.name_prefix
  subnet_ids         = module.network.public_subnet_ids
  rds_sg_id          = module.network.rds_sg_id
  node_instance_type = var.eks_node_instance_type
  desired_nodes      = var.eks_desired_nodes
  tags               = local.common_tags
}
