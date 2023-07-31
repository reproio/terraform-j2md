terraform {
  required_providers {
    env = {
      source = "tchupp/env"
      version = "0.0.2"
    }
  }
}

provider "env" {
  # Configuration options
}

resource "env_variable" "test" {
  name = random_id.test2.hex
}

#resource "random_id" "test" {
#  byte_length = 10
#}

resource "random_id" "test2" {
  byte_length = 10
}
