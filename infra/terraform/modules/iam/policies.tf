resource "aws_iam_policy" "events_s3" {
  name        = "replay-events-s3-${replace(var.events_bucket_arn, ":", "-")}"
  description = "Scoped S3 read/write for capture event prefixes"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "ListBucket"
        Effect = "Allow"
        Action = ["s3:ListBucket"]
        Resource = [var.events_bucket_arn]
        Condition = {
          StringLike = {
            "s3:prefix" = ["projects/*", "incidents/*"]
          }
        }
      },
      {
        Sid    = "ObjectRW"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ]
        Resource = ["${var.events_bucket_arn}/projects/*", "${var.events_bucket_arn}/incidents/*"]
      }
    ]
  })
}

output "events_s3_policy_arn" {
  value = aws_iam_policy.events_s3.arn
}
