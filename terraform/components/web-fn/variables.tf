variable "dns-name" {
  type = string
}

variable "dns-zone" {
  type = string
}

variable "ssm-prefix" {
  type = string
}

variable "tags" {
  type = map(string)
}
