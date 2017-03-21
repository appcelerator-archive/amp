# Outputs

output "public-ip" {
  value = "${aws_instance.m1.public_ip}"
}
