resource "docker_registry_image" "this" {
  name = local.ecr-image-name

  build {
    context  = "${path.module}/docker"
    platform = "linux/amd64"
  }

  lifecycle {
    ignore_changes = all
  }
}
