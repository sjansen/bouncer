variable "cloudwatch-retention" {
  type    = number
  default = 90
}

variable "create-ecr-repo" {
  description = "Controls whether ECR repository for Lambda image should be created"
  type        = bool
  default     = false
}

variable "ecr-repo" {
  description = "Name of ECR repository to use or to create"
  type        = string
  default     = null
}

variable "env-vars" {
  type = map(string)
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

variable "timeout" {
  type    = number
  default = 15
}
