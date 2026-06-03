variable "name_prefix" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "rds_sg_id" {
  type = string
}

variable "kubernetes_version" {
  type    = string
  default = "1.31"
}

variable "node_instance_type" {
  type    = string
  default = "t3.small"
}

variable "desired_nodes" {
  type    = number
  default = 2
}

variable "max_nodes" {
  type    = number
  default = 3
}

variable "tags" {
  type    = map(string)
  default = {}
}
