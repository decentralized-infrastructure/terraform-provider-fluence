package provider

import (
	"context"
	"fmt"
	"time"

	fluenceapi "github.com/decentralized-infrastructure/fluence-api-client-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &VmResource{}
var _ resource.ResourceWithImportState = &VmResource{}

func NewVmResource() resource.Resource {
	return &VmResource{}
}

// VmResource defines the resource implementation.
type VmResource struct {
	client *fluenceapi.Client
}

// VmResourceModel describes the resource data model.
type VmResourceModel struct {
	ID        types.String    `tfsdk:"id"`
	Name      types.String    `tfsdk:"name"`
	Hostname  types.String    `tfsdk:"hostname"`
	OsImage   types.String    `tfsdk:"os_image"`
	SshKeys   []types.String  `tfsdk:"ssh_keys"`
	OpenPorts []OpenPortModel `tfsdk:"open_ports"`
	Instances types.Int64     `tfsdk:"instances"`

	// Constraints (optional)
	BasicConfiguration       types.String   `tfsdk:"basic_configuration"`
	MaxTotalPricePerEpochUsd types.String   `tfsdk:"max_total_price_per_epoch_usd"`
	Countries                []types.String `tfsdk:"datacenter_countries"`

	// Computed fields
	Status          types.String `tfsdk:"status"`
	StatusChangedAt types.String `tfsdk:"status_changed_at"`
	PricePerEpoch   types.String `tfsdk:"price_per_epoch"`
	CreatedAt       types.String `tfsdk:"created_at"`
	NextBillingAt   types.String `tfsdk:"next_billing_at"`
	ReservedBalance types.String `tfsdk:"reserved_balance"`
	TotalSpent      types.String `tfsdk:"total_spent"`
	PublicIp        types.String `tfsdk:"public_ip"`
}

// OpenPortModel represents an open port configuration
type OpenPortModel struct {
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
}

func (r *VmResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm"
}

func (r *VmResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Virtual Machine resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "VM identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "VM name",
				Required:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "VM hostname (optional)",
				Optional:            true,
			},
			"os_image": schema.StringAttribute{
				MarkdownDescription: "Operating system image to use",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ssh_keys": schema.ListAttribute{
				MarkdownDescription: "List of SSH key fingerprints to authorize",
				ElementType:         types.StringType,
				Required:            true,
			},
			"open_ports": schema.ListNestedAttribute{
				MarkdownDescription: "List of ports to open on the VM",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"port": schema.Int64Attribute{
							MarkdownDescription: "Port number",
							Required:            true,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "Protocol (tcp/udp)",
							Required:            true,
						},
					},
				},
			},
			"instances": schema.Int64Attribute{
				MarkdownDescription: "Number of VM instances to create",
				Optional:            true,
				Computed:            true,
			},

			// Constraint attributes
			"basic_configuration": schema.StringAttribute{
				MarkdownDescription: "Basic configuration constraint",
				Optional:            true,
			},
			"max_total_price_per_epoch_usd": schema.StringAttribute{
				MarkdownDescription: "Maximum total price per epoch in USD",
				Optional:            true,
			},
			"datacenter_countries": schema.ListAttribute{
				MarkdownDescription: "List of allowed datacenter countries",
				ElementType:         types.StringType,
				Optional:            true,
			},

			// Computed attributes
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "VM status",
			},
			"status_changed_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "VM status change timestamp",
			},
			"price_per_epoch": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Price per epoch",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "VM creation time",
			},
			"next_billing_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Next billing time",
			},
			"reserved_balance": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Reserved balance",
			},
			"total_spent": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Total amount spent",
			},
			"public_ip": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Public IP address of the VM",
			},
		},
	}
}

