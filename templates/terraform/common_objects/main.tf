terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = ">= 3.30"
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

data "aws_iam_policy_document" "default" {
    version = "2012-10-17"
    statement {
      sid    = "Enable IAM User Permissions one"
      effect = "Allow"
      principals {
        type        = "AWS"
        identifiers = ["arn:aws:iam::${var.account_number}:root"]
      }
      actions   = ["kms:*"]
      resources = ["*"]
    }
    statement {
      sid    = "Allow access for Key Administrators"
      effect = "Allow"
      principals {
        type        = "AWS"
        identifiers = ["arn:aws:iam::${var.account_number}:user/BreakGlass"]
      }
      actions = ["kms:Create*", "kms:Describe*", "kms:Enable*", "kms:List*", "kms:Put*", "kms:Update*", "kms:Revoke*", "kms:Disable*", "kms:Get*", "kms:Delete*", "kms:TagResource", "kms:UntagResource", "kms:ScheduleKeyDeletion", "kms:CancelKeyDeletion"]
      resources = ["*"]
    }
    statement {
      sid    = "Allow use of the key"
      effect = "Allow"
      principals {
        type        = "AWS"
        identifiers = ["arn:aws:iam::${var.account_number}:user/BreakGlass"]
      }
      actions = ["kms:Encrypt", "kms:Decrypt", "kms:ReEncrypt*", "kms:GenerateDataKey*", "kms:DescribeKey"]
      resources = ["*"]
    }
    statement {
      sid    = "Allow attachment of persistent resources one"
      effect = "Allow"
      principals {
        type        = "AWS"
        identifiers = ["arn:aws:iam::${var.account_number}:user/BreakGlass"]
      }
      actions = ["kms:CreateGrant", "kms:ListGrants", "kms:RevokeGrant"]
      resources = ["*"]
      condition {
        test     = "Bool"
        variable = "kms:GrantIsForAWSResource"
        values   = ["true"]
      }
    }
    statement {
      sid    = "Allow service-linked role use of the CMK one"
      effect = "Allow"
      principals {
        type        = "AWS"
        identifiers = ["arn:aws:iam::${var.account_number}:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling"]
      }
      actions = ["kms:Encrypt", "kms:Decrypt", "kms:ReEncrypt*", "kms:GenerateDataKey*", "kms:DescribeKey"]
      resources = ["*"]
    }
    statement {
      sid    = "Allow attachment of persistent resources two"
      effect = "Allow"
      principals {
        type        = "AWS"
        identifiers = ["arn:aws:iam::${var.account_number}:role/aws-service-role/autoscaling.amazonaws.com/AWSServiceRoleForAutoScaling"]
      }
      actions = ["kms:CreateGrant"]
      resources = ["*"]
      condition {
        test     = "Bool"
        variable = "kms:GrantIsForAWSResource"
        values   = ["true"]
      }
    }
    statement {
      sid    = "Allow access for awscloudwatchlog"
      effect = "Allow"
      principals {
        type        = "Service"
        identifiers = ["logs.us-east-1.amazonaws.com"]
      }
      actions = ["kms:Encrypt", "kms:Decrypt", "kms:ReEncrypt*", "kms:GenerateDataKey*", "kms:DescribeKey"]
      resources = ["*"]
      condition {
        test     = "ArnLike"
        variable = "kms:EncryptionContext:aws:logs:arn"
        values   = ["arn:aws:logs:us-east-1:${var.account_number}:*"]
      }
    }
}

module "kms_key" {
    # https://registry.terraform.io/modules/clouddrove/kms/aws/latest
    source      = "clouddrove/kms/aws"
    version     = "0.14.0"
    name        = "kms"
    environment = "${var.tag_environment}"
    label_order = ["name", "environment"]
    enabled     = true
    description             = "KMS key for general purpose"
    deletion_window_in_days = 7
    enable_key_rotation     = true
    alias                   = "alias/kms-tag-general"
    policy                  = data.aws_iam_policy_document.default.json
  }

data "aws_iam_policy_document" "kms_tag_general" {
    version = "2012-10-17"
    statement {
      sid    = "VisualEditor0"
      effect = "Allow"
      actions   = [
                "kms:Decrypt",
                "kms:Encrypt",
                "kms:GenerateDataKey",
                "kms:ReEncryptTo",
                "kms:GenerateDataKeyWithoutPlaintext",
                "kms:DescribeKey",
                "kms:GenerateDataKeyPairWithoutPlaintext",
                "kms:GenerateDataKeyPair",
                "kms:CreateGrant",
                "kms:ReEncryptFrom"
            ]
      resources = ["arn:aws:kms:*:${var.account_number}:key/${module.kms_key.key_id}"]
    }
}

resource "aws_iam_policy" "kms_general" {
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy
  name        = "iam_policy_kms_general"
  path        = "/"
  description = "Allows encryption and decryption to kms key kms-general"

  # Terraform's "jsonencode" function converts a
  # Terraform expression result to valid JSON syntax.
  policy = data.aws_iam_policy_document.kms_tag_general.json
}

resource "aws_iam_role" "kms_general_role" {
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role
  name               = "iam_role_kms_general"
  path               = "/system/"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })

  tags = local.all_tags
}

resource "aws_iam_role_policy_attachment" "tag-policy-attach" {
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment
  role       = aws_iam_role.kms_general_role.name
  policy_arn = aws_iam_policy.kms_general.arn
}

resource "aws_iam_instance_profile" "kms_general_instance_role" {
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_instance_profile
  name = "kms_general_instance_role"
  role = aws_iam_role.kms_general_role.name
}


