# Create SSH keys for VM access - using valid SSH keys
resource "fluence_ssh_key" "admin" {
  name       = "admin-key"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAD0WOgqjaY9EBuhEYg0nTQNuHwGH0Tg/YgtS57VF4g9 admin@example.com"
}

resource "fluence_ssh_key" "deploy" {
  name       = "deploy-key"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKz6FlQh0x1gzrHjWyKW5W1z6DqGLv6Ks1Z1Z1Z1Z1Z1 deploy@example.com"
}

# Create a web server VM
resource "fluence_vm" "web_server" {
  name     = "web-server"
  hostname = "web01"
  os_image = "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img"
  
  # Use the SSH keys we just created
  ssh_keys = [
    fluence_ssh_key.admin.fingerprint,
    fluence_ssh_key.deploy.fingerprint
  ]
  
  # Open web server ports
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
    },
    {
      port     = 3000
      protocol = "tcp"
    }
  ]
  
  # VM configuration
  instances                      = 1
  basic_configuration           = "medium"
  max_total_price_per_epoch_usd = "25.0"
  datacenter_countries          = ["US", "CA", "EU"]
}

# Create a database VM
resource "fluence_vm" "database" {
  name     = "database"
  hostname = "db01"
  os_image = "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img"
  
  # Use only the admin SSH key for the database
  ssh_keys = [
    fluence_ssh_key.admin.fingerprint
  ]
  
  # Open database ports
  open_ports = [
    {
      port     = 22
      protocol = "tcp"
    },
    {
      port     = 5432
      protocol = "tcp"
    },
    {
      port     = 6379
      protocol = "tcp"
    }
  ]
  
  # Database VM configuration
  instances                      = 1
  basic_configuration           = "large"
  max_total_price_per_epoch_usd = "50.0"
  datacenter_countries          = ["US", "CA"]
}

# Read all VMs to see what we've created
data "fluence_vms" "all" {
  depends_on = [
    fluence_vm.web_server,
    fluence_vm.database
  ]
}

# Outputs
output "ssh_keys" {
  value = {
    admin  = fluence_ssh_key.admin.fingerprint
    deploy = fluence_ssh_key.deploy.fingerprint
  }
}

output "web_server" {
  value = {
    id        = fluence_vm.web_server.id
    name      = fluence_vm.web_server.name
    hostname  = fluence_vm.web_server.hostname
    status    = fluence_vm.web_server.status
    public_ip = fluence_vm.web_server.public_ip
    price     = fluence_vm.web_server.price_per_epoch
  }
}

output "database" {
  value = {
    id        = fluence_vm.database.id
    name      = fluence_vm.database.name
    hostname  = fluence_vm.database.hostname
    status    = fluence_vm.database.status
    public_ip = fluence_vm.database.public_ip
    price     = fluence_vm.database.price_per_epoch
  }
}

output "all_vms" {
  value = data.fluence_vms.all
}
