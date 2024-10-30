package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
    client *CoveoClient
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
                Required:    true,
                Description: "The API key for authenticating with the Coveo API.",
            },
            "organization_id": schema.StringAttribute{
                Required:    true,
                Description: "The Coveo organization ID.",
            },
        },
    }
}

// CoveoClient is a simple client to interact with the Coveo API.
type CoveoClient struct {
    ApiKey         string
    OrganizationID string
    HttpClient     *http.Client
}

func NewCoveoClient(apiKey, organizationID string) *CoveoClient {
    return &CoveoClient{
        ApiKey:         apiKey,
        OrganizationID: organizationID,
        HttpClient:     &http.Client{},
    }
}
// DoRequest is a helper to make API requests and parse the response.
func (c *CoveoClient) DoRequest(method, endpoint string, body interface{}) ([]byte, error) {
    // Include the organization ID in the base URL
    baseUrl := fmt.Sprintf("https://api.cloud.coveo.com/push/v1/organizations/%s", c.OrganizationID)
    url := fmt.Sprintf("%s/%s", baseUrl, endpoint)

    var reqBody []byte
    var err error
    if body != nil {
        reqBody, err = json.Marshal(body)
        if err != nil {
            return nil, err
        }
    }
    req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("API request error: %s", resp.Status)
    }
    return ioutil.ReadAll(resp.Body)
}

// Configure prepares a Coveo API client for data sources and resources.
func (p *coveoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    // Retrieve provider configuration values.
    var config struct {
        ApiKey         string `tfsdk:"api_key"`
        OrganizationID string `tfsdk:"organization_id"`
    }


    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    if config.ApiKey == "" || config.OrganizationID == "" {
        resp.Diagnostics.AddError(
            "Missing Configuration",
            "Both the API key and organization ID are required to authenticate with the Coveo API.",
        )
        return
    }

    // Initialize client
    // client := NewCoveoClient(config.ApiKey, config.OrganizationID)
    // if client == nil {
    //     resp.Diagnostics.AddError("Client Initialization Error", "Failed to initialize Coveo client.")
    //     return
    // }
    p.client =  NewCoveoClient(config.ApiKey, config.OrganizationID)
    // Pass the client to resources
    // resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *coveoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return nil
}

func (p *coveoProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        func() resource.Resource { return NewCoveoIndexResource(p.client) },
        func() resource.Resource { return NewCoveoDocumentResource(p.client) },
    }
}
