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

# Use existing SSH keys (recommended approach)
data "fluence_ssh_keys" "existing" {}

# Get default OS images (v1.1.0 feature)
data "fluence_default_images" "available" {}

# Select Ubuntu 24.04 image using semantic selection
locals {
  ubuntu_24_image = [for img in data.fluence_default_images.available.images : img if img.distribution == "Ubuntu" && strcontains(img.name, "24.04")][0]
}

# Create a basic VM
resource "fluence_vm" "basic_example" {
  name     = "basic-vm-example"
  hostname = "basic-vm"
  
  # Use Ubuntu 24.04 image with proper download URL
  os_image = local.ubuntu_24_image.download_url
  
  # Use existing SSH key (first one found)
  ssh_keys = [data.fluence_ssh_keys.existing.ssh_keys[0].fingerprint]
  
  # Open essential ports
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
  
  # Basic configuration
  instances           = 1
  basic_configuration = "cpu-4-ram-8gb-storage-25gb"
  
  # Budget constraint
  max_total_price_per_epoch_usd = "5.0"
  
  # Configure timeouts
  timeouts {
    create = "15m"
  }
}

# Output VM information
output "vm_info" {
  value = {
    id        = fluence_vm.basic_example.id
    name      = fluence_vm.basic_example.name
    status    = fluence_vm.basic_example.status
    public_ip = fluence_vm.basic_example.public_ip
    price     = fluence_vm.basic_example.price_per_epoch
  }
  description = "Basic VM details"
}
