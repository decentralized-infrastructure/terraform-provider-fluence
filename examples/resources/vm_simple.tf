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

# Create an SSH key for the VM
resource "fluence_ssh_key" "vm_key" {
  name       = "simple-vm-key"
  public_key = file("~/.ssh/id_rsa.pub") # Ensure this file exists and contains your public key
}

# Create a simple VM
resource "fluence_vm" "example" {
  name     = "my-terraform-vm"
  hostname = "terraform-vm"
  os_image = "https://cloud-images.ubuntu.com/releases/24.04/release/ubuntu-24.04-server-cloudimg-amd64.img"
  
  # Use the SSH key we created
  ssh_keys = [fluence_ssh_key.vm_key.name]
  
  # Open common ports
  open_ports = [
    {
      port     = 22
      protocol = "tcp"
    },
    {
      port     = 80
      protocol = "tcp"
    }
  ]
  
  # VM configuration
  instances           = 1
  basic_configuration = "cpu-4-ram-8gb-storage-25gb"
  
  # Budget constraint (optional)
  max_total_price_per_epoch_usd = "5.0"
  
  # Location preference (optional)
  datacenter_countries = ["US", "DE"]
  
  # Configure timeouts (optional)
  timeouts {
    create = "15m"  # Wait up to 15 minutes for VM to become active
  }
  
  # Ensure SSH key is created first
  depends_on = [fluence_ssh_key.vm_key]
}

# Output the VM details
output "vm_details" {
  value = {
    id                = fluence_vm.example.id
    name              = fluence_vm.example.name
    status            = fluence_vm.example.status
    public_ip         = fluence_vm.example.public_ip
    price_per_epoch   = fluence_vm.example.price_per_epoch
    created_at        = fluence_vm.example.created_at
    status_changed_at = fluence_vm.example.status_changed_at
  }
  description = "Details of the created VM"
}

output "ssh_key_details" {
  value = {
    name        = fluence_ssh_key.vm_key.name
    fingerprint = fluence_ssh_key.vm_key.fingerprint
    algorithm   = fluence_ssh_key.vm_key.algorithm
    created_at  = fluence_ssh_key.vm_key.created_at
  }
  description = "Details of the SSH key used by the VM"
}
