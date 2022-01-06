variable "dns-name" {
  type = string
}

variable "dns-zone" {
  type = string
}

variable "public-prefixes" {
  type = list(string)
}

variable "public-root" {
  type = bool
}

variable "media-updater-groups" {
  type = list(string)
}

variable "media-updater-roles" {
  type = list(string)
}

variable "ssm-prefix" {
  type = string
}

variable "tags" {
  type = map(string)
}
