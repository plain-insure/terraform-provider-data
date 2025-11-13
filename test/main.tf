

terraform {
  required_providers {
    data = {
      source = "plain-insure/data"
    }
  }
  backend "local" {
    path = "terraform.tfstate"
  }
}


locals {
    demo = ""
}

resource "data_notnull" "validation_token2" {
  value         = local.demo
  default_value = "missing"
}


output "validation_token2" {
  value = data_notnull.validation_token2.result
}
