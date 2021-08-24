variable "dns-name" {
  type = string
}

variable "dns-zone" {
  type = string
}

variable "prefix" {
  type = string
}

variable "public-prefixes" {
  type = list(string)
}

variable "tags" {
  type = map(string)
}
