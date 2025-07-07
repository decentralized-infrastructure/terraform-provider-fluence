# Test the VM deposit estimation data source

# Basic estimate with just instances
data "fluence_vm_estimate_deposit" "basic" {
  instances = 1
}

# Estimate with basic configuration constraint
data "fluence_vm_estimate_deposit" "with_basic_config" {
  instances = 2
  
  constraints = {
    basic_configuration = "medium"
  }
}

# Estimate with full constraints
data "fluence_vm_estimate_deposit" "with_full_constraints" {
  instances = 3
  
  constraints = {
    basic_configuration           = "large"
    max_total_price_per_epoch_usd = "100.0"
    datacenter_countries          = ["US", "CA", "EU"]
    cpu_architecture              = ["x86_64"]
    cpu_manufacturer              = ["Intel", "AMD"]
    memory_type                   = ["DDR4"]
    storage_type                  = ["SSD"]
  }
}

# Output the estimates for comparison
output "basic_estimate" {
  value = {
    instances              = 1
    deposit_amount_usdc    = data.fluence_vm_estimate_deposit.basic.deposit_amount_usdc
    deposit_epochs         = data.fluence_vm_estimate_deposit.basic.deposit_epochs
    total_price_per_epoch  = data.fluence_vm_estimate_deposit.basic.total_price_per_epoch
    max_price_per_epoch    = data.fluence_vm_estimate_deposit.basic.max_price_per_epoch
  }
}

output "medium_config_estimate" {
  value = {
    instances              = 2
    configuration          = "medium"
    deposit_amount_usdc    = data.fluence_vm_estimate_deposit.with_basic_config.deposit_amount_usdc
    deposit_epochs         = data.fluence_vm_estimate_deposit.with_basic_config.deposit_epochs
    total_price_per_epoch  = data.fluence_vm_estimate_deposit.with_basic_config.total_price_per_epoch
    max_price_per_epoch    = data.fluence_vm_estimate_deposit.with_basic_config.max_price_per_epoch
  }
}

output "full_constraints_estimate" {
  value = {
    instances              = 3
    configuration          = "large"
    max_price_limit        = "100.0"
    deposit_amount_usdc    = data.fluence_vm_estimate_deposit.with_full_constraints.deposit_amount_usdc
    deposit_epochs         = data.fluence_vm_estimate_deposit.with_full_constraints.deposit_epochs
    total_price_per_epoch  = data.fluence_vm_estimate_deposit.with_full_constraints.total_price_per_epoch
    max_price_per_epoch    = data.fluence_vm_estimate_deposit.with_full_constraints.max_price_per_epoch
  }
}

# Compare estimates to see how constraints affect pricing
output "cost_comparison" {
  value = {
    basic_cost_per_vm     = tonumber(data.fluence_vm_estimate_deposit.basic.total_price_per_epoch) / 1
    medium_cost_per_vm    = tonumber(data.fluence_vm_estimate_deposit.with_basic_config.total_price_per_epoch) / 2
    large_cost_per_vm     = tonumber(data.fluence_vm_estimate_deposit.with_full_constraints.total_price_per_epoch) / 3
  }
}
