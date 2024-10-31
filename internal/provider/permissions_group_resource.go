package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labbs/terraform-provider-metabase/metabase"
)

var _ resource.ResourceWithImportState = &PermissionsGroupResource{}

func NewPermissionsGroupResource() resource.Resource {
	return &PermissionsGroupResource{
		name: "metabase_permissions_group",
	}
}

type PermissionsGroupResource struct {
	name   string
	client *metabase.Client
}

func (r *PermissionsGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Metabase Group",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Group Id",
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Group name",
				Required:            true,
			},
		},
	}
}

func (r *PermissionsGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan struct {
		ID   types.Int64  `tfsdk:"id"`
		Name types.String `tfsdk:"name"`
	}

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group := metabase.PermissionsGroup{
		Name: plan.Name.ValueString(),
	}

	createdGroup, err := metabase.CreatePermissionsGroup(ctx, r.client, group)
	if err != nil {
		resp.Diagnostics.AddError("failed to create group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &createdGroup)...)
}

func (r *PermissionsGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state struct {
		ID   types.Int64  `tfsdk:"id"`
		Name types.String `tfsdk:"name"`
	}

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := metabase.GetPermissionsGroup(ctx, r.client, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("failed to read group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &group)...)
}

func (r *PermissionsGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan struct {
		ID   types.Int64  `tfsdk:"id"`
		Name types.String `tfsdk:"name"`
	}

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group := metabase.PermissionsGroup{
		ID:   int(plan.ID.ValueInt64()),
		Name: plan.Name.ValueString(),
	}

	updatedGroup, err := metabase.UpdatePermissionsGroup(ctx, r.client, group)
	if err != nil {
		resp.Diagnostics.AddError("failed to update group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedGroup)...)
}

func (r *PermissionsGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state struct {
		ID   types.Int64  `tfsdk:"id"`
		Name types.String `tfsdk:"name"`
	}

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := metabase.DeletePermissionsGroup(ctx, r.client, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("failed to delete group", err.Error())
		return
	}
}

func (r *PermissionsGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permissions_group"
}

func (r *PermissionsGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PermissionsGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*metabase.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *metabase.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}
