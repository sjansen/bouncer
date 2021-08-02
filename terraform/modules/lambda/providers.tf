terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    docker = {
      source = "kreuzwerker/docker"
    }
  }
}

provider "docker" {
  registry_auth {
    address  = local.ecr-address
    username = data.aws_ecr_authorization_token.token.user_name
    password = data.aws_ecr_authorization_token.token.password
  }
}
