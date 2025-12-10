variable "environment" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

resource "aws_security_group" "clickhouse" {
  name        = "replay-${var.environment}-clickhouse"
  description = "ClickHouse analytics for replay staging"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 8123
    to_port     = 8123
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/8"]
  }

  ingress {
    from_port   = 9000
    to_port     = 9000
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/8"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Environment = var.environment
    Service     = "clickhouse"
  }
}

resource "aws_instance" "clickhouse" {
  count         = length(var.subnet_ids) > 0 ? 1 : 0
  ami           = data.aws_ami.amazon_linux.id
  instance_type = "m6i.large"
  subnet_id     = var.subnet_ids[0]
  vpc_security_group_ids = [aws_security_group.clickhouse.id]

  user_data = <<-EOF
    #!/bin/bash
    echo "ClickHouse staging placeholder for ${var.environment}"
  EOF

  tags = {
    Name        = "replay-${var.environment}-clickhouse"
    Environment = var.environment
  }
}

data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}

output "security_group_id" {
  value = aws_security_group.clickhouse.id
}

output "http_endpoint" {
  value = length(aws_instance.clickhouse) > 0 ? "http://${aws_instance.clickhouse[0].private_ip}:8123" : ""
}
