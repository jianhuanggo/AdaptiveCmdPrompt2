output "output_subnet" {
  value = "${tolist(data.aws_subnet_ids.target_subnet.ids)[0]}"
}

output "output_kms" {
  value = "${module.kms_key.key_id}"

}
