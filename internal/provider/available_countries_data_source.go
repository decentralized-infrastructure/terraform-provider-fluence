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
var _ datasource.DataSource = &availableCountriesDataSource{}

func NewAvailableCountriesDataSource() datasource.DataSource {
	return &availableCountriesDataSource{}
}

// availableCountriesDataSource defines the data source implementation.
type availableCountriesDataSource struct {
	client *fluenceapi.Client
}

// AvailableCountriesDataSourceModel describes the data source data model.
type AvailableCountriesDataSourceModel struct {
	Countries []types.String `tfsdk:"countries"`
}

func (d *availableCountriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_available_countries"
}

func (d *availableCountriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetch available datacenter countries from the marketplace",

		Attributes: map[string]schema.Attribute{
			"countries": schema.ListAttribute{
				MarkdownDescription: "List of available datacenter countries (country codes)",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *availableCountriesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *availableCountriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AvailableCountriesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Fetching available countries from marketplace")

	// Call the API
	countries, err := d.client.GetAvailableCountries()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read available countries, got error: %s", err))
		return
	}

	tflog.Debug(ctx, "Retrieved available countries", map[string]interface{}{
		"count":     len(countries),
		"countries": countries,
	})

	// Map response to the model
	data.Countries = make([]types.String, len(countries))
	for i, country := range countries {
		data.Countries[i] = types.StringValue(country)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
