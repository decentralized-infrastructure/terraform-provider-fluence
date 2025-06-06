package provider

import (
  "context"
  "fmt"

  "github.com/0xthresh/fluence-api-client-go"
  "github.com/hashicorp/terraform-plugin-framework/datasource"
  "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
  "github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
  _ datasource.DataSource = &sshDataSource{}
  //_ datasource.DataSourceWithConfigure = &sshDataSource{}
)

// Configure adds the provider configured client to the data source.
func (d *sshDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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


// NewSshDataSource is a helper function to simplify the provider implementation.
func NewSshDataSource() datasource.DataSource {
  return &sshDataSource{}
}


// sshDataSource is the data source implementation.
type sshDataSource struct {
  client *fluenceapi.Client
}


// Metadata returns the data source type name.
func (d *sshDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_ssh_keys"
}

// Schema defines the schema for the data source.
func (d *sshDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
  resp.Schema = schema.Schema{
    Attributes: map[string]schema.Attribute{
      "ssh_keys": schema.ListNestedAttribute{
        Computed: true,
        NestedObject: schema.NestedAttributeObject{
          Attributes: map[string]schema.Attribute{
            "name": schema.StringAttribute{
              Computed: true,
            },
            "fingerprint": schema.StringAttribute{
              Computed: true,
            },
            "algorithm": schema.StringAttribute{
              Computed: true,
            },
            "comment": schema.StringAttribute{
              Computed: true,
            },
            "public_key": schema.StringAttribute{
              Computed: true,
            },
			"active": schema.BoolAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
          },
        },
      },
    },
  }
}

//
type sshKeysDataSourceModel struct {
	SshKeys []sshKeysModel `tfsdk:"ssh_keys"`
}

// SshKey represents an SSH key object
type sshKeysModel struct {
    Name        types.String `tfsdk:"name"`
    Fingerprint types.String `tfsdk:"fingerprint"`
    Algorithm   types.String `tfsdk:"algorithm"`
    Comment     types.String `tfsdk:"comment"`
    PublicKey   types.String `tfsdk:"public_key"`
    Active      types.Bool   `tfsdk:"active"`
    CreatedAt   types.String `tfsdk:"created_at"`
}

// Read refreshes the Terraform state with the latest data.
func (d *sshDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var state sshKeysDataSourceModel

    sshKeys, err := d.client.ListSshKeys()
    if err != nil {
      resp.Diagnostics.AddError(
        "Unable to Read Fluence SSH Keys",
        err.Error(),
      )
      return
    }

    // Map response body to model
    for _, sshKey := range sshKeys {
      sshKeysState := sshKeysModel{
        Name:        types.StringValue(sshKey.Name),
        Fingerprint: types.StringValue(sshKey.Fingerprint),
        Algorithm:   types.StringValue(sshKey.Algorithm),
        Comment:     types.StringValue(sshKey.Comment),
        PublicKey:   types.StringValue(sshKey.PublicKey),
        Active:      types.BoolValue(sshKey.Active),
        CreatedAt:   types.StringValue(sshKey.CreatedAt),
      }

      state.SshKeys = append(state.SshKeys, sshKeysState)
    }

    // Set state
    diags := resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
      return
    }
}

