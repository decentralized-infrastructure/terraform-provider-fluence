---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fluence_vm_estimate_deposit Data Source - terraform-provider-fluence"
subcategory: ""
description: |-
  Estimate the deposit required for creating VMs with given configuration and constraints
---

# fluence_vm_estimate_deposit (Data Source)

Estimate the deposit required for creating VMs with given configuration and constraints



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `instances` (Number) Number of VM instances to estimate for

### Optional

- `constraints` (Attributes) Constraints for the VM estimation (see [below for nested schema](#nestedatt--constraints))

### Read-Only

- `deposit_amount_usdc` (String) Required deposit amount in USDC
- `deposit_epochs` (Number) Number of epochs the deposit covers
- `max_price_per_epoch` (String) Maximum price per epoch for all instances
- `total_price_per_epoch` (String) Total price per epoch for all instances

<a id="nestedatt--constraints"></a>
### Nested Schema for `constraints`

Optional:

- `basic_configuration` (String) Basic configuration constraint (e.g., 'small', 'medium', 'large')
- `cpu_architecture` (List of String) List of allowed CPU architectures
- `cpu_manufacturer` (List of String) List of allowed CPU manufacturers
- `datacenter_countries` (List of String) List of allowed datacenter countries
- `max_total_price_per_epoch_usd` (String) Maximum total price per epoch in USD
- `memory_generation` (List of String) List of allowed memory generations
- `memory_type` (List of String) List of allowed memory types
- `storage_type` (List of String) List of allowed storage types
