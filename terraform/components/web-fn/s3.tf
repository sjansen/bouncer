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
  etag         = filemd5("${path.module}/icons/favicon.ico")
  source       = "${path.module}/icons/favicon.ico"
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
