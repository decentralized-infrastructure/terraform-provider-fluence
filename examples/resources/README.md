# Fluence Terraform Provider Examples

This directory contains example Terraform configurations demonstrating how to use the Fluence provider to manage SSH keys and VMs on the Fluence decentralized compute marketplace.

## Examples Overview

### 1. Simple SSH Key (`ssh_key_simple.tf`)
Standalone example showing how to:
- Create an SSH key with a name
- Access SSH key attributes (fingerprint, algorithm, etc.)
- Output SSH key details

**Use this when**: You only need to manage SSH keys.

### 2. Simple VM (`vm_simple.tf`)
Self-contained example showing how to:
- Create an SSH key for VM access
- Create a VM using that SSH key
- Configure ports, constraints, and basic settings
- Access both SSH key and VM details

**Use this when**: You need a straightforward VM setup with its own SSH key.

### 3. Complete Infrastructure (`complete_infrastructure.tf`)
Comprehensive example demonstrating:
- Using data sources to query marketplace options (countries, hardware, configurations)
- Creating SSH keys and referencing them in VM creation
- Making data-driven infrastructure decisions using available marketplace data
- Comprehensive output of infrastructure details and marketplace options

**Use this when**: You want to leverage marketplace data to make informed infrastructure decisions or need to see the full capabilities of the provider.

## Prerequisites

1. **Fluence API Key**: Get your API key from the [Fluence Console](https://console.fluence.network/settings/api-keys)
2. **Environment Setup**: Set `FLUENCE_API_KEY` environment variable or configure in provider block

## Quick Start

1. Choose an example that fits your use case
2. Copy the example to your working directory
3. Configure your API key:
   ```bash
   export FLUENCE_API_KEY="your-api-key-here"
   ```
4. Initialize and apply:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Provider Configuration

The provider can be configured in multiple ways:

```hcl
provider "fluence" {
  # Option 1: Explicitly set (not recommended for production)
  api_key = "your-api-key"
  host    = "https://api.fluence.dev"  # Optional, defaults to this
  
  # Option 2: Use environment variables (recommended)
  # Set FLUENCE_API_KEY and optionally FLUENCE_HOST
}
```

## Common Patterns

### SSH Key Management
- Always create SSH keys before VMs that need them
- Use meaningful names for easy identification
- Store public keys securely and reference them in configurations

### VM Configuration
- Use data sources to discover available options
- Set budget constraints with `max_total_price_per_epoch_usd`
- Specify location preferences with `datacenter_countries`
- Open only necessary ports for security

### Data-Driven Decisions
- Query `fluence_basic_configurations` to see available VM sizes
- Check `fluence_available_countries` for deployment locations
- Review `fluence_available_hardware` for hardware options
- List existing resources with `fluence_vms` and `fluence_ssh_keys`

## Resource Attributes

### SSH Keys
- **Required**: `name`, `public_key`
- **Computed**: `id`, `fingerprint`, `algorithm`, `comment`, `active`, `created_at`

### VMs
- **Required**: `name`, `os_image`, `ssh_keys`, `open_ports`, `instances`
- **Optional**: `hostname`, `basic_configuration`, `max_total_price_per_epoch_usd`, `datacenter_countries`
- **Computed**: `id`, `status`, `status_changed_at`, `public_ip`, `price_per_epoch`, `created_at`, etc.

## Next Steps

- Check the [Fluence Documentation](https://fluence.dev/docs) for more details
- Review the provider source code for advanced configurations
