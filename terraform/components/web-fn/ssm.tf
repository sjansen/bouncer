resource "aws_ssm_parameter" "SAML" {
  for_each = toset(["SAML_CERTIFICATE", "SAML_METADATA_URL"])

  name      = "/${var.ssm-prefix}/${each.value}"
  type      = "String"
  value     = "invalid"
  overwrite = false
  tags      = var.tags

  lifecycle {
    ignore_changes = [value]
  }
}

resource "aws_ssm_parameter" "SAML_PRIVATE_KEY" {
  name      = "/${var.ssm-prefix}/SAML_PRIVATE_KEY"
  type      = "SecureString"
  value     = "invalid"
  overwrite = false
  tags      = var.tags

  lifecycle {
    ignore_changes = [value]
  }
}
