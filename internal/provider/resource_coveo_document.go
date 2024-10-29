package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type CoveoDocumentResource struct {
    client *CoveoClient
}

func NewCoveoDocumentResource() resource.Resource {
    return &CoveoDocumentResource{}
}

// Metadata sets the resource type name.
func (r *CoveoDocumentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = "coveo_document"
}

// Schema defines the schema for the document resource.
func (r *CoveoDocumentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            // "id": schema.StringAttribute{
            //     Computed:    true,
            //     Description: "The ID of the Coveo document.",
            // },
            "title": schema.StringAttribute{
                Required:    true,
                Description: "The title of the document.",
            },
            "content": schema.StringAttribute{
                Required:    true,
                Description: "The main content of the document.",
            },
            "source_id": schema.StringAttribute{
                Required:    true,
                Description: "The source ID where the document will be stored.",
            },
        },
    }
}

// Create sends a request to create a document in Coveo.
func (r *CoveoDocumentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Verify client initialization
    if r.client == nil {
        resp.Diagnostics.AddError("Client Error", "The Coveo client was not properly initialized.")
        return
    }

    // Extract plan attributes
    var plan struct {
        Title    string `tfsdk:"title"`
        Content  string `tfsdk:"content"`
        SourceID string `tfsdk:"source_id"`
    }
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Define request payload
    requestBody := map[string]interface{}{
        "title":   plan.Title,
        "content": plan.Content,
    }

    // Define the endpoint for document creation in the specified source
    endpoint := fmt.Sprintf("sources/%s/documents", plan.SourceID)
    
    // Make API request
    body, err := r.client.DoRequest("POST", endpoint, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to create document: %s", err))
        return
    }

    // Parse response
    var responseBody map[string]interface{}
    err = json.Unmarshal(body, &responseBody)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", "Could not parse response from Coveo API.")
        return
    }

    // Extract and set document ID
    documentID, ok := responseBody["id"].(string)
    if !ok || documentID == "" {
        resp.Diagnostics.AddError("Invalid Response", "Coveo API response did not include a valid document ID.")
        return
    }
    diags = resp.State.SetAttribute(ctx, path.Root("id"), documentID)
    resp.Diagnostics.Append(diags...)

    // Optionally, set other response attributes if they are part of the API response
    diags = resp.State.SetAttribute(ctx, path.Root("title"), plan.Title)
    resp.Diagnostics.Append(diags...)
}


// Read retrieves the document’s data.
func (r *CoveoDocumentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state struct {
        ID       string `tfsdk:"id"`
        SourceID string `tfsdk:"source_id"`
    }
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    endpoint := fmt.Sprintf("sources/%s/documents/%s", state.SourceID, state.ID)
    body, err := r.client.DoRequest("GET", endpoint, nil)
    if err != nil {
        resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read document: %s", err))
        return
    }

    var responseBody map[string]interface{}
    err = json.Unmarshal(body, &responseBody)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", "Could not parse response from Coveo API.")
        return
    }

    resp.State.SetAttribute(ctx, path.Root("title"), responseBody["title"])
    resp.State.SetAttribute(ctx, path.Root("content"), responseBody["content"])
}

// Update modifies an existing document.
func (r *CoveoDocumentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan struct {
        ID       string `tfsdk:"id"`
        Title    string `tfsdk:"title"`
        Content  string `tfsdk:"content"`
        SourceID string `tfsdk:"source_id"`
    }
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    requestBody := map[string]interface{}{
        "title":   plan.Title,
        "content": plan.Content,
    }

    endpoint := fmt.Sprintf("sources/%s/documents/%s", plan.SourceID, plan.ID)
    _, err := r.client.DoRequest("PUT", endpoint, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to update document: %s", err))
        return
    }
}

// Delete removes a document.
func (r *CoveoDocumentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state struct {
        ID       string `tfsdk:"id"`
        SourceID string `tfsdk:"source_id"`
    }
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    endpoint := fmt.Sprintf("sources/%s/documents/%s", state.SourceID, state.ID)
    _, err := r.client.DoRequest("DELETE", endpoint, nil)
    if err != nil {
        resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to delete document: %s", err))
        return
    }

    resp.State.RemoveResource(ctx)
}
