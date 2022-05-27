terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = ">= 3.40"
    }
  }

}

provider "aws" {
    shared_credentials_file = "${var.tag_credential_filepath}"
    profile                 = "aws-tag-${var.tag_environment}"
    region = "us-east-1"
}

locals {
    
    all_tags = merge(
        tomap({"tag_app"=var.project_name}),
        tomap({"tag_environment"=var.tag_environment})
        )
}

data "aws_vpc" "gtvpc" {
  filter {
    name   = "tag:Name"
    values = ["GT-VPC"]
  }
}


data "aws_subnet_ids" "private_subnet" {
  vpc_id = data.aws_vpc.gtvpc.id

  tags = {
    #Name = "gtreg-sub-pri-aza-ad"
    Name = "gtreg-sub-pri*"
  }

}

data "aws_security_groups" "allowed_ec2" {
  filter {
    name   = "group-name"
    values = ["*app*", "*jenkins*"]
  }
}

resource "aws_rds_cluster_parameter_group" "default" {
  name        = "tag-rds-aurora-mysql-${var.tag_environment}"
  family      = "aurora-mysql5.7"
  description = "tag rds default cluster parameter group"

  parameter {
    name  = "general_log"
    value = "1"
  }

  parameter {
    name  = "lower_case_table_names"
    value = "1"
    apply_method = "pending-reboot"
  }

  parameter {
    name  = "max_allowed_packet"
    value = "524288000"
  }

  parameter {
    name  = "optimizer_switch"
    value = "index_merge=off"
  }

  parameter {
    name  = "performance_schema"
    value = "1"
    apply_method = "pending-reboot"
  }

  parameter {
    name  = "require_secure_transport"
    value = "ON"
  }

  parameter {
    name  = "slow_query_log"
    value = "1"
  }

  parameter {
    name  = "tx_isolation"
    value = "READ-COMMITTED"
  }

}


module "rds" {
  # https://registry.terraform.io/modules/terraform-aws-modules/rds-aurora/aws/latest
  source  = "terraform-aws-modules/rds-aurora/aws"
  version = "~> 5.0"

  name           = "rds-${var.tag_environment}"
  engine         = "aurora-mysql"
  engine_version = "5.7"
  instance_type  = "db.r5.large"

  vpc_id  = "${data.aws_vpc.gtvpc.id}"
  subnets = "${data.aws_subnet_ids.private_subnet.ids}"

  replica_count            = 1
  #replica_scale_enabled    = true
  replica_scale_enabled    = "${var.tag_environment == "prod" ? true : false}"
  replica_scale_min        = 2
  replica_scale_max        = 2
  allowed_security_groups = "${data.aws_security_groups.allowed_ec2.ids}"

  storage_encrypted   = true
  apply_immediately   = true
  monitoring_interval = 10

  db_cluster_parameter_group_name = "tag-rds-aurora-mysql-${var.tag_environment}"

  username = var.tag_rds_username
  create_random_password = false
  password = var.tag_rds_password
  preferred_backup_window = "00:00-03:00"
  preferred_maintenance_window = "thu:03:00-thu:06:00"
  
  #deletion_protection
  enabled_cloudwatch_logs_exports = ["audit", "error", "general", "slowquery"]

  tags = local.all_tags

  depends_on = [
    aws_rds_cluster_parameter_group.default
  ]
}


