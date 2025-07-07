package provider

import (
	"context"
	"os"

	fluenceapi "github.com/0xthresh/fluence-api-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &fluenceProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &fluenceProvider{
			version: version,
		}
	}
}

// fluenceProvider is the provider implementation.
type fluenceProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *fluenceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "fluence"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *fluenceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "The Fluence API host URL. Can also be set via the FLUENCE_HOST environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:    true, // Change from Required to Optional
				Sensitive:   true,
				Description: "The Fluence API key. Can also be set via the FLUENCE_API_KEY environment variable.",
			},
		},
	}
}

// fluenceProviderModel maps provider schema data to a Go type.
type fluenceProviderModel struct {
	Host   types.String `tfsdk:"host"`
	ApiKey types.String `tfsdk:"api_key"`
}

// Configure prepares a Fluence API client for data sources and resources.
func (p *fluenceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config fluenceProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Fluence API Host",
			"The provider cannot create the Fluence API client as there is an unknown configuration value for the Fluence API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FLUENCE_HOST environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Fluence API Key",
			"The provider cannot create the Fluence API client as there is an unknown configuration value for the Fluence API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FLUENCE_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("FLUENCE_HOST")
	api_key := os.Getenv("FLUENCE_API_KEY")

	// Set default host if not provided via env var or config
	if host == "" {
		host = "https://api.fluence.dev"
	}

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ApiKey.IsNull() {
		api_key = config.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Fluence API Host",
			"The provider cannot create the Fluence API client as there is a missing or empty value for the Fluence API host. "+
				"Set the host value in the configuration or use the FLUENCE_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Fluence API ApiKey",
			"The provider cannot create the Fluence API client as there is a missing or empty value for the Fluence API key. "+
				"Set the API key value in the configuration or use the FLUENCE_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Fluence client using the configuration values
	client, err := fluenceapi.NewClient(&host, &api_key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Fluence API Client",
			"An unexpected error occurred when creating the Fluence API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Fluence Client Error: "+err.Error(),
		)
		return
	}

	// Make the Fluence client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *fluenceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSshDataSource,
		NewVmsDataSource,
		NewBasicConfigurationsDataSource,
		NewAvailableCountriesDataSource,
		NewAvailableHardwareDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *fluenceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSshKeyResource,
		NewVmResource,
	}
}
