module "rekey-fn" {
  source = "../components/rekey-fn"
  tags   = var.tags

  prefix     = var.prefix
  ssm-prefix = var.dns-name
}
