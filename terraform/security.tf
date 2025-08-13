resource "aws_security_group" "echo_base_sg" {
  name        = "echo_base_sg"
  description = "Primary security group for Echo Base, the main VPS for UBCEA."

  tags = {
    Name = "echo_base_sg"
  }

  lifecycle {
    ignore_changes = [
      ingress,
      egress
    ]
  }
}

