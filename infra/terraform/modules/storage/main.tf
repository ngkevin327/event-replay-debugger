variable "environment" {
  type = string
}

variable "events_bucket_name" {
  type    = string
  default = null
}

locals {
  bucket_name = coalesce(var.events_bucket_name, "replay-${var.environment}-events")
}

resource "aws_s3_bucket" "events" {
  bucket = local.bucket_name

  tags = {
    Environment = var.environment
    Purpose     = "capture-events"
    ManagedBy   = "terraform"
  }
}

resource "aws_s3_bucket_versioning" "events" {
  bucket = aws_s3_bucket.events.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "events" {
  bucket = aws_s3_bucket.events.id

  rule {
    id     = "expire-old-captures"
    status = "Enabled"

    expiration {
      days = 90
    }

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "events" {
  bucket = aws_s3_bucket.events.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

output "events_bucket_name" {
  value = aws_s3_bucket.events.bucket
}

output "events_bucket_arn" {
  value = aws_s3_bucket.events.arn
}
