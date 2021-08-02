resource "aws_lambda_function" "this" {
  image_uri    = local.ecr-image-name
  package_type = "Image"
  tags         = var.tags

  function_name = var.name
  memory_size   = var.memory-size
  publish       = true
  role          = aws_iam_role.this.arn
  timeout       = var.timeout

  environment {
    variables = var.env-vars
  }

  tracing_config {
    mode = "Active"
  }

  depends_on = [
    aws_cloudwatch_log_group.this,
    docker_registry_image.this,
  ]
}