func (r *VmResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*fluenceapi.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *fluenceapi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *VmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VmResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build the VM configuration
	vmConfig := fluenceapi.VmConfiguration{
		Name:      data.Name.ValueString(),
		OsImage:   data.OsImage.ValueString(),
		OpenPorts: []fluenceapi.OpenPorts{},
		SshKeys:   []string{},
	}

	// Set hostname if provided
	if !data.Hostname.IsNull() && !data.Hostname.IsUnknown() {
		hostname := data.Hostname.ValueString()
		vmConfig.Hostname = &hostname
	}

	// Convert SSH keys
	for _, sshKey := range data.SshKeys {
		vmConfig.SshKeys = append(vmConfig.SshKeys, sshKey.ValueString())
	}

	// Convert open ports
	for _, port := range data.OpenPorts {
		vmConfig.OpenPorts = append(vmConfig.OpenPorts, fluenceapi.OpenPorts{
			Port:     uint16(port.Port.ValueInt64()),
			Protocol: port.Protocol.ValueString(),
		})
	}

	// Build constraints
	var constraints *fluenceapi.OfferConstraints
	if !data.BasicConfiguration.IsNull() || !data.MaxTotalPricePerEpochUsd.IsNull() || len(data.Countries) > 0 {
		constraints = &fluenceapi.OfferConstraints{}

		if !data.BasicConfiguration.IsNull() {
			basicConfig := data.BasicConfiguration.ValueString()
			constraints.BasicConfiguration = &basicConfig
		}

		if !data.MaxTotalPricePerEpochUsd.IsNull() {
			maxPrice := data.MaxTotalPricePerEpochUsd.ValueString()
			constraints.MaxTotalPricePerEpochUsd = &maxPrice
		}

		if len(data.Countries) > 0 {
			countries := []string{}
			for _, country := range data.Countries {
				countries = append(countries, country.ValueString())
			}
			constraints.Datacenter = &fluenceapi.DatacenterConstraint{
				Countries: countries,
			}
		}
	}

	// Set instances (default to 1 if not specified)
	instances := 1
	if !data.Instances.IsNull() && !data.Instances.IsUnknown() {
		instances = int(data.Instances.ValueInt64())
	}

	// Create VM using the client
	createReq := fluenceapi.CreateVmV3{
		Constraints:     constraints,
		Instances:       instances,
		VmConfiguration: vmConfig,
	}

	createdVms, err := r.client.CreateVmV3(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create VM, got error: %s", err))
		return
	}

	if len(createdVms) == 0 {
		resp.Diagnostics.AddError("Client Error", "No VMs were created")
		return
	}

	// For now, we'll handle single VM instances. Use the first created VM as the resource ID
	createdVm := createdVms[0]
	data.ID = types.StringValue(createdVm.VmId)
	data.Instances = types.Int64Value(int64(instances))

	// Set the VM name from the response if available
	if createdVm.VmName != "" {
		data.Name = types.StringValue(createdVm.VmName)
	}

	// Wait for VM to become active before considering creation complete
	err = r.waitForVmActive(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("VM Creation Error", fmt.Sprintf("VM was created but failed to become active: %s", err))
		return
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created VM resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VmResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh the VM data
	err := r.refreshVmData(ctx, &data)
	if err != nil {
		if err.Error() == "VM not found" {
			// VM not found, remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to refresh VM data: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// refreshVmData fetches the current VM data and updates the model
func (r *VmResource) refreshVmData(ctx context.Context, data *VmResourceModel) error {
	vmId := data.ID.ValueString()
	tflog.Debug(ctx, "Attempting to refresh VM data", map[string]interface{}{
		"vm_id": vmId,
	})

	// Retry logic for newly created VMs that might not be immediately available
	maxRetries := 5
	retryDelay := time.Second * 2

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			tflog.Debug(ctx, "Retrying VM data refresh", map[string]interface{}{
				"attempt": attempt + 1,
				"vm_id":   vmId,
			})
			time.Sleep(retryDelay)
		}

		// Get all VMs and find the one matching our ID
		vms, err := r.client.ListVmsV3()
		if err != nil {
			return fmt.Errorf("unable to read VMs: %s", err)
		}

		tflog.Debug(ctx, "Retrieved VMs from API", map[string]interface{}{
			"vm_count": len(vms),
			"vm_id":    vmId,
		})

		// Find the VM by ID
		var foundVm *fluenceapi.RunningInstanceV3
		for _, vm := range vms {
			tflog.Trace(ctx, "Checking VM", map[string]interface{}{
				"api_vm_id":    vm.Id,
				"target_vm_id": vmId,
				"match":        vm.Id == vmId,
			})
			if vm.Id == vmId {
				foundVm = &vm
				break
			}
		}

		if foundVm != nil {
			tflog.Debug(ctx, "Found VM in API", map[string]interface{}{
				"vm_id":   vmId,
				"status":  foundVm.Status,
				"attempt": attempt + 1,
			})

			// Update the model with the current data
			data.Status = types.StringValue(foundVm.Status)
			data.StatusChangedAt = types.StringValue(foundVm.StatusChangedAt)
			data.PricePerEpoch = types.StringValue(foundVm.PricePerEpoch)
			data.CreatedAt = types.StringValue(foundVm.CreatedAt)
			data.NextBillingAt = types.StringValue(foundVm.NextBillingAt)
			data.ReservedBalance = types.StringValue(foundVm.ReservedBalance)
			data.TotalSpent = types.StringValue(foundVm.TotalSpent)

			if foundVm.OsImage != nil {
				data.OsImage = types.StringValue(*foundVm.OsImage)
			}

			if foundVm.PublicIp != nil {
				data.PublicIp = types.StringValue(*foundVm.PublicIp)
			} else {
				data.PublicIp = types.StringNull()
			}

			if foundVm.VmName != nil {
				data.Name = types.StringValue(*foundVm.VmName)
			}

			return nil
		}

		tflog.Debug(ctx, "VM not found in API response", map[string]interface{}{
			"vm_id":   vmId,
			"attempt": attempt + 1,
		})
	}

	return fmt.Errorf("VM not found after %d attempts", maxRetries)
}

func (r *VmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VmResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build update request - only name and open ports can be updated
	updates := []fluenceapi.UpdateVm{
		{
			Id: data.ID.ValueString(),
		},
	}

	// Update VM name if changed
	if !data.Name.IsNull() {
		vmName := data.Name.ValueString()
		updates[0].VmName = &vmName
	}

	// Update open ports if changed
	if len(data.OpenPorts) > 0 {
		openPorts := []fluenceapi.OpenPorts{}
		for _, port := range data.OpenPorts {
			openPorts = append(openPorts, fluenceapi.OpenPorts{
				Port:     uint16(port.Port.ValueInt64()),
				Protocol: port.Protocol.ValueString(),
			})
		}
		updates[0].OpenPorts = &openPorts
	}

	err := r.client.UpdateVms(updates)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update VM, got error: %s", err))
		return
	}

	// Refresh the data after update
	err = r.refreshVmData(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to refresh VM data after update: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VmResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the VM using the VM ID
	vmIds := []string{data.ID.ValueString()}
	_, err := r.client.RemoveVms(vmIds)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete VM, got error: %s", err))
		return
	}
}

func (r *VmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID field for import
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// waitForVmActive waits for a VM to reach the "Active" status
func (r *VmResource) waitForVmActive(ctx context.Context, data *VmResourceModel) error {
	vmId := data.ID.ValueString()
	tflog.Debug(ctx, "Waiting for VM to become active", map[string]interface{}{
		"vm_id": vmId,
	})

	// Use a fixed timeout of 10 minutes for VM creation
	createTimeout := 10 * time.Minute

	retryInterval := time.Second * 10 // Check every 10 seconds
	maxRetries := int(createTimeout / retryInterval)

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			tflog.Debug(ctx, "Waiting for VM status check", map[string]interface{}{
				"attempt":      attempt + 1,
				"max_attempts": maxRetries,
				"vm_id":        vmId,
			})
			time.Sleep(retryInterval)
		}

		// Get all VMs and find the one matching our ID
		vms, err := r.client.ListVmsV3()
		if err != nil {
			tflog.Warn(ctx, "Error retrieving VMs while waiting for availability", map[string]interface{}{
				"error":   err.Error(),
				"attempt": attempt + 1,
				"vm_id":   vmId,
			})
			continue // Continue retrying on API errors
		}

		// Find the VM by ID
		var foundVm *fluenceapi.RunningInstanceV3
		for _, vm := range vms {
			if vm.Id == vmId {
				foundVm = &vm
				break
			}
		}

		if foundVm == nil {
			tflog.Warn(ctx, "VM not found while waiting for availability", map[string]interface{}{
				"attempt": attempt + 1,
				"vm_id":   vmId,
			})
			continue // Continue retrying if VM not found
		}

		tflog.Debug(ctx, "Checking VM status", map[string]interface{}{
			"vm_id":   vmId,
			"status":  foundVm.Status,
			"attempt": attempt + 1,
		})

		// Update the model with current data
		data.Status = types.StringValue(foundVm.Status)
		data.StatusChangedAt = types.StringValue(foundVm.StatusChangedAt)
		data.PricePerEpoch = types.StringValue(foundVm.PricePerEpoch)
		data.CreatedAt = types.StringValue(foundVm.CreatedAt)
		data.NextBillingAt = types.StringValue(foundVm.NextBillingAt)
		data.ReservedBalance = types.StringValue(foundVm.ReservedBalance)
		data.TotalSpent = types.StringValue(foundVm.TotalSpent)

		if foundVm.OsImage != nil {
			data.OsImage = types.StringValue(*foundVm.OsImage)
		}

		if foundVm.PublicIp != nil {
			data.PublicIp = types.StringValue(*foundVm.PublicIp)
		} else {
			data.PublicIp = types.StringNull()
		}

		if foundVm.VmName != nil {
			data.Name = types.StringValue(*foundVm.VmName)
		}

		// Check if VM has reached active status
		if foundVm.Status == "Active" {
			tflog.Info(ctx, "VM is now active", map[string]interface{}{
				"vm_id":        vmId,
				"total_time":   time.Duration(attempt) * retryInterval,
				"final_status": foundVm.Status,
			})
			return nil
		}

		// Check for failure states that we should not wait through
		if foundVm.Status == "failed" || foundVm.Status == "error" || foundVm.Status == "Failed" || foundVm.Status == "Error" {
			return fmt.Errorf("VM creation failed with status: %s", foundVm.Status)
		}

		tflog.Debug(ctx, "VM not yet active, continuing to wait", map[string]interface{}{
			"vm_id":          vmId,
			"current_status": foundVm.Status,
			"time_elapsed":   time.Duration(attempt+1) * retryInterval,
			"time_remaining": createTimeout - (time.Duration(attempt+1) * retryInterval),
		})
	}

	// If we've exhausted all retries, return the current status in the error
	currentStatus := "unknown"
	if !data.Status.IsNull() {
		currentStatus = data.Status.ValueString()
	}

	return fmt.Errorf("VM did not become active within %v (current status: %s)", createTimeout, currentStatus)
}
