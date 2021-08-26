output "docker-registry" {
  value = module.rekey-fn.docker-registry
}

output "rekey-fn-name" {
  value = module.rekey-fn.name
}

output "rekey-fn-repo-url" {
  value = module.rekey-fn.repo-url
}

output "web-fn-name" {
  value = module.web-fn.name
}

output "web-fn-repo-url" {
  value = module.web-fn.repo-url
}
