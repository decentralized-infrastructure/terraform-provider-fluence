# Create multiple SSH keys
resource "fluence_ssh_key" "primary" {
  name       = "primary-key"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAD0WOgqjaY9EBuhEYg0nTQNuHwGH0Tg/YgtS57VF4g9 test@decentralizedinfra.com"
}

# Read all SSH keys (including the ones we just created)
data "fluence_ssh_keys" "all_keys" {
  depends_on = [
    fluence_ssh_key.primary,
  ]
}

# Output the created SSH keys
output "created_ssh_keys" {
  value = {
    primary = {
      id          = fluence_ssh_key.primary.id
      name        = fluence_ssh_key.primary.name
      fingerprint = fluence_ssh_key.primary.fingerprint
      algorithm   = fluence_ssh_key.primary.algorithm
    }
  }
}

# Output all SSH keys from the data source
output "all_ssh_keys" {
  value = data.fluence_ssh_keys.all_keys
}
