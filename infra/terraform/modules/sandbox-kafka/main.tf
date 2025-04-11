terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

resource "aws_msk_cluster" "sandbox" {
  cluster_name           = "${var.topic_prefix}-${var.environment}"
  kafka_version          = "3.5.1"
  number_of_broker_nodes = 2

  broker_node_group_info {
    instance_type   = var.broker_instance_type
    client_subnets  = []
    security_groups = []
    storage_info {
      ebs_storage_info {
        volume_size = 100
      }
    }
  }
}

locals {
  sandbox_topics = [
    "${var.topic_prefix}.payments",
    "${var.topic_prefix}.ledger",
    "${var.topic_prefix}.notifications",
  ]
}

output "cluster_arn" {
  value = aws_msk_cluster.sandbox.arn
}

output "topic_names" {
  value = local.sandbox_topics
}
