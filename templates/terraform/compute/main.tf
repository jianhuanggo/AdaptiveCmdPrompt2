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
      
    tag_ec2_config = flatten([
        for resource in keys(var.tag_ec2_map) : [
            for elem in var.tag_ec2_map[resource] : {
                resource          = resource
                name              = elem.name
                machine_type      = elem.machine_type
                key_pair          = elem.key_pair
                ami               = elem.ami
                subnet_type       = elem.subnet_type
                ingress_cidr_blocks = elem.ingress_cidr_blocks
                ingress_rules     = elem.ingress_rules
            }
            ]
        ])

    tag_ec2_config_map = {
        for s in local.tag_ec2_config: "${s.resource}:${s.name}" => s
    }

    tag_alb_config = flatten([
        for resource in keys(var.tag_alb_map) : [
            for elem in var.tag_alb_map[resource] : {
                resource          = resource
                name              = elem.name
            }
            ]
        ])

    tag_alb_config_map = {
        for s in local.tag_alb_config: "${s.resource}:${s.name}" => s
    }    

}

##################################################################
# Data sources to get VPC, subnet, security group and AMI details
##################################################################
data "aws_vpc" "gtvpc" {
  filter {
    name   = "tag:Name"
    values = ["GT-VPC"]
  }
}

data "aws_subnet_ids" "target_subnet" {
  vpc_id = data.aws_vpc.gtvpc.id

  tags = {
    Name = "gtreg-sub-pub-aza-dmz"
  }

}

data "aws_subnet_ids" "public_subnet" {
  vpc_id = data.aws_vpc.gtvpc.id

  tags = {
    #Name = "gtreg-sub-pub-aza-dmz"
    Name = "gtreg-sub-pub*"
  }

}

data "aws_subnet_ids" "private_subnet" {
  vpc_id = data.aws_vpc.gtvpc.id

  tags = {
    #Name = "gtreg-sub-pri-aza-ad"
    Name = "gtreg-sub-pri*"
  }

}

module "s3_bucket_for_logs" {
  source = "terraform-aws-modules/s3-bucket/aws"

  for_each    = local.tag_alb_config_map
  bucket = "s3-alb-tag-app-${each.value.name}-${var.tag_environment}"
  acl    = "log-delivery-write"

  # Allow deletion of non-empty bucket
  force_destroy = true

  attach_elb_log_delivery_policy = true
}

module "security_group" {
  # https://registry.terraform.io/modules/terraform-aws-modules/security-group/aws/latest
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  for_each    = local.tag_ec2_config_map
  name        = "sg_tag_${var.app_name}_${each.value.name}"
  description = "security group for ${each.value.name} ec2 instance in ${"project_name"} ${var.tag_environment} enviroment"
  vpc_id      = data.aws_vpc.gtvpc.id

  ingress_cidr_blocks = "${each.value.ingress_cidr_blocks == ["default"] ? [data.aws_vpc.gtvpc.cidr_block] : each.value.ingress_cidr_blocks}"
  ingress_rules       = "${each.value.ingress_rules == ["default"] ? ["ssh-tcp"] : each.value.ingress_rules}"
  egress_rules        = ["all-all"]
}

data "aws_ami" "amazon-linux-2" {
 most_recent = true
 owners = ["amazon"]

 filter {
   name   = "owner-alias"
   values = ["amazon"]
 }


 filter {
   name   = "name"
   values = ["amzn2-ami-hvm*"]
 }
}

data "aws_kms_alias" "kms_general" {
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/kms_alias
  # data.aws_kms_alias.kms_general.arn  
  name = "alias/kms-tag-general"
}

data "aws_iam_instance_profile" "kms_general_instance_role" {
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_instance_profile
  # data.aws_iam_instance_profile.kms_general_instance_role.arn
  name = "kms_general_instance_role"
}

resource "aws_launch_template" "ec2-instance" {
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/launch_template#capacity_reservation_preference

  for_each    = local.tag_ec2_config_map
  #name = "lt-${var.tag_usage}-${var.tag_environment}"
  name = "lt-tag-${each.value.name}-${var.tag_environment}"

  description = "launch template for ${each.value.name} ec2 instance in ${"project_name"} ${var.tag_environment} enviroment"

  update_default_version = true

  block_device_mappings {
    device_name = "/dev/xvda"

    ebs {

      volume_type = "gp2"
      delete_on_termination = true
      encrypted = true
      kms_key_id = "${data.aws_kms_alias.kms_general.arn}"
      volume_size = 200

    }
  }

  disable_api_termination = false

  ebs_optimized = true

  iam_instance_profile {
    arn = "${data.aws_iam_instance_profile.kms_general_instance_role.arn}"
  }

  image_id = "${each.value.ami != "default" ? each.value.ami : data.aws_ami.amazon-linux-2.id}"
  instance_initiated_shutdown_behavior = "terminate"
  instance_type = "${each.value.machine_type}"

  key_name = "${each.value.key_pair}"

  monitoring {
    enabled = true
  }

  network_interfaces {
    #associate_public_ip_address = true
    associate_public_ip_address =  "${each.value.name == "bastion" ? true : false}"
    security_groups = ["${module.security_group["${each.value.name}:${each.value.name}"].security_group_id}"]
  }

  tag_specifications {
    resource_type = "instance"

    tags = merge(
      tomap({"Name"="${var.project_name}-${each.value.name}"}),
      local.all_tags
    )

  }

  tag_specifications {
    resource_type = "volume"

    tags = local.all_tags
    
  }
  
  user_data = filebase64("${path.module}/user_data.sh")

  tags = local.all_tags
}

