module "rekey-fn" {
  source = "../components/rekey-fn"
  tags   = var.tags

  dns-name   = var.dns-name
  ssm-prefix = var.dns-name
}

module "web-fn" {
  source = "../components/web-fn"
  tags   = var.tags

  dns-name        = var.dns-name
  dns-zone        = var.dns-zone
  public-prefixes = var.public-prefixes
  ssm-prefix      = var.dns-name

  providers = {
    aws            = aws
    aws.cloudfront = aws.cloudfront
    aws.route53    = aws.route53
  }
}
