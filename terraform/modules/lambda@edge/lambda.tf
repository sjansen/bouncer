resource "aws_lambda_function" "this" {
  provider = aws.us-east-1
  tags     = var.tags

  filename         = var.zip_path
  source_code_hash = var.zip_hash

  function_name = var.name
  handler       = var.handler
  memory_size   = var.memory-size
  publish       = true
  role          = aws_iam_role.this.arn
  runtime       = "nodejs14.x"
}
