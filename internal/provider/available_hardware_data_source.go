package provider

import (
	"context"
	"fmt"

	fluenceapi "github.com/0xthresh/fluence-api-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &availableHardwareDataSource{}

func NewAvailableHardwareDataSource() datasource.DataSource {
	return &availableHardwareDataSource{}
}

// availableHardwareDataSource defines the data source implementation.
type availableHardwareDataSource struct {
	client *fluenceapi.Client
}

// AvailableHardwareDataSourceModel describes the data source data model.
type AvailableHardwareDataSourceModel struct {
	Cpu     []CpuHardwareModel     `tfsdk:"cpu"`
	Memory  []MemoryHardwareModel  `tfsdk:"memory"`
	Storage []StorageHardwareModel `tfsdk:"storage"`
}

// CpuHardwareModel represents CPU hardware options
type CpuHardwareModel struct {
	Architecture types.String `tfsdk:"architecture"`
	Manufacturer types.String `tfsdk:"manufacturer"`
}

// MemoryHardwareModel represents memory hardware options
type MemoryHardwareModel struct {
	Type       types.String `tfsdk:"type"`
	Generation types.String `tfsdk:"generation"`
}

// StorageHardwareModel represents storage hardware options
type StorageHardwareModel struct {
	Type types.String `tfsdk:"type"`
}

func (d *availableHardwareDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_available_hardware"
}

func (d *availableHardwareDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetch available hardware options from the marketplace",

		Attributes: map[string]schema.Attribute{
			"cpu": schema.ListNestedAttribute{
				MarkdownDescription: "Available CPU hardware options",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"architecture": schema.StringAttribute{
							MarkdownDescription: "CPU architecture",
							Computed:            true,
						},
						"manufacturer": schema.StringAttribute{
							MarkdownDescription: "CPU manufacturer",
							Computed:            true,
						},
					},
				},
			},
			"memory": schema.ListNestedAttribute{
				MarkdownDescription: "Available memory hardware options",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Memory type",
							Computed:            true,
						},
						"generation": schema.StringAttribute{
							MarkdownDescription: "Memory generation",
							Computed:            true,
						},
					},
				},
			},
			"storage": schema.ListNestedAttribute{
				MarkdownDescription: "Available storage hardware options",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Storage type",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *availableHardwareDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *availableHardwareDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AvailableHardwareDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Fetching available hardware from marketplace")

	// Call the API
	hardware, err := d.client.GetAvailableHardware()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read available hardware, got error: %s", err))
		return
	}

	tflog.Debug(ctx, "Retrieved available hardware", map[string]interface{}{
		"cpu_count":     len(hardware.Cpu),
		"memory_count":  len(hardware.Memory),
		"storage_count": len(hardware.Storage),
	})

	// Map CPU hardware
	data.Cpu = make([]CpuHardwareModel, len(hardware.Cpu))
	for i, cpu := range hardware.Cpu {
		data.Cpu[i] = CpuHardwareModel{
			Architecture: types.StringValue(cpu.Architecture),
			Manufacturer: types.StringValue(cpu.Manufacturer),
		}
	}

	// Map Memory hardware
	data.Memory = make([]MemoryHardwareModel, len(hardware.Memory))
	for i, mem := range hardware.Memory {
		data.Memory[i] = MemoryHardwareModel{
			Type:       types.StringValue(mem.Type),
			Generation: types.StringValue(mem.Generation),
		}
	}

	// Map Storage hardware
	data.Storage = make([]StorageHardwareModel, len(hardware.Storage))
	for i, storage := range hardware.Storage {
		data.Storage[i] = StorageHardwareModel{
			Type: types.StringValue(storage.Type),
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
