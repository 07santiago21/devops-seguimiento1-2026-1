resource "aws_db_subnet_group" "this" {
  name       = "${var.name_prefix}-rds-subnet-group"
  subnet_ids = var.subnet_ids

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-rds-subnet-group"
  })
}

resource "aws_db_parameter_group" "this" {
  name        = "${var.name_prefix}-postgres-params"
  family      = "postgres${var.engine_version}"
  description = "Custom parameter group for ${var.name_prefix} postgres"

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-postgres-params"
  })
}

resource "aws_db_instance" "this" {
  identifier             = "${var.name_prefix}-postgres"
  allocated_storage      = var.allocated_storage
  engine                 = var.engine
  engine_version         = var.engine_version
  instance_class         = var.instance_class
  db_name                = var.db_name
  username               = var.username
  password               = var.password
  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = var.vpc_security_group_ids
  skip_final_snapshot    = var.skip_final_snapshot
  publicly_accessible    = false
  parameter_group_name   = aws_db_parameter_group.this.name

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-postgres"
  })
}
