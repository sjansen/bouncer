variable "cloudwatch-retention" {
  type    = number
  default = 90
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

variable "media-bucket" {
  type = string
}

variable "tags" {
  type = map(string)
}
