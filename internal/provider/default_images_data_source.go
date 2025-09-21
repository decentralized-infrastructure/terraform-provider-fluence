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
var _ datasource.DataSource = &defaultImagesDataSource{}

func NewDefaultImagesDataSource() datasource.DataSource {
	return &defaultImagesDataSource{}
}

// defaultImagesDataSource defines the data source implementation.
type defaultImagesDataSource struct {
	client *fluenceapi.Client
}

// DefaultImagesDataSourceModel describes the data source data model.
type DefaultImagesDataSourceModel struct {
	Images []DefaultImageModel `tfsdk:"images"`
}

// DefaultImageModel represents a default OS image
type DefaultImageModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Distribution types.String `tfsdk:"distribution"`
	Slug         types.String `tfsdk:"slug"`
	DownloadUrl  types.String `tfsdk:"download_url"`
	Username     types.String `tfsdk:"username"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

func (d *defaultImagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_images"
}

func (d *defaultImagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetch list of default OS images",

		Attributes: map[string]schema.Attribute{
			"images": schema.ListNestedAttribute{
				MarkdownDescription: "List of available default OS images",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Image ID",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Image name",
							Computed:            true,
						},
						"distribution": schema.StringAttribute{
							MarkdownDescription: "OS distribution",
							Computed:            true,
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "Image slug identifier",
							Computed:            true,
						},
						"download_url": schema.StringAttribute{
							MarkdownDescription: "Image download URL",
							Computed:            true,
						},
						"username": schema.StringAttribute{
							MarkdownDescription: "Default username for the image",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "Image creation timestamp",
							Computed:            true,
						},
						"updated_at": schema.StringAttribute{
							MarkdownDescription: "Image last update timestamp",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *defaultImagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *defaultImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DefaultImagesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch default images from the API
	images, err := d.client.GetDefaultImages()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read default images, got error: %s", err))
		return
	}

	// Convert API response to Terraform model
	data.Images = []DefaultImageModel{}
	for _, img := range images {
		image := DefaultImageModel{
			Id:           types.StringValue(img.Id),
			Name:         types.StringValue(img.Name),
			Distribution: types.StringValue(img.Distribution),
			Slug:         types.StringValue(img.Slug),
			DownloadUrl:  types.StringValue(img.DownloadUrl),
			Username:     types.StringValue(img.Username),
			CreatedAt:    types.StringValue(img.CreatedAt),
			UpdatedAt:    types.StringValue(img.UpdatedAt),
		}

		data.Images = append(data.Images, image)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "read default images data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
