package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	porkbun "github.com/kyswtn/terraform-provider-porkbun/internal/client"
	"github.com/kyswtn/terraform-provider-porkbun/internal/consts"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &NameserversResource{}
	_ resource.ResourceWithImportState = &NameserversResource{}
)

type NameserversResource struct {
	client *porkbun.Client
}

func NewNameserversResource() resource.Resource {
	return &NameserversResource{}
}

func (r *NameserversResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nameservers"
}

func (r *NameserversResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Update nameservers for your domain.",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				MarkdownDescription: "The FQDN of the domain.",
				Required:            true,
			},
			"nameservers": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "An array of nameservers that you would like to update your domain with.",
				Required:            true,
			},
		},
	}
}

type NameserversResourceModel struct {
	Domain      types.String   `tfsdk:"domain"`
	Nameservers []types.String `tfsdk:"nameservers"`
}

func (r *NameserversResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*porkbun.Client)
	if !ok {
		resp.Diagnostics.AddError(
			consts.ErrUnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *porkbun.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *NameserversResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NameserversResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nameservers := make([]string, len(data.Nameservers))
	for i, tfNameserver := range data.Nameservers {
		nameservers[i] = tfNameserver.ValueString()
	}

	err := r.client.UpdateNameservers(ctx, data.Domain.ValueString(), nameservers)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update nameservers", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NameserversResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NameserversResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nameservers, err := r.client.GetNameservers(ctx, data.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to get nameservers", err.Error())
		return
	}

	tfNameservers := make([]types.String, len(nameservers))
	for i, ns := range nameservers {
		tfNameservers[i] = types.StringValue(ns)
	}
	data.Nameservers = tfNameservers

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NameserversResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NameserversResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nameservers := make([]string, len(data.Nameservers))
	for i, tfNameserver := range data.Nameservers {
		nameservers[i] = tfNameserver.ValueString()
	}

	err := r.client.UpdateNameservers(ctx, data.Domain.ValueString(), nameservers)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update nameservers", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NameserversResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NameserversResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateNameservers(ctx, data.Domain.ValueString(), consts.GetDefaultNameservers())
	if err != nil {
		resp.Diagnostics.AddError("Unable to update nameservers", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NameserversResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}
