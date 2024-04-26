package provider

import (
	"context"
	"net/url"
	"os"
	"strconv"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	porkbun "github.com/kyswtn/terraform-provider-porkbun/internal/client"
	"github.com/kyswtn/terraform-provider-porkbun/internal/consts"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &PorkbunProvider{}

type PorkbunProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PorkbunProvider{
			version,
		}
	}
}

func (p *PorkbunProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "porkbun"
	resp.Version = p.version
}

func (p *PorkbunProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "`apikey` required by Porkbun API. " +
					"Can also be configured using the `PORKBUN_API_KEY` environment variable.",
				Required:  true,
				Sensitive: true,
			},
			"secret_api_key": schema.StringAttribute{
				MarkdownDescription: "`secretapikey` required by Porkbun API. " +
					"Can also be configured using the `PORKBUN_SECRET_API_KEY` environment variable.",
				Required:  true,
				Sensitive: true,
			},
			"custom_base_url": schema.StringAttribute{
				MarkdownDescription: "Override the default base URL (https://porkbun.com/api/json/v3) used by Porkbun API client. " +
					"Can also be configured using the `PORKBUN_CUSTOM_BASE_URL` environment variable.",
				Optional: true,
			},
			"max_retries": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of retries to perform when an API request fails (default to 4). " +
					"Can also be configured using the `PORKBUN_MAX_RETRIES` environment variable.",
				Optional: true,
			},
		},
	}
}

type PorkbunProviderConfigurationModel struct {
	APIKey        types.String `tfsdk:"api_key"`
	SecretAPIKey  types.String `tfsdk:"secret_api_key"`
	CustomBaseURL types.String `tfsdk:"custom_base_url"`
	MaxRetries    types.Int64  `tfsdk:"max_retries"`
}

func (p *PorkbunProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config PorkbunProviderConfigurationModel

	// Load configuration values into `config`.
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that none of the configuration values are unknown.
	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			consts.ErrUnknownConfigurationValue,
			`The provider cannot create Porkbun API client as there is an unknown configuration value for "api_key". `+
				"Either target apply the source of the value first, set the value statically in the configuration, "+
				"or use the PORKBUN_API_KEY environment variable.",
		)
	}
	if config.SecretAPIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_api_key"),
			consts.ErrUnknownConfigurationValue,
			`The provider cannot create Porkbun API client as there is an unknown configuration value for "secret_api_key". `+
				"Either target apply the source of the value first, set the value statically in the configuration, "+
				"or use the PORKBUN_SECRET_API_KEY environment variable.",
		)
	}
	if config.CustomBaseURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("custom_base_url"),
			consts.ErrUnknownConfigurationValue,
			`The provider cannot create Porkbun API client as there is an unknown configuration value for "custom_base_url". `+
				"Either target apply the source of the value first, set the value statically in the configuration, "+
				"or use the PORKBUN_CUSTOM_BASE_URL environment variable.",
		)
	}
	if config.MaxRetries.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("max_retries"),
			consts.ErrUnknownConfigurationValue,
			`The provider cannot create Porkbun API client as there is an unknown configuration value for "max_retries". `+
				"Either target apply the source of the value first, set the value statically in the configuration, "+
				"or use the PORKBUN_MAX_RETRIES environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Load configuration values from either terraform files or environment variables.
	var apiKey string
	if config.APIKey.IsNull() {
		apiKey = os.Getenv("PORKBUN_API_KEY")
	} else {
		apiKey = config.APIKey.ValueString()
	}

	var secretAPIKey string
	if config.SecretAPIKey.IsNull() {
		secretAPIKey = os.Getenv("PORKBUN_SECRET_API_KEY")
	} else {
		secretAPIKey = config.SecretAPIKey.ValueString()
	}

	var customBaseURL string
	if config.CustomBaseURL.IsNull() {
		customBaseURL = os.Getenv("PORKBUN_CUSTOM_BASE_URL")
	} else {
		customBaseURL = config.CustomBaseURL.ValueString()
	}

	var maxRetries int64 = 4
	if config.MaxRetries.IsNull() {
		if value, ok := os.LookupEnv("PORKBUN_MAX_RETRIES"); ok {
			maxRetries, _ = strconv.ParseInt(value, 10, 64)
		}
	} else {
		maxRetries = config.MaxRetries.ValueInt64()
	}

	// Validate that required values are populated.
	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			consts.ErrInvalidConfigurationValue,
			`The provider cannot create Porkbun API client as there is a missing or empty value for "api_key". `+
				"Set the value in the configuration or use the PORKBUN_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if secretAPIKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_api_key"),
			consts.ErrInvalidConfigurationValue,
			`The provider cannot create Porkbun API client as there is a missing or empty value for "secret_api_key". `+
				"Set the value in the configuration or use the PORKBUN_SECRET_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client := porkbun.New(apiKey, secretAPIKey)

	// If `customBaseURL` is set, parse it and set the value in client.
	if customBaseURL != "" {
		urlParsed, err := url.Parse(customBaseURL)
		if err != nil {
			if _, ok := os.LookupEnv("PORKBUN_CUSTOM_BASE_URL"); ok {
				resp.Diagnostics.AddError(
					consts.ErrInvalidConfigurationValue,
					"The provider cannot override the base URL of Porkbun API client as the value configured for "+
						"PORKBUN_CUSTOM_BASE_URL environment variable is not a valid URL.",
				)
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("custom_base_url"),
					consts.ErrInvalidConfigurationValue,
					`The provider cannot override the base URL of Porkbun API client as the value configured for "custom_base_url" is not a valid URL.`,
				)
			}
			return
		}
		client.SetCustomBaseURL(urlParsed)
	}

	// Replace client's `httpClient` with `retryablehttp.Client`.
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = int(maxRetries)
	client.SetCustomHTTPClient(retryClient.StandardClient())

	resp.ResourceData = &client
	resp.DataSourceData = &client
}

func (p *PorkbunProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNameserversResource,
		NewDNSRecordResource,
	}
}

func (p *PorkbunProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNameserversDataSource,
	}
}
