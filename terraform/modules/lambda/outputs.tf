output "docker-registry" {
  value = var.create-ecr-repo ? split("/", aws_ecr_repository.this[0].repository_url)[0] : ""
}

output "function" {
  value = aws_lambda_function.this
}

output "repo-url" {
  value = var.create-ecr-repo ? aws_ecr_repository.this[0].repository_url : ""
}

output "role" {
  value = aws_iam_role.this
}
