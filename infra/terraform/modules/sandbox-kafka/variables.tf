variable "environment" {
  type        = string
  description = "Deployment environment name"
}

variable "topic_prefix" {
  type        = string
  description = "Prefix for isolated sandbox Kafka topics"
  default     = "replay-sandbox"
}

variable "broker_instance_type" {
  type    = string
  default = "kafka.m5.large"
}
