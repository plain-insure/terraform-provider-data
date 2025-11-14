// Copyright (c) Plain Technologies Aps

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NotNullResource{}

func NewNotNullResource() resource.Resource {
	return &NotNullResource{}
}

// NotNullResource defines the resource implementation.
type NotNullResource struct{}

// NotNullResourceModel describes the resource data model.
type NotNullResourceModel struct {
	Value        types.String `tfsdk:"value"`
	DefaultValue types.String `tfsdk:"default_value"`
	Result       types.String `tfsdk:"result"`
	ID           types.String `tfsdk:"id"`
}

func (r *NotNullResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notnull"
}

func (r *NotNullResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "NotNull resource that returns a non-null value based on input value and default. It will maintain prior non-null value when input changes to null. Working like a sticky value with a default fallback for the initial run.",

		Attributes: map[string]schema.Attribute{
			"value": schema.StringAttribute{
				MarkdownDescription: "The primary value to use for result will keep last non-null value when changed to null",
				Optional:            true,
			},
			"default_value": schema.StringAttribute{
				MarkdownDescription: "The default value to use when value is null and no previous value exists",
				Optional:            true,
			},
			"result": schema.StringAttribute{
				MarkdownDescription: "The computed result - never null",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal identifier",
			},
		},
	}
}

func (r *NotNullResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NotNullResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the result value
	result := r.computeResult(&data, nil)
	data.Result = types.StringValue(result)

	// Set a static ID since this resource doesn't represent a real external resource
	data.ID = types.StringValue("notnull")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotNullResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NotNullResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotNullResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NotNullResourceModel
	var state NotNullResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read prior state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the result value using prior state
	result := r.computeResult(&data, &state)
	data.Result = types.StringValue(result)

	// Maintain the same ID
	data.ID = state.ID

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotNullResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NotNullResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Nothing to do for deletion since this is a logical resource
}

// computeResult determines the result value based on the inputs and prior state
// Logic:
// 1. If value is not null and not unknown and not empty string, use value
// 2. If value is null/unknown and we have prior state with a result (meaning value changed from non-null to null/unknown), use the prior result
// 3. If value is null/unknown and we have no prior state, use default_value
// 4. If all are null, use empty string
func (r *NotNullResource) computeResult(data *NotNullResourceModel, priorState *NotNullResourceModel) string {
	// If value is provided (not null and not unknown and not empty string), use it
	if !data.Value.IsNull() && !data.Value.IsUnknown() && data.Value.ValueString() != "" {
		return data.Value.ValueString()
	}

	// Value is null or empty, check if we have a stored state with a non-empty result
	// (meaning value changed from non-null to null/empty)
	if priorState != nil && !priorState.Result.IsNull() && priorState.Result.ValueString() != "" {
		// Value changed from something to null/unknown, return the previously stored result
		return priorState.Result.ValueString()
	}

	// No prior state or value is null/empty from the start, use default_value
	if !data.DefaultValue.IsNull() && data.DefaultValue.ValueString() != "" {
		return data.DefaultValue.ValueString()
	}

	// Everything is null or empty, return empty string
	return ""
}
