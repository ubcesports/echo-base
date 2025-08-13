terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.8.0"
    }
  }

  required_version = ">= 1.2"
}

provider "aws" {
  region = var.aws_region
}

resource "aws_instance" "echo_base" {
  ami                    = "ami-0aff18ec83b712f05"
  instance_type          = var.instance_type
  vpc_security_group_ids = [aws_security_group.echo_base_sg.id]

  tags = {
    "Name" = "echo_base"
  }
}
