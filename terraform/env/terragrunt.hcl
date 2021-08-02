locals {
  aws_version        = "~> 3.52"
  docker_version     = "~> 2.14"
  terraform_version  = "~> 1.0.3"
  terragrunt_version = "~> 0.31.3"

  env  = path_relative_to_include()
  proj = "bouncer"
  region = {
    production = "us-west-2"
    staging    = "us-east-2"
  }[local.env]

  config = find_in_parent_folders("terragrunt-local.json", "")
  prefix = local.config == "" ? local.proj : jsondecode(file(local.config)).prefix
}

generate "locals-provider" {
  path      = "locals-generated.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
locals {
  env  = "${local.env}"
  proj = "${local.proj}"
}
EOF
}

generate "providers" {
  path      = "providers-generated.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
provider "aws" {
  region = "${local.region}"
}

provider "aws" {
  alias  = "us-east-1"
  region = "us-east-1"
}

terraform {
  required_version = "${local.terraform_version}"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "${local.aws_version}"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "${local.docker_version}"
    }
  }
}
EOF
}

remote_state {
  backend = "s3"
  config = {
    region         = local.region
    dynamodb_table = "terraform"
    bucket         = "${local.prefix}-terraform-${local.region}"
    key            = "${local.proj}/${local.env}.tfstate"
    encrypt        = true
  }
}

terragrunt_version_constraint = local.terragrunt_version
