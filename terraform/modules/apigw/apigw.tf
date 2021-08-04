resource "aws_api_gateway_account" "this" {
  cloudwatch_role_arn = aws_iam_role.apigw.arn
}

resource "aws_api_gateway_base_path_mapping" "this" {
  api_id      = aws_api_gateway_rest_api.this.id
  stage_name  = aws_api_gateway_stage.this.stage_name
  domain_name = aws_api_gateway_domain_name.this.domain_name
}

resource "aws_api_gateway_domain_name" "this" {
  domain_name              = var.dns-name
  regional_certificate_arn = aws_acm_certificate_validation.cert-apigw.certificate_arn
  security_policy          = "TLS_1_2"
  tags                     = var.tags

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "this" {
  rest_api_id = join("", aws_api_gateway_rest_api.this.*.id)

  lifecycle {
    create_before_destroy = true
  }
  triggers = {
    redeployment = sha1(jsonencode([
      var.lambda-invoke-arn,
      aws_api_gateway_integration.this.id,
      aws_api_gateway_integration.this.id,
      aws_api_gateway_method.this.id,
      aws_api_gateway_method.this.id,
      aws_api_gateway_resource.this.id,
    ]))
  }

  depends_on = [
    aws_api_gateway_integration.root,
    aws_api_gateway_integration.this,
  ]
}

resource "aws_api_gateway_integration" "root" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_method.root.resource_id
  http_method = aws_api_gateway_method.root.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda-invoke-arn
}

resource "aws_api_gateway_integration" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_method.this.resource_id
  http_method = aws_api_gateway_method.this.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = var.lambda-invoke-arn
}

resource "aws_api_gateway_method" "root" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  resource_id   = aws_api_gateway_rest_api.this.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
  request_parameters = {
    "method.request.header.Host" = true
  }
}

resource "aws_api_gateway_method" "this" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  resource_id   = aws_api_gateway_resource.this.id
  http_method   = "ANY"
  authorization = "NONE"
  request_parameters = {
    "method.request.header.Host" = true
  }
}

resource "aws_api_gateway_method_settings" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  stage_name  = aws_api_gateway_stage.this.stage_name
  method_path = "*/*"

  settings {
    metrics_enabled = true
    logging_level   = "ERROR"
  }
}

resource "aws_api_gateway_resource" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_rest_api.this.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_rest_api" "this" {
  name = var.dns-name
  tags = var.tags

  minimum_compression_size = 65536
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_stage" "this" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  deployment_id = aws_api_gateway_deployment.this.id
  stage_name    = "default"

  xray_tracing_enabled = true
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.apigw.arn
    format          = <<EOT
$context.identity.sourceIp $context.identity.caller $context.identity.user [$context.requestTime] "$context.httpMethod $context.resourcePath $context.protocol" $context.status $context.responseLength $context.requestId
EOT
  }

  depends_on = [aws_cloudwatch_log_group.apigw]
}
