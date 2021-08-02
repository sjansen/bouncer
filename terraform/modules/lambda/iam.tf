data "aws_iam_policy_document" "AssumeRole-lambda" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "ecr-lambda" {
  statement {
    sid = "LambdaECRImageRetrievalPolicy"
    actions = [
      "ecr:BatchGetImage",
      "ecr:DeleteRepositoryPolicy",
      "ecr:GetDownloadUrlForLayer",
      "ecr:GetRepositoryPolicy",
      "ecr:SetRepositoryPolicy",
    ]
    principals {
      type = "Service"
      identifiers = [
        "lambda.amazonaws.com",
      ]
    }
  }
}

resource "aws_iam_role" "this" {
  name = var.name
  tags = var.tags

  assume_role_policy = data.aws_iam_policy_document.AssumeRole-lambda.json
}

resource "aws_iam_role_policy_attachment" "lambda-xray" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}
