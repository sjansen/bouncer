resource "aws_cloudwatch_event_rule" "rekey" {
  name = "${var.dns-name}-rekey"
  tags = var.tags

  description         = "Rotate JWT keys."
  schedule_expression = "rate(4 hours)"
}

resource "aws_cloudwatch_event_target" "rekey" {
  arn  = module.lambda.function.arn
  rule = aws_cloudwatch_event_rule.rekey.name
}
