package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ provider.Provider = &coveoProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
    return func() provider.Provider {
        return &coveoProvider{
            version: version,
        }
    }
}

// coveoProvider is the provider implementation.
type coveoProvider struct {
    version string
}

// Metadata returns the provider type name.
func (p *coveoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "coveo"
    resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *coveoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Required: true,
				Description: "The API key to authenticate with the Coveo API.",
			},
		},
	}
}

type CoveoClient struct {
	ApiKey string
	HttpClient *http.Client
}

func NewCoveoClient(apiKey string) *CoveoClient {
	return &CoveoClient{
		ApiKey: apiKey,
		HttpClient: &http.Client{},
	}
}

// Configure prepares a coveo API client for data sources and resources.
func (p *coveoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config struct {
		ApiKey string `tfsdk:"api_key"`
	}
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ApiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"An API key must be provided to authenticate with the Coveo API.",
		)
		return
	}

	client := NewCoveoClient(config.ApiKey)

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *coveoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return nil
}

// Resources defines the resources implemented in the provider.
func (p *coveoProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCoveoIndexResource,
	}
}
