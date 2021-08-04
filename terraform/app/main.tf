module "rekey-fn" {
  source = "../components/rekey-fn"
  tags   = var.tags

  dns-name   = var.dns-name
  ssm-prefix = var.dns-name
}

module "web-fn" {
  source = "../components/web-fn"
  tags   = var.tags

  dns-name   = var.dns-name
  dns-zone   = var.dns-zone
  ssm-prefix = var.dns-name

  providers = {
    aws           = aws
    aws.us-east-1 = aws.us-east-1
  }
}
