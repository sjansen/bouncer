locals {
  tags = merge({
    Environment = local.env
    Project     = local.proj
  }, var.tags)
}
