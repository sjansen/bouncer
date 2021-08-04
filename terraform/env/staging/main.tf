module "app" {
  source = "../../app"
  tags   = local.tags

  dns-name = var.dns-name
  dns-zone = var.dns-zone
  prefix   = local.env

  providers = {
    aws           = aws
    aws.us-east-1 = aws.us-east-1
  }
}
