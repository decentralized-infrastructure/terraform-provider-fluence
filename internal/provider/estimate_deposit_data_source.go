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
var _ datasource.DataSource = &estimateDepositDataSource{}

func NewEstimateDepositDataSource() datasource.DataSource {
	return &estimateDepositDataSource{}
}

// estimateDepositDataSource defines the data source implementation.
type estimateDepositDataSource struct {
	client *fluenceapi.Client
}

// EstimateDepositDataSourceModel describes the data source data model.
type EstimateDepositDataSourceModel struct {
	// Input parameters
	Instances   types.Int64                      `tfsdk:"instances"`
	Constraints *EstimateDepositConstraintsModel `tfsdk:"constraints"`

	// Output results
	DepositAmountUsdc  types.String `tfsdk:"deposit_amount_usdc"`
	DepositEpochs      types.Int64  `tfsdk:"deposit_epochs"`
	TotalPricePerEpoch types.String `tfsdk:"total_price_per_epoch"`
	MaxPricePerEpoch   types.String `tfsdk:"max_price_per_epoch"`
}

// EstimateDepositConstraintsModel represents the constraints for deposit estimation
type EstimateDepositConstraintsModel struct {
	BasicConfiguration       types.String   `tfsdk:"basic_configuration"`
	MaxTotalPricePerEpochUsd types.String   `tfsdk:"max_total_price_per_epoch_usd"`
	DatacenterCountries      []types.String `tfsdk:"datacenter_countries"`

	// Hardware constraints (optional)
	CpuArchitecture  []types.String `tfsdk:"cpu_architecture"`
	CpuManufacturer  []types.String `tfsdk:"cpu_manufacturer"`
	MemoryType       []types.String `tfsdk:"memory_type"`
	MemoryGeneration []types.String `tfsdk:"memory_generation"`
	StorageType      []types.String `tfsdk:"storage_type"`
}

func (d *estimateDepositDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm_estimate_deposit"
}

func (d *estimateDepositDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Estimate the deposit required for creating VMs with given configuration and constraints",

		Attributes: map[string]schema.Attribute{
			"instances": schema.Int64Attribute{
				MarkdownDescription: "Number of VM instances to estimate for",
				Required:            true,
			},
			"constraints": schema.SingleNestedAttribute{
				MarkdownDescription: "Constraints for the VM estimation",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"basic_configuration": schema.StringAttribute{
						MarkdownDescription: "Basic configuration constraint (e.g., 'small', 'medium', 'large')",
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
					"cpu_architecture": schema.ListAttribute{
						MarkdownDescription: "List of allowed CPU architectures",
						ElementType:         types.StringType,
						Optional:            true,
					},
					"cpu_manufacturer": schema.ListAttribute{
						MarkdownDescription: "List of allowed CPU manufacturers",
						ElementType:         types.StringType,
						Optional:            true,
					},
					"memory_type": schema.ListAttribute{
						MarkdownDescription: "List of allowed memory types",
						ElementType:         types.StringType,
						Optional:            true,
					},
					"memory_generation": schema.ListAttribute{
						MarkdownDescription: "List of allowed memory generations",
						ElementType:         types.StringType,
						Optional:            true,
					},
					"storage_type": schema.ListAttribute{
						MarkdownDescription: "List of allowed storage types",
						ElementType:         types.StringType,
						Optional:            true,
					},
				},
			},

			// Output attributes
			"deposit_amount_usdc": schema.StringAttribute{
				MarkdownDescription: "Required deposit amount in USDC",
				Computed:            true,
			},
			"deposit_epochs": schema.Int64Attribute{
				MarkdownDescription: "Number of epochs the deposit covers",
				Computed:            true,
			},
			"total_price_per_epoch": schema.StringAttribute{
				MarkdownDescription: "Total price per epoch for all instances",
				Computed:            true,
			},
			"max_price_per_epoch": schema.StringAttribute{
				MarkdownDescription: "Maximum price per epoch for all instances",
				Computed:            true,
			},
		},
	}
}

