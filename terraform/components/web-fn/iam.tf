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

data "aws_iam_policy_document" "update-media" {
  statement {
    actions = [
      "s3:DeleteObject",
      "s3:GetObject",
      "s3:ListBucket",
      "s3:PutObject",
    ]
    resources = [
      aws_s3_bucket.media.arn,
      "${aws_s3_bucket.media.arn}/*",
    ]
  }
}

resource "aws_iam_group_policy_attachment" "update-media" {
  for_each = toset(var.media-updater-groups)

  group      = each.value
  policy_arn = aws_iam_policy.update-media.arn
}

resource "aws_iam_policy" "update-media" {
  name        = "update-${aws_s3_bucket.media.id}"
  description = "Upload and remove objects from ${aws_s3_bucket.media.id}"
  policy      = data.aws_iam_policy_document.update-media.json
}

resource "aws_iam_role_policy" "this" {
  name   = "all-the-things"
  role   = module.lambda.role.name
  policy = data.aws_iam_policy_document.this.json
}

resource "aws_iam_role_policy_attachment" "update-media" {
  for_each = toset(var.media-updater-roles)

  role       = each.value
  policy_arn = aws_iam_policy.update-media.arn
}