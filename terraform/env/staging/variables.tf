variable "dns-name" {
  type = string
}

variable "dns-zone" {
  type = string
}

variable "public-prefixes" {
  type    = list(string)
  default = []
}

variable "public-root" {
  type    = bool
  default = false
}
