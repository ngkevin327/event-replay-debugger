terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "environment" {
  type    = string
  default = "staging"
}

module "vpc" {
  source      = "../../modules/vpc"
  environment = var.environment
}

module "eks" {
  source      = "../../modules/eks"
  environment = var.environment
  vpc_id      = module.vpc.vpc_id
  subnet_ids  = module.vpc.private_subnet_ids
}

module "rds" {
  source      = "../../modules/rds"
  environment = var.environment
  vpc_id      = module.vpc.vpc_id
  subnet_ids  = module.vpc.private_subnet_ids
}

output "vpc_id" {
  value = module.vpc.vpc_id
}

output "eks_cluster_name" {
  value = module.eks.cluster_name
}

module "storage" {
  source      = "../../modules/storage"
  environment = var.environment
}

module "clickhouse" {
  source      = "../../modules/clickhouse"
  environment = var.environment
  vpc_id      = module.vpc.vpc_id
  subnet_ids  = module.vpc.private_subnet_ids
}

output "rds_endpoint" {
  value     = module.rds.endpoint
  sensitive = true
}

output "events_bucket" {
  value = module.storage.events_bucket_name
}

output "clickhouse_endpoint" {
  value = module.clickhouse.http_endpoint
}
