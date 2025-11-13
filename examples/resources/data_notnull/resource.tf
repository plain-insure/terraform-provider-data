# Copyright (c) Plain Technologies Aps

terraform {
  required_providers {
    data = {
      source = "plain-insure/data"
    }
  }
}

provider "data" {}

# Example with value provided
resource "data_notnull" "example1" {
  value         = "hello"
  default_value = "default"
}

output "result1" {
  value = data_notnull.example1.result
}

# Example with only default_value (no value)
resource "data_notnull" "example2" {
  default_value = "fallback"
}

output "result2" {
  value = data_notnull.example2.result
}

# Example with value that could change to null
resource "data_notnull" "example3" {
  value         = "initial_value"
  default_value = "default"
}

output "result3" {
  value = data_notnull.example3.result
}
