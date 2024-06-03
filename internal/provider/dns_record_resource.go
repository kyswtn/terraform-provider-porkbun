package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	porkbun "github.com/kyswtn/terraform-provider-porkbun/internal/client"
	"github.com/kyswtn/terraform-provider-porkbun/internal/consts"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &DNSRecordResource{}
	_ resource.ResourceWithImportState = &DNSRecordResource{}
)

type DNSRecordResource struct {
	client *porkbun.Client
}

func NewDNSRecordResource() resource.Resource {
	return &DNSRecordResource{}
}

func (r *DNSRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *DNSRecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create DNS records for your domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the record.",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The domain of the record.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The subdomain for the record being created/updated/deleted, not including the domain itself. " +
					"Leave blank to target the root domain. Use * for a wildcard record.",
				Optional: true,
				Default:  stringdefault.StaticString(""),
				Computed: true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of record being created." +
					"Valid types are `A`, `MX`, `CNAME`, `ALIAS`, `TXT`, `NS`, `AAAA`, `SRV`, `TLSA`, and `CAA`",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("A", "MX", "CNAME", "ALIAS", "TXT", "NS", "AAAA", "SRV", "TLSA", "CAA"),
				},
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The answer content for the record. " +
					"Please see the DNS management popup from the domain management console for proper formatting of each record type.",
				Required: true,
			},
			"ttl": schema.Int64Attribute{
				MarkdownDescription: "The time to live in seconds for the record. The minimum and the default is 600 seconds.",
				Optional:            true,
				Default:             int64default.StaticInt64(600),
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(600),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The priority of the record for those that support it.",
				Optional:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Comments or notes about the DNS record. This field has no effect on DNS responses.",
				Optional:            true,
			},
		},
	}
}

type DNSRecordResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Domain   types.String `tfsdk:"domain"`
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	Content  types.String `tfsdk:"content"`
	TTL      types.Int64  `tfsdk:"ttl"`
	Priority types.Int64  `tfsdk:"priority"`
	Notes    types.String `tfsdk:"notes"`
}

func (r *DNSRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DNSRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DNSRecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	record := porkbun.DNSRecord{
		Name:    data.Name.ValueString(),
		Type:    data.Type.ValueString(),
		Content: data.Content.ValueString(),
		TTL:     strconv.Itoa(int(data.TTL.ValueInt64())),
	}

	priority := int(data.Priority.ValueInt64())
	if priority != 0 {
		record.Priority = strconv.Itoa(priority)
	}

	notes := data.Notes.ValueString()
	if notes != "" {
		record.Notes = data.Notes.ValueString()
	}

	ID, err := r.client.CreateDNSRecord(ctx, data.Domain.ValueString(), record)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create DNS record", err.Error())
		return
	}

	data.ID = types.StringValue(strconv.Itoa(ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DNSRecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()
	record, err := r.client.RetrieveDNSRecord(ctx, domain, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve DNS record", err.Error())
		return
	}

	name := record.Name
	if name == domain {
		name = ""
	} else {
		name = strings.ReplaceAll(name, fmt.Sprintf(".%s", domain), "")
	}

	data.Name = types.StringValue(name)
	data.Type = types.StringValue(record.Type)
	data.Content = types.StringValue(record.Content)

	ttl, _ := strconv.Atoi(record.TTL)
	data.TTL = types.Int64Value(int64(ttl))

	if record.Priority != "" {
		priority, _ := strconv.Atoi(record.Priority)
		data.Priority = types.Int64Value(int64(priority))
	}

	if record.Notes != "" {
		data.Notes = types.StringValue(record.Notes)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DNSRecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	record := porkbun.DNSRecord{
		Name:    data.Name.ValueString(),
		Type:    data.Type.ValueString(),
		Content: data.Content.ValueString(),
		TTL:     strconv.Itoa(int(data.TTL.ValueInt64())),
	}

	priority := int(data.Priority.ValueInt64())
	if priority != 0 {
		record.Priority = strconv.Itoa(priority)
	}

	notes := data.Notes.ValueString()
	if notes != "" {
		record.Notes = data.Notes.ValueString()
	}

	err := r.client.EditDNSRecord(ctx, data.Domain.ValueString(), data.ID.ValueString(), record)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update DNS record", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DNSRecordResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDNSRecord(ctx, data.Domain.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete DNS record", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importId := strings.SplitN(req.ID, "/", 2)

	var domain string
	var id string
	if len(importId) == 2 {
		domain, id = importId[0], importId[1]
	} else {
		resp.Diagnostics.AddError(
			"Invalid import ID specified",
			"Use the import ID of format \"FQDN/recordID\" to import DNS records.",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), domain)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
