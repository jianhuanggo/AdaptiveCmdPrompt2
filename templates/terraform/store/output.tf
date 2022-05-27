output "output_security_group_web_jenkins" {
  #value = "${aws_launch_template.ec2-instance.*}"
  value = ["${data.aws_security_groups.allowed_ec2.ids}"]
}


output "output_rds_master_username" {
  value = "${module.rds.rds_cluster_master_username}"
    sensitive = true
}

output "output_rds_master_password" {
  value = "${module.rds.rds_cluster_master_password}"
  sensitive = true
}

output "output_rds_master_port" {
  value = "${module.rds.rds_cluster_port}"

}

output "output_rds_master_endpoint" {
  value = "${module.rds.rds_cluster_instance_endpoints}"

}

output "output_rds_master_engine_version" {
  value = "${module.rds.rds_cluster_engine_version}"

}