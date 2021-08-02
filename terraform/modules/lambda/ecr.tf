data "aws_ecr_authorization_token" "token" {}

resource "aws_ecr_lifecycle_policy" "this" {
  count = var.create-ecr-repo ? 1 : 0

  repository = aws_ecr_repository.this[0].name
  policy     = <<EOF
{
    "rules": [{
        "rulePriority": 10,
        "description": "Expire untagged images after 3 days",
        "selection": {
            "tagStatus": "untagged",
            "countType": "sinceImagePushed",
            "countUnit": "days",
            "countNumber": 3
        },
        "action": {
            "type": "expire"
        }
    }, {
        "rulePriority": 100,
        "description": "Keep last 3 tagged images",
        "selection": {
            "tagStatus": "any",
            "countType": "imageCountMoreThan",
            "countNumber": 3
        },
        "action": {
            "type": "expire"
        }
    }]
}
EOF
}

resource "aws_ecr_repository" "this" {
  count = var.create-ecr-repo ? 1 : 0
  tags  = var.tags
  name  = var.name

  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository_policy" "this" {
  repository = aws_ecr_repository.this[0].name
  policy     = data.aws_iam_policy_document.ecr-lambda.json
}
