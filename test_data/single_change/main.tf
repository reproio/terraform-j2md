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
  name = "test1_changed"
}
