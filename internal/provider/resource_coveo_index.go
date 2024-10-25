package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type CoveoIndexResource struct {
	client *CoveoClient
}

func NewCoveoIndexResource() resource.Resource {
	return &CoveoIndexResource{}
}


func (r *CoveoIndexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "coveo_index"
}

func (r *CoveoIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				Description: "The name of the index.",
			},
		},
	}
}

func (r *CoveoIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Implement the creation of an index using the Coveo API
}

func (r *CoveoIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Implement the read logic
}

func (r *CoveoIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Implement the update logic
}

func (r *CoveoIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Implement the delete logic
}