terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

data "aws_caller_identity" "current" {}

resource "aws_ecr_repository" "saas_backend" {
  name = "tenet/saas-backend"
}

output "ecr_url" {
  value = aws_ecr_repository.saas_backend.repository_url
}