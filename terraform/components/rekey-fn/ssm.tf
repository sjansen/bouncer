resource "aws_ssm_parameter" "JWK" {
  for_each = toset(["JWK1", "JWK2"])

  name      = "/${var.ssm-prefix}/${each.value}"
  type      = "SecureString"
  value     = "invalid"
  overwrite = false
  tags      = var.tags

  lifecycle {
    ignore_changes = [value]
  }
}
