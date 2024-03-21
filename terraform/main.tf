provider "aws" {
  region     = "ap-southeast-1"
  access_key = var.access_key
  secret_key = var.secret_key
}

resource "aws_default_vpc" "default" {
  force_destroy = false
  tags = {
    Name = "Default VPC"
  }
}


data "aws_ami" "ec2_ecs_optimised" {
  most_recent = true

  filter {
    name   = "name"
    values = ["amzn2-ami-ecs-kernel-5.10-hvm-2.0.20240221-x86_64-ebs"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["591542846629"] # Amazon
}

resource "aws_instance" "demo-instances" {
  for_each               = toset(["scylla-demo", "cassandra-demo", "loader-demo"])
  ami                    = data.aws_ami.ec2_ecs_optimised.id
  vpc_security_group_ids = [aws_default_vpc.default.default_security_group_id]
  instance_type          = "t2.xlarge"
  key_name               = var.key_name
  user_data              = <<EOF
  #!/bin/bash
  yum update

  # set hostname
  hostnamectl set-hostname ${each.key}

  # setup aio max nr
  echo "fs.aio-max-nr = 1048576" >> /etc/sysctl.conf

  # install pkgs
  yum install curl git go vim tmux unzip -y

  # install cqlsh
  python3 -m pip install cqlsh-expansion

  # install docker compose 
  curl -L https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
  chmod +x /usr/local/bin/docker-compose
  chgrp docker /usr/local/bin/docker-compose
  chmod 750 /usr/local/bin/docker-compose

  # clone repo
  su ec2-user -c 'git clone https://github.com/soonann/scylla-topic-research.git /home/ec2-user/scylla-topic-research'
  EOF

  tags = {
    Name = each.key
  }
}
