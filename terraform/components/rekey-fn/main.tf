module "lambda" {
  source = "../../modules/lambda"
  tags   = var.tags
  name   = "${local.dns-name-underscored}-rekey-fn"

  create-ecr-repo = true
  env-vars = {
    BOUNCER_SSM_PREFIX = "/${var.ssm-prefix}/"
  }
}
