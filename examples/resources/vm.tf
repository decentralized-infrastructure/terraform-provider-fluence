# Virtual Machine Management Examples

# Create a VM with basic configuration
resource "fluence_vm" "example" {
  name     = "my-terraform-vm"
  os_image = "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img"
  
  # Reference SSH keys (replace with your SSH key fingerprint)
  ssh_keys = [
    fluence_ssh_key.example.fingerprint  # Reference the SSH key created in ssh_key.tf
  ]
  
  # Open common ports
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
  
  # Optional: VM constraints
  instances                      = 1
  basic_configuration           = "small"
  max_total_price_per_epoch_usd = "10.0"
  datacenter_countries          = ["US", "CA"]
}

# Data source to list all VMs
data "fluence_vms" "all" {}

# Output VM details
output "vm_details" {
  value = {
    id              = fluence_vm.example.id
    name            = fluence_vm.example.name
    status          = fluence_vm.example.status
    public_ip       = fluence_vm.example.public_ip
    price_per_epoch = fluence_vm.example.price_per_epoch
    created_at      = fluence_vm.example.created_at
  }
}

# Output all VMs
output "all_vms" {
  value = data.fluence_vms.all.vms
}
