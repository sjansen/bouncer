data "archive_file" "origin-request" {
  type        = "zip"
  output_path = "${path.module}/origin-request.zip"
  source {
    filename = "index.js"
    content  = file("${path.module}/origin-request.js")
  }
}

data "archive_file" "viewer-request" {
  type        = "zip"
  output_path = "${path.module}/viewer-request.zip"
  source {
    filename = "index.js"
    content  = file("${path.module}/../../../cloudfront/viewer-request/dist.js")
  }
  source {
    filename = "config.js"
    content  = <<EOT
'use strict';
exports.JWKS_ENDPOINT = new URL('https://${var.dns-name}/b/jwks/');
exports.PUBLIC_PREFIXES = new Set(${jsonencode(var.public-prefixes)});
exports.ROOT_IS_PUBLIC = ${jsonencode(var.public-root)};
    EOT
  }
}
