output "name" {
  value = module.lambda.function.function_name
}

output "repo-url" {
  value = module.lambda.repo-url
}
