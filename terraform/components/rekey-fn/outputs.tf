output "docker-registry" {
  value = module.lambda.docker-registry
}

output "name" {
  value = module.lambda.function.function_name
}

output "repo-url" {
  value = module.lambda.repo-url
}
