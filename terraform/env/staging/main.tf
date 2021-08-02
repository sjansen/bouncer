module "app" {
  source = "../../app"
  tags   = local.tags

  dns-name = var.dns-name
  prefix   = local.env
}
