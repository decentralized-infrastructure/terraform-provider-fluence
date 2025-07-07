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

# Query all available data sources
data "fluence_ssh_keys" "all_keys" {}
data "fluence_vms" "all_vms" {}
data "fluence_available_countries" "countries" {}
data "fluence_available_hardware" "hardware" {}
data "fluence_basic_configurations" "configurations" {}

# Output comprehensive marketplace and resource information
output "marketplace_info" {
  value = {
    available_countries = data.fluence_available_countries.countries.countries
    basic_configurations = data.fluence_basic_configurations.configurations.configurations
    hardware_options = {
      cpu_types = data.fluence_available_hardware.hardware.cpu
      memory_types = data.fluence_available_hardware.hardware.memory
      storage_types = data.fluence_available_hardware.hardware.storage
    }
  }
  description = "Available marketplace options"
}

output "existing_resources" {
  value = {
    ssh_keys = {
      count = length(data.fluence_ssh_keys.all_keys.ssh_keys)
      keys = [for key in data.fluence_ssh_keys.all_keys.ssh_keys : {
        name = key.name
        fingerprint = key.fingerprint
        algorithm = key.algorithm
        created_at = key.created_at
      }]
    }
    vms = {
      count = length(data.fluence_vms.all_vms.vms)
      summary = [for vm in data.fluence_vms.all_vms.vms : {
        id = vm.id
        name = vm.vm_name
        status = vm.status
        public_ip = vm.public_ip
        price_per_epoch = vm.price_per_epoch
        created_at = vm.created_at
      }]
    }
  }
  description = "Existing SSH keys and VMs in your account"
}
