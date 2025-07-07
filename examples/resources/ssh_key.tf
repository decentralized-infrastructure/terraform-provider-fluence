# SSH Key Management Examples

# Create a new SSH key
resource "fluence_ssh_key" "example" {
  name       = "my-terraform-key"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAD0WOgqjaY9EBuhEYg0nTQNuHwGH0Tg/YgtS57VF4g9 user@example.com"
}

# Data source to list all SSH keys
data "fluence_ssh_keys" "all" {}

# Output the created SSH key details
output "ssh_key_details" {
  value = {
    id          = fluence_ssh_key.example.id
    name        = fluence_ssh_key.example.name
    fingerprint = fluence_ssh_key.example.fingerprint
    algorithm   = fluence_ssh_key.example.algorithm
    active      = fluence_ssh_key.example.active
    created_at  = fluence_ssh_key.example.created_at
  }
}

# Output all SSH keys
output "all_ssh_keys" {
  value = data.fluence_ssh_keys.all.ssh_keys
}
