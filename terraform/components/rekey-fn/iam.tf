data "aws_iam_policy_document" "this" {
  statement {
    actions = [
      "ssm:DescribeParameters",
      "ssm:PutParameters",
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
