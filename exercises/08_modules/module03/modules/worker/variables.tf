variable "worker_name" {
  type        = string
  description = "Unique worker identifier"
}

variable "concurrency" {
  type        = number
  description = "Concurrent job limit"
}
