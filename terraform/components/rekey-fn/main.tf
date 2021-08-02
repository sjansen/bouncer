module "lambda" {
  source = "../../modules/lambda"
  tags   = var.tags
  name   = "${var.prefix}-rekey-fn"

  create-ecr-repo = true
  env-vars = {
    BOUNCER_SSM_PREFIX = "/${var.ssm-prefix}/"
  }
}
