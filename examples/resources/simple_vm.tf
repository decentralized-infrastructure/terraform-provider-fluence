# First, read existing SSH keys to use their fingerprints
data "fluence_ssh_keys" "existing" {}

# Create a simple VM using an existing SSH key
resource "fluence_vm" "test2" {
  name     = "test-vm-2"
  os_image = "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img"
  
  # Use the first available SSH key
  ssh_keys = length(data.fluence_ssh_keys.existing.ssh_keys) > 0 ? [
    data.fluence_ssh_keys.existing.ssh_keys[0].fingerprint
  ] : []
  
  # Basic port configuration
  open_ports = [
    {
      port     = 22
      protocol = "tcp"
    }
  ]
  
  instances = 1
}

# Output VM details
output "vm_info" {
  value = {
    id     = fluence_vm.test2.id
    name   = fluence_vm.test2.name
    status = fluence_vm.test2.status
  }
}

output "available_ssh_keys" {
  value = data.fluence_ssh_keys.existing.ssh_keys
}
