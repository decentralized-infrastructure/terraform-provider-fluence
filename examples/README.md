# Fluence Terraform Provider Examples

This directory contains examples of how to use the Fluence Terraform provider to manage SSH keys and virtual machines.

## Prerequisites

1. Set your Fluence API key as an environment variable:
   ```bash
   export FLUENCE_API_KEY="your-api-key-here"
   ```

2. Build and install the provider (from the root directory):
   ```bash
   go build -o terraform-provider-fluence
   ```

## Examples

### Data Sources

- `data-sources/data-source.tf` - Shows how to read existing SSH keys from the Fluence API

### Resources

- `resources/provider.tf` - Provider configuration
- `resources/ssh_key.tf` - SSH key management examples
- `resources/vm.tf` - Virtual machine management examples  
- `resources/marketplace.tf` - Marketplace data sources examples
- `resources/complete_infrastructure.tf` - Complete infrastructure example with multiple VMs and SSH keys
- `resources/simple_vm.tf` - Simple VM creation example
- `resources/complete_ssh_example.tf` - Comprehensive SSH key examples

## Features

### SSH Key Management
- Create, read, update, and delete SSH keys
- List all SSH keys in your account
- Import existing SSH keys into Terraform state

### Virtual Machine Management
- Create, read, update, and delete virtual machines
- Configure VM constraints (CPU, memory, storage, datacenter location)
- Set pricing limits and basic configurations
- Open specific ports on VMs
- List all VMs in your account

### Marketplace Integration
- Query available VM configurations
- Get list of available datacenter countries
- Retrieve available hardware options (CPU, memory, storage)

## Usage

1. Navigate to the example directory you want to try
2. Replace the placeholder values with your actual configuration:
   - SSH public keys
   - VM names and configurations
   - Pricing limits
3. Run the standard Terraform commands:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## SSH Key Format

Make sure your SSH public keys are in the correct format:

- RSA keys: `ssh-rsa AAAAB3NzaC1yc2EAAAA... comment`
- ED25519 keys: `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA... comment`
- ECDSA keys: `ssh-ecdsa AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAI... comment`

You can generate a new SSH key pair using:
```bash
ssh-keygen -t ed25519 -C "your_email@example.com"
```

The public key will be in `~/.ssh/id_ed25519.pub` (or similar, depending on the key type).
