package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewCoveoIndexResource(client *CoveoClient) resource.Resource {
    return &CoveoIndexResource{client: client}
}

type CoveoIndexResource struct {
    client *CoveoClient
}

func (r *CoveoIndexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = "coveo_index"
}

func (r *CoveoIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "name": schema.StringAttribute{
                Required: true,
                Description: "The name of the Coveo index.",
            },
            // Add other necessary attributes here
        },
    }
}

func (r *CoveoIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Extract the attributes from the Terraform configuration.
    if r.client == nil {
        resp.Diagnostics.AddError("Client Error", "The Coveo client was not properly initialized.")
        return
    }
	var plan struct {
        Name        string `tfsdk:"name"`
    }
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Define the request body for the API to create the index.
    requestBody := map[string]interface{}{
        "name":        plan.Name,
    }

    // Construct the URL for the Coveo index creation API.
    endpoint := "indexes"
    body, err := r.client.DoRequest("POST", endpoint, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to create Coveo index: %s", err))
        return
    }

    // Parse the response to extract the index ID.
    var responseBody map[string]interface{}
    err = json.Unmarshal(body, &responseBody)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", "Could not parse response from Coveo API.")
        return
    }

    // Get the index ID from the response.
    indexID, ok := responseBody["id"].(string)
    if !ok || indexID == "" {
        resp.Diagnostics.AddError("Invalid Response", "Coveo API response did not include a valid index ID.")
        return
    }

    // Update the Terraform state with the new index ID and name.
    diags = resp.State.SetAttribute(ctx, path.Root("id"), indexID)
    resp.Diagnostics.Append(diags...)
    diags = resp.State.SetAttribute(ctx, path.Root("name"), plan.Name)
    resp.Diagnostics.Append(diags...)
}




func (r *CoveoIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Extract the ID from the current state.
	var state struct {
		ID string `tfsdk:"id"`
	}
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Make the API request to get the index details.
	url := fmt.Sprintf("https://platform.cloud.coveo.com/rest/organizations/<org_id>/indexes/%s", state.ID)
	body, err := r.client.DoRequest("GET", url, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read Coveo index: %s", err))
		return
	}

	// Parse and set response attributes in the Terraform state.
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", "Could not parse response from Coveo API.")
		return
	}

	resp.State.SetAttribute(ctx, path.Root("name"), responseBody["name"])
}

func (r *CoveoIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Extract current state and planned changes.
	var plan, state struct {
		ID   string `tfsdk:"id"`
		Name string `tfsdk:"name"`
	}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Define request body with updated values.
	requestBody := map[string]interface{}{
		"name": plan.Name,
	}

	// Make the API request to update the index.
	url := fmt.Sprintf("https://platform.cloud.coveo.com/rest/organizations/<org_id>/indexes/%s", state.ID)
	_, err := r.client.DoRequest("PUT", url, requestBody)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to update Coveo index: %s", err))
		return
	}

	// Update state with new values.
	resp.State.SetAttribute(ctx, path.Root("name"), plan.Name)
}


func (r *CoveoIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Extract the ID from the current state.
	var state struct {
		ID string `tfsdk:"id"`
	}
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Make the API request to delete the index.
	url := fmt.Sprintf("https://platform.cloud.coveo.com/rest/organizations/<org_id>/indexes/%s", state.ID)
	_, err := r.client.DoRequest("DELETE", url, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to delete Coveo index: %s", err))
		return
	}

	// Remove the ID from Terraform state.
	resp.State.RemoveResource(ctx)
}
