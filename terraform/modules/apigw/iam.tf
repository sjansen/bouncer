data "aws_iam_policy_document" "AssumeRole-apigw" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["apigateway.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "apigw" {
  name = "${var.dns-name}-APIGateway"
  tags = var.tags

  assume_role_policy = data.aws_iam_policy_document.AssumeRole-apigw.json
}

resource "aws_iam_role_policy_attachment" "apigw" {
  role       = aws_iam_role.apigw.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}
