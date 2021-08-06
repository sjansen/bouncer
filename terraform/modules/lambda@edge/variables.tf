variable "handler" {
  type = string
}

variable "memory-size" {
  type    = number
  default = 128
}

variable "name" {
  type = string
}

variable "tags" {
  type = map(string)
}

variable "zip_hash" {
  type = string
}

variable "zip_path" {
  type = string
}
