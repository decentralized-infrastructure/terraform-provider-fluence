# Fluence Terraform Provider Examples

This directory contains streamlined examples demonstrating how to use the Fluence provider v1.1.0 to manage VMs on the Fluence decentralized compute marketplace.

## Quick Start

1. **Set your API key**:
   ```bash
   export FLUENCE_API_KEY="your-api-key-here"
   ```

2. **Choose an example** and navigate to its directory:
   ```bash
   cd resources/basic-vm
   # or
   cd resources/advanced-infrastructure
   # or  
   cd data-sources/marketplace-showcase
   ```

3. **Run the example**:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Examples Overview

### Resources

- **`resources/basic-vm/`** - Simple VM creation using existing SSH keys and semantic image selection
- **`resources/advanced-infrastructure/`** - Multiple VMs with hardware constraints, additional resources, and cost estimation

### Data Sources

- **`data-sources/marketplace-showcase/`** - Comprehensive demonstration of all available data sources and marketplace discovery features

## Prerequisites

- **Fluence API Key**: Get yours from the [Fluence Console](https://console.fluence.network/settings/api-keys)
- **Existing SSH Keys**: Have at least one SSH key in your Fluence account (examples use existing keys)
- **Terraform**: Version 1.0+ with provider version ~> 1.0
