resource "aws_cloudfront_distribution" "this" {
  provider = aws.cloudfront

  enabled         = true
  is_ipv6_enabled = true
  price_class     = "PriceClass_100"
  tags            = var.tags

  aliases = [
    var.dns-name
  ]

  dynamic "custom_error_response" {
    for_each = toset([400, 403, 500, 502, 503, 504])
    content {
      error_code            = custom_error_response.value
      error_caching_min_ttl = 0
    }
  }

  custom_error_response {
    error_code            = 404
    error_caching_min_ttl = 60
  }

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD", "OPTIONS"]
    target_origin_id = "s3-bucket"

    compress               = true
    default_ttl            = var.default-ttl
    max_ttl                = var.max-ttl
    min_ttl                = 0
    viewer_protocol_policy = "redirect-to-https"

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    dynamic "lambda_function_association" {
      for_each = var.edge-lambdas
      content {
        event_type = lambda_function_association.key
        lambda_arn = lambda_function_association.value
      }
    }
  }

  logging_config {
    include_cookies = false
    bucket          = var.logs-bucket
  }

  dynamic "ordered_cache_behavior" {
    for_each = toset(var.apigw-paths)
    content {
      path_pattern     = ordered_cache_behavior.value
      allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
      cached_methods   = ["GET", "HEAD", "OPTIONS"]
      target_origin_id = "APIGW"

      compress               = true
      default_ttl            = 0
      max_ttl                = 3600
      min_ttl                = 0
      viewer_protocol_policy = "redirect-to-https"

      forwarded_values {
        query_string = true
        cookies {
          forward = "all"
        }
        headers = [
          "Host",
        ]
      }
    }
  }

  origin {
    domain_name = trimsuffix(trimprefix(aws_api_gateway_stage.this.invoke_url, "https://"), "/default")
    origin_id   = "APIGW"
    custom_origin_config {
      http_port              = "80"
      https_port             = "443"
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  origin {
    domain_name = var.media-bucket
    origin_id   = "s3-bucket"
    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.this.cloudfront_access_identity_path
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate_validation.cert-cloudfront.certificate_arn
    minimum_protocol_version = "TLSv1.2_2021"
    ssl_support_method       = "sni-only"
  }
}

resource "aws_cloudfront_origin_access_identity" "this" {
  provider = aws.cloudfront
  comment  = "OAI for ${var.dns-name}"
}
