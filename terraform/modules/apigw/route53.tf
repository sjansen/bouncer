data "aws_route53_zone" "zone" {
  provider     = aws.route53
  name         = var.dns-zone
  private_zone = false
}


resource "aws_route53_record" "A" {
  provider = aws.route53
  zone_id  = data.aws_route53_zone.zone.id
  name     = var.dns-name
  type     = "A"

  alias {
    name                   = aws_cloudfront_distribution.this.domain_name
    zone_id                = aws_cloudfront_distribution.this.hosted_zone_id
    evaluate_target_health = false
  }
}


resource "aws_route53_record" "cert-apigw" {
  provider = aws.route53
  for_each = {
    for dvo in aws_acm_certificate.cert-apigw.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.zone.zone_id
}


resource "aws_route53_record" "cert-cloudfront" {
  provider = aws.route53
  for_each = {
    for dvo in aws_acm_certificate.cert-cloudfront.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.zone.zone_id
}
