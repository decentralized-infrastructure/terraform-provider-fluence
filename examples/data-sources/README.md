# Fluence Data Sources Example

This example demonstrates how to query all available data sources from the Fluence provider to get comprehensive information about the marketplace and your existing resources.

## What This Example Shows

### Marketplace Data Sources
- **`fluence_available_countries`**: Lists countries where VMs can be deployed
- **`fluence_available_hardware`**: Shows available CPU, memory, and storage options
- **`fluence_basic_configurations`**: Lists predefined VM configuration sizes

### Resource Data Sources
- **`fluence_ssh_keys`**: Lists all SSH keys in your account
- **`fluence_vms`**: Lists all VMs in your account with their current status

## Usage

1. Set your API key:
   ```bash
   export FLUENCE_API_KEY="your-api-key-here"
   ```

2. Run the example:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

3. Review the outputs to understand:
   - What deployment options are available
   - What resources you currently have
   - How to make informed infrastructure decisions

## Sample Output

The example outputs structured information like:

```
marketplace_info = {
  available_countries = ["BE", "PL", "US"]
  basic_configurations = [
    "cpu-2-ram-4gb-storage-25gb",
    "cpu-4-ram-8gb-storage-25gb",
    # ... more configurations
  ]
  hardware_options = {
    cpu_types = [
      { architecture = "ZEN", manufacturer = "AMD" },
      # ... more CPU options
    ]
    # ... memory and storage options
  }
}

existing_resources = {
  ssh_keys = {
    count = 2
    keys = [
      {
        name = "my-key"
        fingerprint = "SHA256:..."
        algorithm = "ssh-ed25519"
        created_at = "2024-01-01T00:00:00Z"
      }
      # ... more keys
    ]
  }
  vms = {
    count = 1
    summary = [
      {
        id = "0x..."
        name = "my-vm"
        status = "Active"
        public_ip = "1.2.3.4"
        price_per_epoch = "1.23"
        created_at = "2024-01-01T00:00:00Z"
      }
      # ... more VMs
    ]
  }
}
```