func (d *estimateDepositDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *estimateDepositDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EstimateDepositDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build the estimation request
	estimateRequest := fluenceapi.EstimateDepositRequestV3{
		Instances: int(data.Instances.ValueInt64()),
	}

	// Build constraints if provided
	if data.Constraints != nil {
		constraints := &fluenceapi.OfferConstraints{}

		// Basic configuration
		if !data.Constraints.BasicConfiguration.IsNull() {
			basicConfig := data.Constraints.BasicConfiguration.ValueString()
			constraints.BasicConfiguration = &basicConfig
		}

		// Max price per epoch
		if !data.Constraints.MaxTotalPricePerEpochUsd.IsNull() {
			maxPrice := data.Constraints.MaxTotalPricePerEpochUsd.ValueString()
			constraints.MaxTotalPricePerEpochUsd = &maxPrice
		}

		// Datacenter countries
		if len(data.Constraints.DatacenterCountries) > 0 {
			countries := make([]string, len(data.Constraints.DatacenterCountries))
			for i, country := range data.Constraints.DatacenterCountries {
				countries[i] = country.ValueString()
			}
			constraints.Datacenter = &fluenceapi.DatacenterConstraint{
				Countries: countries,
			}
		}

		// Hardware constraints
		hasHardwareConstraints := false
		hardwareConstraints := &fluenceapi.HardwareConstraints{}

		// CPU constraints
		if len(data.Constraints.CpuArchitecture) > 0 || len(data.Constraints.CpuManufacturer) > 0 {
			cpuHardware := []fluenceapi.CpuHardware{}

			// If both are specified, create combinations
			if len(data.Constraints.CpuArchitecture) > 0 && len(data.Constraints.CpuManufacturer) > 0 {
				for _, arch := range data.Constraints.CpuArchitecture {
					for _, mfr := range data.Constraints.CpuManufacturer {
						cpuHardware = append(cpuHardware, fluenceapi.CpuHardware{
							Architecture: arch.ValueString(),
							Manufacturer: mfr.ValueString(),
						})
					}
				}
			} else if len(data.Constraints.CpuArchitecture) > 0 {
				for _, arch := range data.Constraints.CpuArchitecture {
					cpuHardware = append(cpuHardware, fluenceapi.CpuHardware{
						Architecture: arch.ValueString(),
					})
				}
			} else {
				for _, mfr := range data.Constraints.CpuManufacturer {
					cpuHardware = append(cpuHardware, fluenceapi.CpuHardware{
						Manufacturer: mfr.ValueString(),
					})
				}
			}

			hardwareConstraints.Cpu = cpuHardware
			hasHardwareConstraints = true
		}

		// Memory constraints
		if len(data.Constraints.MemoryType) > 0 || len(data.Constraints.MemoryGeneration) > 0 {
			memoryHardware := []fluenceapi.MemoryHardware{}

			// If both are specified, create combinations
			if len(data.Constraints.MemoryType) > 0 && len(data.Constraints.MemoryGeneration) > 0 {
				for _, memType := range data.Constraints.MemoryType {
					for _, gen := range data.Constraints.MemoryGeneration {
						memoryHardware = append(memoryHardware, fluenceapi.MemoryHardware{
							Type:       memType.ValueString(),
							Generation: gen.ValueString(),
						})
					}
				}
			} else if len(data.Constraints.MemoryType) > 0 {
				for _, memType := range data.Constraints.MemoryType {
					memoryHardware = append(memoryHardware, fluenceapi.MemoryHardware{
						Type: memType.ValueString(),
					})
				}
			} else {
				for _, gen := range data.Constraints.MemoryGeneration {
					memoryHardware = append(memoryHardware, fluenceapi.MemoryHardware{
						Generation: gen.ValueString(),
					})
				}
			}

			hardwareConstraints.Memory = memoryHardware
			hasHardwareConstraints = true
		}

		// Storage constraints
		if len(data.Constraints.StorageType) > 0 {
			storageHardware := []fluenceapi.StorageHardware{}
			for _, storageType := range data.Constraints.StorageType {
				storageHardware = append(storageHardware, fluenceapi.StorageHardware{
					Type: storageType.ValueString(),
				})
			}
			hardwareConstraints.Storage = storageHardware
			hasHardwareConstraints = true
		}

		if hasHardwareConstraints {
			constraints.Hardware = hardwareConstraints
		}

		estimateRequest.Constraints = constraints
	}

	tflog.Debug(ctx, "Estimating VM deposit", map[string]interface{}{
		"instances":   estimateRequest.Instances,
		"constraints": estimateRequest.Constraints,
	})

	// Call the API
	estimate, err := d.client.EstimateDeposit(estimateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to estimate deposit, got error: %s", err))
		return
	}

	tflog.Debug(ctx, "Received deposit estimate", map[string]interface{}{
		"deposit_amount_usdc":   estimate.DepositAmountUsdc,
		"deposit_epochs":        estimate.DepositEpochs,
		"total_price_per_epoch": estimate.TotalPricePerEpoch,
		"max_price_per_epoch":   estimate.MaxPricePerEpoch,
	})

	// Map response to the model
	data.DepositAmountUsdc = types.StringValue(estimate.DepositAmountUsdc)
	data.DepositEpochs = types.Int64Value(int64(estimate.DepositEpochs))
	data.TotalPricePerEpoch = types.StringValue(estimate.TotalPricePerEpoch)
	data.MaxPricePerEpoch = types.StringValue(estimate.MaxPricePerEpoch)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
