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

# Get comprehensive marketplace data using v1.1.0 features
data "fluence_datacenters" "available" {}
data "fluence_default_images" "available" {}
data "fluence_available_countries" "available" {}
data "fluence_available_hardware" "available" {}
data "fluence_basic_configurations" "available" {}
data "fluence_ssh_keys" "existing" {}

# Select different OS images using semantic selection
locals {
  ubuntu_24_image = [for img in data.fluence_default_images.available.images : img if img.distribution == "Ubuntu" && strcontains(img.name, "24.04")][0]
  ubuntu_22_image = [for img in data.fluence_default_images.available.images : img if img.distribution == "Ubuntu" && strcontains(img.name, "22.04")][0]
  debian_image    = [for img in data.fluence_default_images.available.images : img if img.distribution == "Debian"][0]
}

# Create a web server VM with advanced configuration
resource "fluence_vm" "web_server" {
  name     = "advanced-web-server"
  hostname = "web-01"
  
  # Use Ubuntu 24.04 for web server
  os_image = local.ubuntu_24_image.download_url
  
  # Use existing SSH key
  ssh_keys = [data.fluence_ssh_keys.existing.ssh_keys[0].fingerprint]
  
  # Open ports for web services
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
  
  # Advanced configuration with constraints
  instances           = 1
  basic_configuration = "cpu-8-ram-16gb-storage-50gb"
  
  # Hardware constraints (v1.1.0 feature)
  hardware_constraints = [
    {
      cpu_architecture  = ["x86_64"]
      cpu_manufacturer  = ["Intel", "AMD"]
      memory_type       = ["DDR4"]
      memory_generation = ["4"]
      storage_type      = ["SSD"]
    }
  ]
  
  # Additional resources (v1.1.0 feature)
  additional_resources = [
    {
      storage = [
        {
          supply = 25
          units  = "GB"
          type   = "SSD"
        }
      ]
    }
  ]
  
  # Location and budget constraints
  datacenter_countries          = ["US", "DE", "CA"]
  max_total_price_per_epoch_usd = "15.0"
  
  # Configure timeouts
  timeouts {
    create = "20m"
  }
}

# Create a database server VM with different OS
resource "fluence_vm" "database_server" {
  name     = "database-server"
  hostname = "db-01"
  
  # Use Ubuntu 22.04 for database
  os_image = local.ubuntu_22_image.download_url
  
  # Use existing SSH key
  ssh_keys = [data.fluence_ssh_keys.existing.ssh_keys[0].fingerprint]
  
  # Database-specific ports
  open_ports = [
    {
      port     = 22
      protocol = "tcp"
    },
    {
      port     = 3306  # MySQL
      protocol = "tcp"
    },
    {
      port     = 5432  # PostgreSQL
      protocol = "tcp"
    }
  ]
  
  # High-performance configuration for database
  instances           = 1
  basic_configuration = "cpu-8-ram-32gb-storage-100gb"
  
  # Optimized for database workloads
  hardware_constraints = [
    {
      cpu_architecture  = ["x86_64"]
      cpu_manufacturer  = ["Intel"]
      memory_type       = ["DDR4"]
      storage_type      = ["NVMe", "SSD"]
    }
  ]
  
  # Location preference for low latency
  datacenter_countries          = ["US"]
  max_total_price_per_epoch_usd = "25.0"
  
  timeouts {
    create = "20m"
  }
}

# Output comprehensive information
output "web_server_info" {
  value = {
    id        = fluence_vm.web_server.id
    hostname  = fluence_vm.web_server.hostname
    status    = fluence_vm.web_server.status
    public_ip = fluence_vm.web_server.public_ip
    price     = fluence_vm.web_server.price_per_epoch
    os_image  = fluence_vm.web_server.os_image
  }
  description = "Web server VM details"
}

output "database_server_info" {
  value = {
    id        = fluence_vm.database_server.id
    hostname  = fluence_vm.database_server.hostname
    status    = fluence_vm.database_server.status
    public_ip = fluence_vm.database_server.public_ip
    price     = fluence_vm.database_server.price_per_epoch
    os_image  = fluence_vm.database_server.os_image
  }
  description = "Database server VM details"
}

output "marketplace_summary" {
  value = {
    available_countries     = length(data.fluence_available_countries.available.countries)
    available_configurations = length(data.fluence_basic_configurations.available.configurations)
    available_images        = length(data.fluence_default_images.available.images)
    available_datacenters   = length(data.fluence_datacenters.available.datacenters)
  }
  description = "Summary of available marketplace options"
}
