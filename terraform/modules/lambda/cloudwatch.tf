resource "aws_cloudwatch_log_group" "this" {
  name = "/aws/lambda/${var.name}"
  tags = var.tags

  retention_in_days = var.cloudwatch-retention
}
