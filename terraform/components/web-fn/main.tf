module "apigw" {
  source = "../../modules/apigw"
  tags   = var.tags

  dns-name          = var.dns-name
  dns-zone          = var.dns-zone
  lambda-arn        = module.lambda.function.arn
  lambda-invoke-arn = module.lambda.function.invoke_arn
  logs-bucket       = aws_s3_bucket.logs.bucket_regional_domain_name
  media-bucket      = aws_s3_bucket.media.bucket_regional_domain_name

  edge-lambdas = {
    "viewer-request" : module.viewer-request.function.qualified_arn
  }

  providers = {
    aws           = aws
    aws.us-east-1 = aws.us-east-1
  }
}

module "lambda" {
  source = "../../modules/lambda"
  tags   = var.tags
  name   = "${local.dns-name-underscored}-web-fn"

  create-ecr-repo = true
  env-vars = {
    BOUNCER_APP_URL            = "https://${var.dns-name}/"
    BOUNCER_SAML_CERTIFICATE   = "ssm"
    BOUNCER_SAML_METADATA_URL  = "ssm"
    BOUNCER_SAML_PRIVATE_KEY   = "ssm"
    BOUNCER_SESSION_TABLE_NAME = aws_dynamodb_table.sessions.name
    BOUNCER_SSM_PREFIX         = "/${var.ssm-prefix}/"
  }
}

module "origin-request" {
  source = "../../modules/lambda@edge"
  tags   = var.tags

  handler  = "index.handler"
  name     = "${local.dns-name-underscored}-origin-request"
  zip_hash = data.archive_file.origin-request.output_base64sha256
  zip_path = data.archive_file.origin-request.output_path

  providers = {
    aws           = aws
    aws.us-east-1 = aws.us-east-1
  }
}

module "viewer-request" {
  source = "../../modules/lambda@edge"
  tags   = var.tags

  handler  = "index.handler"
  name     = "${local.dns-name-underscored}-viewer-request"
  zip_hash = data.archive_file.viewer-request.output_base64sha256
  zip_path = data.archive_file.viewer-request.output_path

  providers = {
    aws           = aws
    aws.us-east-1 = aws.us-east-1
  }
}
