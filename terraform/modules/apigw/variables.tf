variable "apigw-paths" {
  type    = list(string)
  default = ["/api/*"]
}

variable "cloudwatch-retention" {
  type    = number
  default = 90
}

variable "default-ttl" {
  type    = number
  default = 86400
}

variable "dns-name" {
  description = "for example: bouncer.example.com"
  type        = string
}

variable "dns-zone" {
  description = "for example: example.com"
  type        = string
}

variable "edge-lambdas" {
  description = "event_type => lambda_arn"
  type        = map(string)
  default     = {}
}

variable "lambda-arn" {
  type = string
}

variable "lambda-invoke-arn" {
  type = string
}

variable "logs-bucket" {
  type = string
}

variable "max-ttl" {
  type    = number
  default = 604800
}

variable "media-bucket" {
  type = string
}

variable "tags" {
  type = map(string)
}
