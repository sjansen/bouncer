locals {
  dns-name-underscored = replace(var.dns-name, "/[^-_a-zA-Z0-9]+/", "_")
}
