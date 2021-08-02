locals {
  ecr-address    = format("%v.dkr.ecr.%v.amazonaws.com", data.aws_caller_identity.this.account_id, data.aws_region.current.name)
  ecr-repo       = var.create-ecr-repo ? aws_ecr_repository.this[0].id : var.ecr-repo
  ecr-image-name = format("%v/%v:latest", local.ecr-address, local.ecr-repo)
}
