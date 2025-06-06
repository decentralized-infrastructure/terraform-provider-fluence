terraform {
  required_providers {
    fluence = {
      source = "hashicorp.com/0xthresh/fluence"
    }
  }
}

provider "fluence" {
  # Set the FLUENCE_API_KEY environment variable in the terminal to avoid exposing keys
}

data "fluence_ssh_keys" "keys" {}

output "fluence_ssh_keys" {
  value = data.fluence_ssh_keys.keys
}
