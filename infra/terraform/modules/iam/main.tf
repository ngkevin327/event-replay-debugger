variable "environment" {
  type = string
}

variable "oidc_provider_arn" {
  type = string
}

variable "events_bucket_arn" {
  type = string
}

locals {
  oidc_issuer = replace(var.oidc_provider_arn, "https://", "")
}

resource "aws_iam_role" "ingestion" {
  name = "replay-${var.environment}-ingestion-irsa"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Federated = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:oidc-provider/${local.oidc_issuer}"
      }
      Action = "sts:AssumeRoleWithWebIdentity"
      Condition = {
        StringEquals = {
          "${local.oidc_issuer}:sub" = "system:serviceaccount:replay:ingestion"
        }
      }
    }]
  })
}

resource "aws_iam_role" "replay_worker" {
  name = "replay-${var.environment}-replay-worker-irsa"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Federated = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:oidc-provider/${local.oidc_issuer}"
      }
      Action = "sts:AssumeRoleWithWebIdentity"
      Condition = {
        StringEquals = {
          "${local.oidc_issuer}:sub" = "system:serviceaccount:replay:replay-worker"
        }
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ingestion" {
  role       = aws_iam_role.ingestion.name
  policy_arn = aws_iam_policy.events_s3.arn
}

resource "aws_iam_role_policy_attachment" "replay_worker" {
  role       = aws_iam_role.replay_worker.name
  policy_arn = aws_iam_policy.events_s3.arn
}

data "aws_caller_identity" "current" {}

output "ingestion_role_arn" {
  value = aws_iam_role.ingestion.arn
}

output "replay_worker_role_arn" {
  value = aws_iam_role.replay_worker.arn
}
