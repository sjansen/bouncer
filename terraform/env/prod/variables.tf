variable "dns-name" {
  type = string
}

variable "dns-zone" {
  type = string
}

variable "media-updater-groups" {
  type    = list(string)
  default = []
}

variable "media-updater-roles" {
  type    = list(string)
  default = []
}

variable "public-prefixes" {
  type    = list(string)
  default = []
}

variable "public-root" {
  type    = bool
  default = false
}

variable "tags" {
  type    = map(string)
  default = {}
}
