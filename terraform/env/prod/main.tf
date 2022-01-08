module "app" {
  source = "../../app"
  tags   = local.tags

  dns-name             = var.dns-name
  dns-zone             = var.dns-zone
  media-updater-groups = var.media-updater-groups
  media-updater-roles  = var.media-updater-roles
  prefix               = local.env
  public-prefixes      = var.public-prefixes
  public-root          = var.public-root

  providers = {
    aws            = aws
    aws.cloudfront = aws.cloudfront
    aws.route53    = aws.route53
  }
}
