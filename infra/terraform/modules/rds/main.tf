variable "environment" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "subnet_ids" {
  type    = list(string)
  default = []
}

variable "db_name" {
  type    = string
  default = "replay"
}

resource "aws_security_group" "postgres" {
  name        = "replay-${var.environment}-postgres"
  description = "Postgres access for replay control plane"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 5432
    to_port     = 5432
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
  }
}

resource "aws_db_subnet_group" "this" {
  name       = "replay-${var.environment}-postgres"
  subnet_ids = var.subnet_ids

  tags = {
    Environment = var.environment
  }
}

resource "aws_db_instance" "postgres" {
  identifier              = "replay-${var.environment}"
  engine                  = "postgres"
  engine_version          = "15"
  instance_class          = "db.t4g.medium"
  allocated_storage       = 50
  db_name                 = var.db_name
  username                = "replay"
  manage_master_user_password = true
  vpc_security_group_ids  = [aws_security_group.postgres.id]
  db_subnet_group_name    = aws_db_subnet_group.this.name
  skip_final_snapshot     = true
  publicly_accessible     = false

  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

output "endpoint" {
  value = aws_db_instance.postgres.address
}

output "security_group_id" {
  value = aws_security_group.postgres.id
}
