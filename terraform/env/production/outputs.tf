output "docker-registry" {
  value = module.app.docker-registry
}

output "rekey-fn-name" {
  value = module.app.rekey-fn-name
}

output "rekey-fn-repo-url" {
  value = module.app.rekey-fn-repo-url
}

output "web-fn-name" {
  value = module.app.web-fn-name
}

output "web-fn-repo-url" {
  value = module.app.web-fn-repo-url
}
