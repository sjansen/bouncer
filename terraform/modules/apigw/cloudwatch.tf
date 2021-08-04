resource "aws_cloudwatch_log_group" "apigw" {
  name = "/aws/apigateway/${local.dns-name-underscored}"
  tags = var.tags

  retention_in_days = var.cloudwatch-retention
}
