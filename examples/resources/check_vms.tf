# Just test the VMs data source
data "fluence_vms" "check" {}

output "existing_vms" {
  value = data.fluence_vms.check.vms
}
