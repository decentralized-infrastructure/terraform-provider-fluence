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

# Query all available data sources including new v1.1.0 features
data "fluence_ssh_keys" "all_keys" {}
data "fluence_vms" "all_vms" {}
data "fluence_available_countries" "countries" {}
data "fluence_available_hardware" "hardware" {}
data "fluence_basic_configurations" "configurations" {}

# New v1.1.0 data sources
data "fluence_datacenters" "datacenters" {}
data "fluence_default_images" "images" {}

# Estimate deposit for different VM configurations
data "fluence_vm_estimate_deposit" "small_vm" {
  instances = 1
  
  constraints = {
    basic_configuration           = "cpu-2-ram-4gb-storage-25gb"
    max_total_price_per_epoch_usd = "3.0"
    datacenter_countries          = ["US"]
  }
}

data "fluence_vm_estimate_deposit" "large_cluster" {
  instances = 3
  
  constraints = {
    basic_configuration           = "cpu-8-ram-16gb-storage-50gb"
    max_total_price_per_epoch_usd = "50.0"
    datacenter_countries          = ["US", "DE", "CA"]
  }
}

# Output all available marketplace information
output "ssh_keys_summary" {
  value = {
    total_keys = length(data.fluence_ssh_keys.all_keys.ssh_keys)
    active_keys = length([for key in data.fluence_ssh_keys.all_keys.ssh_keys : key if key.active])
    keys = [for key in data.fluence_ssh_keys.all_keys.ssh_keys : {
      name        = key.name
      fingerprint = key.fingerprint
      algorithm   = key.algorithm
      active      = key.active
    }]
  }
  description = "Summary of all SSH keys"
}

output "vms_summary" {
  value = {
    total_vms = length(data.fluence_vms.all_vms.vms)
    active_vms = length([for vm in data.fluence_vms.all_vms.vms : vm if vm.status == "Active"])
    vms = [for vm in data.fluence_vms.all_vms.vms : {
      id        = vm.id
      name      = vm.name
      status    = vm.status
      public_ip = vm.public_ip
      price     = vm.price_per_epoch
    }]
  }
  description = "Summary of all VMs"
}

output "marketplace_data" {
  value = {
    countries = {
      total = length(data.fluence_available_countries.countries.countries)
      list  = data.fluence_available_countries.countries.countries
    }
    
    configurations = {
      total = length(data.fluence_basic_configurations.configurations.configurations)
      list  = [for config in data.fluence_basic_configurations.configurations.configurations : {
        name        = config.name
        cpu         = config.cpu
        memory_gb   = config.memory_gb
        storage_gb  = config.storage_gb
      }]
    }
    
    hardware = {
      cpu_architectures = data.fluence_available_hardware.hardware.cpu_architectures
      cpu_manufacturers = data.fluence_available_hardware.hardware.cpu_manufacturers
      memory_types      = data.fluence_available_hardware.hardware.memory_types
      memory_generations = data.fluence_available_hardware.hardware.memory_generations
      storage_types     = data.fluence_available_hardware.hardware.storage_types
    }
    
    datacenters = {
      total = length(data.fluence_datacenters.datacenters.datacenters)
      by_country = {for dc in data.fluence_datacenters.datacenters.datacenters : dc.country => dc.name...}
    }
    
    images = {
      total = length(data.fluence_default_images.images.images)
      by_distribution = {for img in data.fluence_default_images.images.images : img.distribution => img.name...}
      ubuntu_images = [for img in data.fluence_default_images.images.images : img if img.distribution == "Ubuntu"]
      debian_images = [for img in data.fluence_default_images.images.images : img if img.distribution == "Debian"]
    }
  }
  description = "Complete marketplace information"
}

output "cost_estimates" {
  value = {
    small_vm = {
      instances         = data.fluence_vm_estimate_deposit.small_vm.instances
      deposit_amount    = data.fluence_vm_estimate_deposit.small_vm.deposit_amount_usdc
      configuration     = "cpu-2-ram-4gb-storage-25gb"
    }
    
    large_cluster = {
      instances         = data.fluence_vm_estimate_deposit.large_cluster.instances
      deposit_amount    = data.fluence_vm_estimate_deposit.large_cluster.deposit_amount_usdc
      configuration     = "cpu-8-ram-16gb-storage-50gb"
    }
  }
  description = "Cost estimates for different deployment scenarios"
}

# Demonstrate semantic image selection patterns
output "image_selection_examples" {
  value = {
    latest_ubuntu = [for img in data.fluence_default_images.images.images : img if img.distribution == "Ubuntu" && strcontains(img.name, "24.04")][0]
    lts_ubuntu    = [for img in data.fluence_default_images.images.images : img if img.distribution == "Ubuntu" && strcontains(img.name, "22.04")][0]
    debian_stable = [for img in data.fluence_default_images.images.images : img if img.distribution == "Debian"][0]
    
    # Pattern for finding specific versions
    all_ubuntu_versions = [for img in data.fluence_default_images.images.images : {
      name = img.name
      url  = img.download_url
    } if img.distribution == "Ubuntu"]
  }
  description = "Examples of semantic image selection patterns"
}
