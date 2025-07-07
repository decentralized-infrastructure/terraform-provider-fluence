package provider

import (
	"context"
	"fmt"

	fluenceapi "github.com/decentralized-infrastructure/fluence-api-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &basicConfigurationsDataSource{}

func NewBasicConfigurationsDataSource() datasource.DataSource {
	return &basicConfigurationsDataSource{}
}

// basicConfigurationsDataSource defines the data source implementation.
type basicConfigurationsDataSource struct {
	client *fluenceapi.Client
}

// BasicConfigurationsDataSourceModel describes the data source data model.
type BasicConfigurationsDataSourceModel struct {
	Configurations []types.String `tfsdk:"configurations"`
}

func (d *basicConfigurationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_basic_configurations"
}

func (d *basicConfigurationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetch available basic VM configurations from the marketplace",

		Attributes: map[string]schema.Attribute{
			"configurations": schema.ListAttribute{
				MarkdownDescription: "List of available basic VM configurations",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *basicConfigurationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *basicConfigurationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BasicConfigurationsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Fetching basic configurations from marketplace")

	// Call the API
	configs, err := d.client.GetBasicConfigurations()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read basic configurations, got error: %s", err))
		return
	}

	tflog.Debug(ctx, "Retrieved basic configurations", map[string]interface{}{
		"count":          len(configs),
		"configurations": configs,
	})

	// Map response to the model
	data.Configurations = make([]types.String, len(configs))
	for i, config := range configs {
		data.Configurations[i] = types.StringValue(config)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
