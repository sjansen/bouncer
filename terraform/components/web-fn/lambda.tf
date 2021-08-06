data "archive_file" "origin-request" {
  type        = "zip"
  output_path = "${path.module}/origin-request.zip"
  source {
    filename = "index.js"
    content  = file("${path.module}/cloudfront.js")
  }
}

data "archive_file" "viewer-request" {
  type        = "zip"
  output_path = "${path.module}/viewer-request.zip"
  source {
    filename = "index.js"
    content  = file("${path.module}/../../../cloudfront/viewer-request/dist.js")
  }
}
