package provider

import (
	"context"
	"fmt"

	fluenceapi "github.com/0xthresh/fluence-api-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &vmsDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *vmsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*fluenceapi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *fluenceapi.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// NewVmsDataSource is a helper function to simplify the provider implementation.
func NewVmsDataSource() datasource.DataSource {
	return &vmsDataSource{}
}

// vmsDataSource is the data source implementation.
type vmsDataSource struct {
	client *fluenceapi.Client
}

// Metadata returns the data source type name.
func (d *vmsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vms"
}

// Schema defines the schema for the data source.
func (d *vmsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vms": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"price_per_epoch": schema.StringAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"next_billing_at": schema.StringAttribute{
							Computed: true,
						},
						"reserved_balance": schema.StringAttribute{
							Computed: true,
						},
						"total_spent": schema.StringAttribute{
							Computed: true,
						},
						"os_image": schema.StringAttribute{
							Computed: true,
						},
						"public_ip": schema.StringAttribute{
							Computed: true,
						},
						"vm_name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// vmsDataSourceModel maps the data source schema data.
type vmsDataSourceModel struct {
	Vms []vmModel `tfsdk:"vms"`
}

// vmModel maps VM data.
type vmModel struct {
	ID              types.String `tfsdk:"id"`
	Status          types.String `tfsdk:"status"`
	PricePerEpoch   types.String `tfsdk:"price_per_epoch"`
	CreatedAt       types.String `tfsdk:"created_at"`
	NextBillingAt   types.String `tfsdk:"next_billing_at"`
	ReservedBalance types.String `tfsdk:"reserved_balance"`
	TotalSpent      types.String `tfsdk:"total_spent"`
	OsImage         types.String `tfsdk:"os_image"`
	PublicIp        types.String `tfsdk:"public_ip"`
	VmName          types.String `tfsdk:"vm_name"`
}

// Read refreshes the Terraform state with the latest data.
func (d *vmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vmsDataSourceModel

	vms, err := d.client.ListVmsV3()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Fluence VMs",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, vm := range vms {
		vmState := vmModel{
			ID:              types.StringValue(vm.Id),
			Status:          types.StringValue(vm.Status),
			PricePerEpoch:   types.StringValue(vm.PricePerEpoch),
			CreatedAt:       types.StringValue(vm.CreatedAt),
			NextBillingAt:   types.StringValue(vm.NextBillingAt),
			ReservedBalance: types.StringValue(vm.ReservedBalance),
			TotalSpent:      types.StringValue(vm.TotalSpent),
		}

		if vm.OsImage != nil {
			vmState.OsImage = types.StringValue(*vm.OsImage)
		} else {
			vmState.OsImage = types.StringNull()
		}

		if vm.PublicIp != nil {
			vmState.PublicIp = types.StringValue(*vm.PublicIp)
		} else {
			vmState.PublicIp = types.StringNull()
		}

		if vm.VmName != nil {
			vmState.VmName = types.StringValue(*vm.VmName)
		} else {
			vmState.VmName = types.StringNull()
		}

		state.Vms = append(state.Vms, vmState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
