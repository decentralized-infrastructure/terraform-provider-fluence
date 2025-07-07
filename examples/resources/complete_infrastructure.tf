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

# Get available marketplace data to make informed decisions
data "fluence_available_countries" "all" {}
data "fluence_available_hardware" "all" {}
data "fluence_basic_configurations" "all" {}

# List existing VMs
data "fluence_vms" "existing" {}

# List existing SSH keys
data "fluence_ssh_keys" "existing" {}

# Create an SSH key for the VM
resource "fluence_ssh_key" "vm_key" {
  name       = "terraform-infrastructure-key"
  # Replace with your actual public key associated with a private key you control
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKgJIjnDg1DjqOOxINs78oU3f7PJXIyq9uiNocNVhXNx user@example.com"
}

# Create a VM using the SSH key and marketplace data
resource "fluence_vm" "web_server" {
  name     = "terraform-web-server"
  hostname = "web-01"
  os_image = "https://cloud-images.ubuntu.com/releases/24.04/release/ubuntu-24.04-server-cloudimg-amd64.img"
  
  # Reference the SSH key we created
  ssh_keys = [fluence_ssh_key.vm_key.name]
  
  # Open ports for web server
  open_ports = [
    {
      port     = 22
      protocol = "tcp"
    },
    {
      port     = 80
      protocol = "tcp"
    },
    {
      port     = 443
      protocol = "tcp"
    }
  ]
  
  # VM configuration - use available basic configurations
  instances           = 1
  basic_configuration = data.fluence_basic_configurations.all.configurations[1] # cpu-4-ram-8gb-storage-25gb
  
  # Budget and location constraints using available data
  max_total_price_per_epoch_usd = "10.0"
  datacenter_countries          = slice(data.fluence_available_countries.all.countries, 0, 2) # First 2 countries
  
  # Ensure SSH key is created first
  depends_on = [fluence_ssh_key.vm_key]
}

# Outputs showing the complete infrastructure
output "infrastructure_summary" {
  value = {
    ssh_key = {
      name        = fluence_ssh_key.vm_key.name
      fingerprint = fluence_ssh_key.vm_key.fingerprint
      created_at  = fluence_ssh_key.vm_key.created_at
    }
    vm = {
      id                = fluence_vm.web_server.id
      name              = fluence_vm.web_server.name
      status            = fluence_vm.web_server.status
      public_ip         = fluence_vm.web_server.public_ip
      price_per_epoch   = fluence_vm.web_server.price_per_epoch
      created_at        = fluence_vm.web_server.created_at
      status_changed_at = fluence_vm.web_server.status_changed_at
    }
  }
  description = "Complete infrastructure details"
}

output "marketplace_data" {
  value = {
    available_countries      = data.fluence_available_countries.all.countries
    available_configurations = data.fluence_basic_configurations.all.configurations
    hardware_options = {
      cpu_types     = data.fluence_available_hardware.all.cpu
      memory_types  = data.fluence_available_hardware.all.memory
      storage_types = data.fluence_available_hardware.all.storage
    }
  }
  description = "Available marketplace options"
}

output "existing_resources" {
  value = {
    existing_vms     = length(data.fluence_vms.existing.vms)
    existing_ssh_keys = length(data.fluence_ssh_keys.existing.ssh_keys)
  }
  description = "Count of existing resources"
}
