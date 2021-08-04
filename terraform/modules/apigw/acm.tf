resource "aws_acm_certificate" "cert-apigw" {
  domain_name       = var.dns-name
  validation_method = "DNS"
  tags              = var.tags
}


resource "aws_acm_certificate" "cert-cloudfront" {
  provider          = aws.us-east-1
  domain_name       = var.dns-name
  validation_method = "DNS"
  tags              = var.tags
}


resource "aws_acm_certificate_validation" "cert-apigw" {
  certificate_arn         = aws_acm_certificate.cert-apigw.arn
  validation_record_fqdns = [for record in aws_route53_record.cert-apigw : record.fqdn]
}


resource "aws_acm_certificate_validation" "cert-cloudfront" {
  provider                = aws.us-east-1
  certificate_arn         = aws_acm_certificate.cert-cloudfront.arn
  validation_record_fqdns = [for record in aws_route53_record.cert-cloudfront : record.fqdn]
}
