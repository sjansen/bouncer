data "aws_iam_policy_document" "this" {
  statement {
    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:BatchWriteItem",
      "dynamodb:DeleteItem",
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:Query",
      "dynamodb:UpdateItem",
    ]
    resources = [aws_dynamodb_table.sessions.arn]
  }
  statement {
    actions = [
      "ssm:GetParameter",
      "ssm:GetParameters",
    ]
    resources = [
      "arn:aws:ssm:*:*:parameter/${var.ssm-prefix}/*",
    ]
  }
}

resource "aws_iam_role_policy" "this" {
  name   = "all-the-things"
  role   = module.lambda.role.name
  policy = data.aws_iam_policy_document.this.json
}
