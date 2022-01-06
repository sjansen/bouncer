resource "aws_s3_bucket" "logs" {
  bucket = "${var.dns-name}-logs"
  tags   = var.tags

  acl = "log-delivery-write"
  lifecycle_rule {
    id                                     = "cleanup"
    enabled                                = true
    abort_incomplete_multipart_upload_days = 3
    expiration {
      days = 90
    }
    noncurrent_version_expiration {
      days = 30
    }
  }
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket" "media" {
  bucket = "${var.dns-name}-media"
  tags   = var.tags

  acl = "private"
  lifecycle_rule {
    id                                     = "cleanup"
    enabled                                = true
    abort_incomplete_multipart_upload_days = 3
    noncurrent_version_expiration {
      days = 30
    }
  }
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket_object" "favicon" {
  bucket       = aws_s3_bucket.media.id
  key          = "favicon.ico"
  content_type = "image/x-icon"
  source       = "${path.module}/files/favicon.ico"
  source_hash  = filemd5("${path.module}/files/favicon.ico")
  lifecycle {
    ignore_changes = [source_hash, version_id]
  }
}

resource "aws_s3_bucket_object" "index" {
  bucket       = aws_s3_bucket.media.id
  key          = "index.html"
  content_type = "text/html"
  source       = "${path.module}/files/index.html"
  source_hash  = filemd5("${path.module}/files/index.html")
  lifecycle {
    ignore_changes = [source_hash, version_id]
  }
}

resource "aws_s3_bucket_object" "readme" {
  bucket       = aws_s3_bucket.media.id
  key          = "/b/README"
  content_type = "text/plain; charset=utf-8"
  source       = "${path.module}/files/README"
  source_hash  = filemd5("${path.module}/files/README")
}

resource "aws_s3_bucket_policy" "media" {
  bucket = aws_s3_bucket.media.id
  policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [{
    "Effect":"Allow",
    "Action": "s3:GetObject",
    "Principal": {
      "AWS": "${module.apigw.cf_oai_arn}"
    },
    "Resource": [
      "${aws_s3_bucket.media.arn}/*"
    ]
  }]
}
EOF
}

resource "aws_s3_bucket_public_access_block" "logs" {
  bucket = aws_s3_bucket.logs.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_public_access_block" "media" {
  bucket = aws_s3_bucket.media.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}