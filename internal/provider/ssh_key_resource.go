package provider

import (
	"context"
	"fmt"

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
var _ resource.Resource = &SshKeyResource{}
var _ resource.ResourceWithImportState = &SshKeyResource{}

func NewSshKeyResource() resource.Resource {
	return &SshKeyResource{}
}

// SshKeyResource defines the resource implementation.
type SshKeyResource struct {
	client *fluenceapi.Client
}

// SshKeyResourceModel describes the resource data model.
type SshKeyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	PublicKey   types.String `tfsdk:"public_key"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Algorithm   types.String `tfsdk:"algorithm"`
	Comment     types.String `tfsdk:"comment"`
	Active      types.Bool   `tfsdk:"active"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (r *SshKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (r *SshKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SSH Key resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SSH Key identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "SSH Key name (optional)",
				Optional:            true,
			},
			"fingerprint": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SSH Key fingerprint",
			},
			"algorithm": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SSH Key algorithm",
			},
			"comment": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SSH Key comment",
			},
			"active": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the SSH Key is active",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SSH Key creation time",
			},
			"public_key": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "SSH public key content",
			},
		},
	}
}

func (r *SshKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SshKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SshKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the public key directly from the plan since it's not in the model
	var publicKey types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("public_key"), &publicKey)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create SSH key using the client
	createReq := fluenceapi.AddSshKey{
		PublicKey: publicKey.ValueString(),
	}

	// Set name if provided
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		name := data.Name.ValueString()
		createReq.Name = name
	}

	sshKey, err := r.client.CreateSshKey(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create SSH key, got error: %s", err))
		return
	}

	// Map response to resource model
	// Use fingerprint as ID since it's unique and always returned
	data.ID = types.StringValue(sshKey.Fingerprint)
	if sshKey.Name != nil {
		data.Name = types.StringValue(*sshKey.Name)
	} else {
		data.Name = types.StringNull()
	}
	data.PublicKey = types.StringValue(sshKey.PublicKey)
	data.Fingerprint = types.StringValue(sshKey.Fingerprint)
	data.Algorithm = types.StringValue(sshKey.Algorithm)
	data.Comment = types.StringValue(sshKey.Comment)
	data.Active = types.BoolValue(sshKey.Active)
	data.CreatedAt = types.StringValue(sshKey.CreatedAt)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created SSH key resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SshKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SshKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get all SSH keys and find the one matching our ID
	sshKeys, err := r.client.ListSshKeys()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SSH keys, got error: %s", err))
		return
	}

	// Find the SSH key by fingerprint (which we use as ID)
	var foundKey *fluenceapi.SshKey
	for _, key := range sshKeys {
		if key.Fingerprint == data.ID.ValueString() {
			foundKey = &key
			break
		}
	}

	if foundKey == nil {
		// SSH key not found, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	// Update the model with the current data
	if foundKey.Name != nil {
		data.Name = types.StringValue(*foundKey.Name)
	} else {
		data.Name = types.StringNull()
	}
	data.PublicKey = types.StringValue(foundKey.PublicKey)
	data.Fingerprint = types.StringValue(foundKey.Fingerprint)
	data.Algorithm = types.StringValue(foundKey.Algorithm)
	data.Comment = types.StringValue(foundKey.Comment)
	data.Active = types.BoolValue(foundKey.Active)
	data.CreatedAt = types.StringValue(foundKey.CreatedAt)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SshKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// For now, SSH keys don't support update operations in the API
	// We'll need to implement this if the API supports it in the future
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"SSH key updates are not currently supported. Please recreate the resource.",
	)
}

func (r *SshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SshKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the SSH key using the fingerprint
	err := r.client.RemoveSshKey(data.Fingerprint.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete SSH key, got error: %s", err))
		return
	}
}

func (r *SshKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID field for import
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
