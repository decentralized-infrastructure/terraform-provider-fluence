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
var _ datasource.DataSource = &datacentersDataSource{}

func NewDatacentersDataSource() datasource.DataSource {
	return &datacentersDataSource{}
}

// datacentersDataSource defines the data source implementation.
type datacentersDataSource struct {
	client *fluenceapi.Client
}

// DatacentersDataSourceModel describes the data source data model.
type DatacentersDataSourceModel struct {
	Datacenters []DatacenterModel `tfsdk:"datacenters"`
}

// DatacenterModel represents a datacenter
type DatacenterModel struct {
	Id             types.String   `tfsdk:"id"`
	CountryCode    types.String   `tfsdk:"country_code"`
	CityCode       types.String   `tfsdk:"city_code"`
	Index          types.Int64    `tfsdk:"index"`
	Tier           types.Int64    `tfsdk:"tier"`
	Certifications []types.String `tfsdk:"certifications"`
	Slug           types.String   `tfsdk:"slug"`
}

func (d *datacentersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datacenters"
}

func (d *datacentersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetch list of registered datacenters",

		Attributes: map[string]schema.Attribute{
			"datacenters": schema.ListNestedAttribute{
				MarkdownDescription: "List of available datacenters",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Datacenter ID",
							Computed:            true,
						},
						"country_code": schema.StringAttribute{
							MarkdownDescription: "Country code",
							Computed:            true,
						},
						"city_code": schema.StringAttribute{
							MarkdownDescription: "City code",
							Computed:            true,
						},
						"index": schema.Int64Attribute{
							MarkdownDescription: "Datacenter index",
							Computed:            true,
						},
						"tier": schema.Int64Attribute{
							MarkdownDescription: "Datacenter tier",
							Computed:            true,
						},
						"certifications": schema.ListAttribute{
							MarkdownDescription: "List of datacenter certifications",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "Datacenter slug identifier",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *datacentersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *datacentersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatacentersDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch datacenters from the API
	datacenters, err := d.client.GetDatacenters()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read datacenters, got error: %s", err))
		return
	}

	// Convert API response to Terraform model
	data.Datacenters = []DatacenterModel{}
	for _, dc := range datacenters {
		certs := []types.String{}
		for _, cert := range dc.Certifications {
			certs = append(certs, types.StringValue(cert))
		}

		datacenter := DatacenterModel{
			Id:             types.StringValue(string(dc.Id)),
			CountryCode:    types.StringValue(dc.CountryCode),
			CityCode:       types.StringValue(dc.CityCode),
			Index:          types.Int64Value(dc.Index),
			Tier:           types.Int64Value(dc.Tier),
			Certifications: certs,
			Slug:           types.StringValue(dc.Slug),
		}

		data.Datacenters = append(data.Datacenters, datacenter)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "read datacenters data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
