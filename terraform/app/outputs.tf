output "docker-registry" {
  value = module.rekey-fn.docker-registry
}

output "web-fn-name" {
  value = module.web-fn.name
}

output "web-fn-repo-url" {
  value = module.web-fn.repo-url
}
