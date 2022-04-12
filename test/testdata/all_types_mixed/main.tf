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

resource "env_variable" "test1" {
  name = "test1"
}

resource "env_variable" "test2" {
  name = "test2_changed"
}

resource "random_id" "test4" {
  byte_length = 10
}

resource "env_variable" "test5" {
  name = "test5"
}
