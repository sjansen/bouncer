variable "dns-name" {
  type = string
}

variable "dns-zone" {
  type = string
}

variable "public-prefixes" {
  type = list(string)
}

variable "ssm-prefix" {
  type = string
}

variable "tags" {
  type = map(string)
}
