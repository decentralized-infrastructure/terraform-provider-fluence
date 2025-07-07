terraform {
  required_providers {
    fluence = {
      source = "hashicorp.com/decentralized-infrastructure/fluence"
      version = "~> 1.0"
    }
  }
}

provider "fluence" {
  # Configuration options:
  # host     = "https://api.fluence.dev"  # Optional, defaults to this value
  # api_key  = "your-api-key"             # Or set FLUENCE_API_KEY env var
}

# Create an SSH key
resource "fluence_ssh_key" "example" {
  name       = "my-terraform-key"
  # Replace with your actual public key associated with a private key you control
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKgJIjnDg1DjqOOxINs78oU3f7PJXIyq9uiNocNVhXNx user@example.com"
}

# Output the SSH key details
output "ssh_key_details" {
  value = {
    name        = fluence_ssh_key.example.name
    fingerprint = fluence_ssh_key.example.fingerprint
    algorithm   = fluence_ssh_key.example.algorithm
    active      = fluence_ssh_key.example.active
    created_at  = fluence_ssh_key.example.created_at
  }
  description = "Details of the created SSH key"
}
