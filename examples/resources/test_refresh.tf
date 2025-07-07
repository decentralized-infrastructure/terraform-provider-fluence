# Simple test to check VM refresh functionality
resource "fluence_vm" "refresh_test" {
  name     = "refresh-test-vm"
  os_image = "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img"
  
  # Use existing SSH key fingerprint
  ssh_keys = [
    "SHA256:qncvcamDwxCu8nhUTsZUh+EA7705p+i4tjQj04lOr6U"
  ]
  
  open_ports = [
    {
      port     = 22
      protocol = "tcp"
    }
  ]
  
  instances = 1
}

# Output VM details to check refresh data
output "vm_refresh_test" {
  value = {
    id             = fluence_vm.refresh_test.id
    name           = fluence_vm.refresh_test.name
    status         = fluence_vm.refresh_test.status
    public_ip      = fluence_vm.refresh_test.public_ip
    price_per_epoch = fluence_vm.refresh_test.price_per_epoch
    created_at     = fluence_vm.refresh_test.created_at
  }
}
