output "vpc_id" {
  description = "VPC id created by the network module"
  value       = aws_vpc.this.id
}

output "public_subnet_id" {
  description = "Public subnet id"
  value       = aws_subnet.public.id
}

output "private_subnet_id" {
  description = "Private subnet id"
  value       = aws_subnet.private.id
}

output "bastion_public_ip" {
  description = "Public IP of the bastion host"
  value       = aws_instance.bastion.public_ip
}

output "bastion_sg_id" {
  description = "Security group id for bastion"
  value       = aws_security_group.bastion_sg.id
}

output "lambda_sg_id" {
  description = "Security group id for Lambda"
  value       = aws_security_group.lambda_sg.id
}

output "rds_sg_id" {
  description = "Security group id for RDS"
  value       = aws_security_group.rds_sg.id
}
