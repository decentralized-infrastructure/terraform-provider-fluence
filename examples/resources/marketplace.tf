# Marketplace Data Sources Examples

# Get available basic configurations
data "fluence_basic_configurations" "available" {}

# Get available countries
data "fluence_available_countries" "available" {}

# Get available hardware options
data "fluence_available_hardware" "available" {}

# Output marketplace information
output "marketplace_info" {
  value = {
    basic_configurations = data.fluence_basic_configurations.available.configurations
    available_countries  = data.fluence_available_countries.available.countries
    hardware_options = {
      cpu     = data.fluence_available_hardware.available.cpu
      memory  = data.fluence_available_hardware.available.memory
      storage = data.fluence_available_hardware.available.storage
    }
  }
}

# Example: Use marketplace data in VM configuration
resource "fluence_vm" "marketplace_example" {
  name     = "marketplace-configured-vm"
  os_image = "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img"
  
  ssh_keys = [
    fluence_ssh_key.example.fingerprint
  ]
  
  open_ports = [
    {
      port     = 22
      protocol = "tcp"
    }
  ]
  
  # Use marketplace data to set configuration
  instances               = 1
  basic_configuration     = data.fluence_basic_configurations.available.configurations[0]  # Use first available config
  datacenter_countries    = [data.fluence_available_countries.available.countries[0]]      # Use first available country
}

output "marketplace_vm_details" {
  value = {
    vm_id         = fluence_vm.marketplace_example.id
    configuration = fluence_vm.marketplace_example.basic_configuration
    countries     = fluence_vm.marketplace_example.datacenter_countries
    status        = fluence_vm.marketplace_example.status
  }
}
