package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	porkbun "github.com/kyswtn/terraform-provider-porkbun/internal/client"
	"github.com/kyswtn/terraform-provider-porkbun/internal/consts"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &NameserversDataSource{}
	_ datasource.DataSourceWithConfigure = &NameserversDataSource{}
)

type NameserversDataSource struct {
	client *porkbun.Client
}

func NewNameserversDataSource() datasource.DataSource {
	return &NameserversDataSource{}
}

func (d *NameserversDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nameservers"
}

func (d *NameserversDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get nameservers for your domain.",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				MarkdownDescription: "The FQDN of the domain.",
				Required:            true,
			},
			"nameservers": schema.ListAttribute{
				MarkdownDescription: "An array of nameserver host names.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

type NameserversDataSourceModel struct {
	Domain      types.String   `tfsdk:"domain"`
	Nameservers []types.String `tfsdk:"nameservers"`
}

func (d *NameserversDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*porkbun.Client)
	if !ok {
		resp.Diagnostics.AddError(
			consts.ErrUnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *porkbun.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *NameserversDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state NameserversDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()
	nameservers, err := d.client.GetNameservers(ctx, domain)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get nameservers", err.Error())
		return
	}

	tfNameservers := make([]types.String, len(nameservers))
	for i, ns := range nameservers {
		tfNameservers[i] = types.StringValue(ns)
	}
	state.Nameservers = tfNameservers

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
