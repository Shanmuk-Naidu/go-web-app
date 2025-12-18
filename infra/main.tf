terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }
  required_version = ">= 1.2.0"
}

provider "aws" {
  region  = "ap-south-1"
}

resource "aws_security_group" "web_sg" {
  name        = "go-web-app-sg-terraform"
  description = "Allow SSH, HTTP, and App traffic"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"]
}

resource "aws_instance" "app_server" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
  
  # Your Key Pair Name
  key_name      = "Ubuntu-kp" 

  vpc_security_group_ids = [aws_security_group.web_sg.id]

  tags = {
    Name = "Go-Web-App-Terraform"
  }

  user_data = <<-EOF
              #!/bin/bash
              # 1. Update and install Docker
              sudo apt-get update -y
              sudo apt-get install -y docker.io
              
              # 2. Start Docker Service
              sudo systemctl start docker
              sudo systemctl enable docker
              sudo usermod -aG docker ubuntu

              # 3. Pull and Run Your App
              # Retries added in case network is slow on startup
              until sudo docker run -d -p 8080:8080 --name go-web-app shanmuk9/go-web-app:latest; do
                echo "Retrying Docker run..."
                sleep 5
              done
              EOF
}

output "instance_ip" {
  value = aws_instance.app_server.public_ip
}