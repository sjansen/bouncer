module "app" {
  source = "../../app"
  tags   = local.tags

  dns-name        = var.dns-name
  dns-zone        = var.dns-zone
  prefix          = local.env
  public-prefixes = var.public-prefixes
  public-root     = var.public-root

  providers = {
    aws            = aws
    aws.cloudfront = aws.cloudfront
    aws.route53    = aws.route53
  }
}
