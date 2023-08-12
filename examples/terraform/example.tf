terraform {
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "3.2.1"
    }
  }
}

provider "null" {}


resource "null_resource" "foo" {}

resource "null_resource" "bar" {}

resource "null_resource" "im_a_count" {
  count = 5
}
