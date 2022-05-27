output "output_subnet" {
  value = "${tolist(data.aws_subnet_ids.target_subnet.ids)[0]}"
}

output "output_ami" {
  value = "${data.aws_ami.amazon-linux-2.id}"

}

output "output_kms" {
  value = "${data.aws_kms_alias.kms_general.arn}"

}

output "output_local_tag_ec2_config" {
  value = "${local.tag_ec2_config}"

}

output "output_local_tag_ec2_config_map" {
  value = "${local.tag_ec2_config_map}"

}

output "output_launch_template_app" {
  value = "${aws_launch_template.ec2-instance["app:app"].name}"
}

output "output_security_group_app" {
  value = ["${module.security_group["app:app"].security_group_id}"]
}

output "output_vpc_cidr" {

  value = ["${data.aws_vpc.gtvpc.cidr_block}"]
}

output "output_alb_target_group" {

  value = "${concat(module.alb["internal:internal"].target_group_arns, module.alb["external:external"].target_group_arns)}"
}
