locals {
  proj = "bouncer"

  aws_version        = "~> 3.70"
  docker_version     = "~> 2.15"
  terraform_version  = "~> 1.1.2"
  terragrunt_version = "~> 0.35.16"

  defaults = {
    prefix    = local.proj
    providers = {}
    region = {
      production = "us-west-2"
      staging    = "us-east-2"
    }
  }

  env     = path_relative_to_include()
  found   = find_in_parent_folders("terragrunt-local.json", "")
  encoded = local.found == "" ? jsonencode(local.defaults) : file(local.found)

  decoded = jsondecode(local.encoded)
  prefix  = (
    contains(keys(local.decoded), "prefix")
    ? local.decoded.prefix
    : local.defaults.prefix
  )
  providers  = (
    contains(keys(local.decoded), "providers")
    ? local.decoded.providers
    : local.defaults.providers
  )
  region  = (
    contains(keys(local.decoded), "region")
    ? local.decoded.region[local.env]
    : local.defaults.region[local.env]
  )
  state = {
    region = lookup(
      lookup(
        lookup(local.decoded, "state", {}),
        local.env, {}
      ),
      "region",
      local.region
    )
    bucket = lookup(
      lookup(
        lookup(local.decoded, "state", {}),
        local.env, {}
      ),
      "bucket",
      "${local.prefix}-terraform-${local.region}"
    )
    key = lookup(
      lookup(
        lookup(local.decoded, "state", {}),
        local.env, {}
      ),
      "key",
      "${local.proj}/${local.env}.tfstate"
    )
    table = lookup(
      lookup(
        lookup(local.decoded, "state", {}),
        local.env, {}
      ),
      "table",
      "terraform"
    )
  }
}

generate "locals" {
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
  alias  = "cloudfront"
  region = "us-east-1"
}

provider "aws" {
  alias  = "route53"
  region = "us-east-1"
%{ if contains(keys(local.providers), "route53") ~}
%{ if contains(keys(local.providers.route53), "profile") }
  profile = "${local.providers.route53.profile}"
%{ endif ~}
%{ if contains(keys(local.providers.route53), "role_arn") }
  assume_role {
    role_arn = "${local.providers.route53.role_arn}"
  }
%{ endif ~}
%{ endif }
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
    region         = local.state.region
    dynamodb_table = local.state.table
    bucket         = local.state.bucket
    key            = local.state.key
    encrypt        = true
  }
}

terragrunt_version_constraint = local.terragrunt_version