module "alb_security_group_external" {
  # https://registry.terraform.io/modules/terraform-aws-modules/security-group/aws/latest
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"
  
  name        = "sg_tag_alb_external"
  description = "security group for external alb in ${"project_name"} ${var.tag_environment} enviroment"
  vpc_id      = data.aws_vpc.gtvpc.id


  ingress_cidr_blocks = ["0.0.0.0/0"]
  ingress_rules       = ["http-80-tcp", "https-443-tcp"]
  egress_rules        = ["all-all"]
}

module "alb_security_group_internal" {
  # https://registry.terraform.io/modules/terraform-aws-modules/security-group/aws/latest
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name        = "sg_tag_alb_internal"
  description = "security group for internal alb in ${"project_name"} ${var.tag_environment} enviroment"
  vpc_id      = data.aws_vpc.gtvpc.id


  ingress_cidr_blocks = ["172.16.0.0/16", "10.2.0.0/16", "143.215.0.0/16"]
  ingress_rules       = ["http-80-tcp", "https-443-tcp"]
  egress_rules        = ["all-all"]
}

data "aws_acm_certificate" "issued" {
  domain   = "*.tag.gatech.edu"
  statuses = ["ISSUED"]
}



module "alb" {
  # https://registry.terraform.io/modules/terraform-aws-modules/alb/aws/latest
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 6.0"

  for_each    = local.tag_alb_config_map
  name = "alb-tag-app-${each.value.name}-${var.tag_environment}"

  load_balancer_type = "application"

  vpc_id             = "${data.aws_vpc.gtvpc.id}"
  subnets             = "${each.value.name == "internal" ? data.aws_subnet_ids.private_subnet.ids : data.aws_subnet_ids.public_subnet.ids}"
  security_groups    = ["${each.value.name == "internal" ? module.alb_security_group_internal.security_group_id : module.alb_security_group_external.security_group_id}"]
  internal           = "${each.value.name == "internal" ? true : false}"

  access_logs = {
    bucket = "s3-alb-tag-app-${each.value.name}-${var.tag_environment}"
  }

  target_groups = [
    {
      name_prefix      = "tg-"
      backend_protocol = "HTTP"
      backend_port     = 80
      target_type      = "instance"
      targets = []
      #matcher = "200,301,302"
    },
    {
      name_prefix      = "tg-"
      backend_protocol = "HTTPS"
      backend_port     = 443
      target_type      = "instance"
      targets = []
    }
  ]

  https_listeners = [
    {
      port               = 443
      protocol           = "HTTPS"
      certificate_arn    = "${data.aws_acm_certificate.issued.arn}"
      target_group_index = 1
    }
  ]

  http_tcp_listeners = [
    {
      port               = 80
      protocol           = "HTTP"
      target_group_index = 0
    }
  ]

  tags = local.all_tags
  depends_on = [
    module.s3_bucket_for_logs,
    module.alb_security_group_external,
    module.alb_security_group_internal
  ]

}


module "asg" {
  # https://registry.terraform.io/modules/terraform-aws-modules/autoscaling/aws/latest
  source  = "terraform-aws-modules/autoscaling/aws"
  version = "~> 4.0"

  # Autoscaling group
  for_each    = local.tag_ec2_config_map
  name = "asg-tag-${var.app_name}-${each.value.name}-${var.tag_environment}"




  min_size                  = 1
  max_size                  = 1
  desired_capacity          = 1
  wait_for_capacity_timeout = 0
  health_check_type         = "EC2"
  vpc_zone_identifier = "${each.value.subnet_type != "public" ? data.aws_subnet_ids.private_subnet.ids : data.aws_subnet_ids.public_subnet.ids}"

  instance_refresh = {
    strategy = "Rolling"
    preferences = {
      min_healthy_percentage = 50
    }
    triggers = ["tag"]
  }

  # Launch template
  use_lt          = true
  launch_template = "${aws_launch_template.ec2-instance["${each.value.name}:${each.value.name}"].name}"
  default_version = "12"
  target_group_arns = "${each.value.name == "web" ? concat(module.alb["internal:internal"].target_group_arns, module.alb["external:external"].target_group_arns) : []}"
  
  tags_as_map = local.all_tags
  
  depends_on = [
    module.alb,
    aws_launch_template.ec2-instance
  ]
    
  # However, if software running in this EC2 instance needs access
  # to the S3 API in order to boot properly, there is also a "hidden"
  # dependency on the aws_iam_role_policy that Terraform cannot
  # automatically infer, so it must be declared explicitly:
  #depends_on = [
  #  aws_iam_role_policy.example,
  #]

}
