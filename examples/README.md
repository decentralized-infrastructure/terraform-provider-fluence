# Fluence Terraform Provider Examples

This directory contains focused examples demonstrating how to use the Fluence Terraform provider to manage infrastructure on the Fluence decentralized compute marketplace.

## 📁 Directory Structure

```
examples/
├── README.md                     # This file
├── data-sources/                 # Data source examples
│   ├── README.md                 # Data sources documentation
│   └── data-source.tf           # Query marketplace and existing resources
└── resources/                    # Resource management examples
    ├── README.md                 # Resources documentation
    ├── provider.tf               # Provider configuration
    ├── ssh_key_simple.tf         # Simple SSH key example
    ├── vm_simple.tf              # Simple VM example
    └── complete_infrastructure.tf # Complete infrastructure example
```

## 🚀 Quick Start

1. **Get API Key**: Obtain your API key from [Fluence Console](https://console.fluence.network/settings/api-keys)

2. **Set Environment Variable**:
   ```bash
   export FLUENCE_API_KEY="your-api-key-here"
   ```

3. **Choose an Example**: Navigate to the appropriate directory and copy the example that fits your needs

4. **Initialize and Apply**:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## 📚 Examples Overview

### 🔐 SSH Key Management (`resources/ssh_key_simple.tf`)
**Perfect for**: Learning SSH key basics only
- Create and manage SSH keys
- Access key attributes (fingerprint, algorithm, etc.)
- Basic output examples

### 🖥️ Simple VM (`resources/vm_simple.tf`)  
**Perfect for**: Self-contained VM setup
- Create an SSH key and VM together
- Configure ports, SSH access, and constraints
- See both SSH key and VM outputs

### 🏗️ Complete Infrastructure (`resources/complete_infrastructure.tf`)
**Perfect for**: Data-driven, production-ready setups
- Use marketplace data sources for optimal configuration
- Create SSH keys and reference them in VMs
- Make informed infrastructure decisions using available data
- Comprehensive infrastructure and marketplace outputs

### 📊 Data Sources (`data-sources/data-source.tf`)
**Perfect for**: Discovery and inventory
- Query all available marketplace options
- List existing resources in your account
- Make informed infrastructure planning decisions

## 🎯 Learning Path

1. **Start Here**: `data-sources/` - Understand what's available
2. **SSH Keys**: `resources/ssh_key_simple.tf` - Learn key management basics
3. **Basic VM**: `resources/vm_simple.tf` - Create a self-contained VM setup
4. **Advanced**: `resources/complete_infrastructure.tf` - Data-driven production patterns

## 🔧 Prerequisites

- **Terraform**: Version 1.0 or later
- **Fluence API Key**: From [Fluence Console](https://console.fluence.network)
- **SSH Public Key**: For VM access (can generate with `ssh-keygen`)

## 💡 Common Patterns

- **Environment Variables**: Use `FLUENCE_API_KEY` instead of hardcoding credentials
- **Data Sources First**: Query available options before creating resources
- **SSH Keys Before VMs**: Always create SSH keys before VMs that reference them
- **Budget Constraints**: Set `max_total_price_per_epoch_usd` to control costs
- **Location Preferences**: Use `datacenter_countries` for optimal latency

## 📖 Next Steps

- Review individual example READMEs for detailed usage
- Check the [Fluence Documentation](https://fluence.dev/docs) for platform details
- Explore the `swagger.json` file for complete API reference
- Join the [Fluence Community](https://fluence.dev/community) for support

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
