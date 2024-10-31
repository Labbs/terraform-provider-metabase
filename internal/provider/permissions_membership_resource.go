package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labbs/terraform-provider-metabase/metabase"
)

var _ resource.ResourceWithImportState = &PermissionsMembershipResource{}

func NewPermissionsMembershipResource() resource.Resource {
	return &PermissionsMembershipResource{
		name: "metabase_permissions_membership",
	}
}

type PermissionsMembershipResource struct {
	name   string
	client *metabase.Client
}

func (r *PermissionsMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Metabase Permissions Membership",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Permissions Membership Id",
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"group_id": schema.Int64Attribute{
				MarkdownDescription: "Group Id",
				Required:            true,
			},
			"user_id": schema.Int64Attribute{
				MarkdownDescription: "User Id",
				Required:            true,
			},
			"is_group_manager": schema.BoolAttribute{
				MarkdownDescription: "Is Group Manager",
				Optional:            true,
			},
		},
	}
}

func (r *PermissionsMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan struct {
		ID             types.Int64 `tfsdk:"id"`
		GroupID        types.Int64 `tfsdk:"group_id"`
		UserID         types.Int64 `tfsdk:"user_id"`
		IsGroupManager types.Bool  `tfsdk:"is_group_manager"`
	}

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	membership := metabase.PermissionsMembership{
		GroupID:        int(plan.GroupID.ValueInt64()),
		UserID:         int(plan.UserID.ValueInt64()),
		IsGroupManager: plan.IsGroupManager.ValueBool(),
	}

	createMembership, err := metabase.CreatePermissionsMembership(ctx, r.client, membership)
	if err != nil {
		resp.Diagnostics.AddError("failed to create membership", err.Error())
		return
	}

	plan.ID = types.Int64Value(int64(createMembership.ID))
	if createMembership.IsGroupManager {
		plan.IsGroupManager = types.BoolValue(true)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PermissionsMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state struct {
		ID             types.Int64 `tfsdk:"id"`
		GroupID        types.Int64 `tfsdk:"group_id"`
		UserID         types.Int64 `tfsdk:"user_id"`
		IsGroupManager types.Bool  `tfsdk:"is_group_manager"`
	}

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	membership, err := metabase.GetPermissionsMembership(ctx, r.client, int(state.ID.ValueInt64()), int(state.GroupID.ValueInt64()), int(state.UserID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("failed to get membership", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &membership)...)
}

func (r *PermissionsMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan struct {
		ID             types.Int64 `tfsdk:"id"`
		GroupID        types.Int64 `tfsdk:"group_id"`
		UserID         types.Int64 `tfsdk:"user_id"`
		IsGroupManager types.Bool  `tfsdk:"is_group_manager"`
	}

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	membership := metabase.PermissionsMembership{
		ID:             int(plan.ID.ValueInt64()),
		GroupID:        int(plan.GroupID.ValueInt64()),
		UserID:         int(plan.UserID.ValueInt64()),
		IsGroupManager: plan.IsGroupManager.ValueBool(),
	}

	updateMembership, err := metabase.UpdatePermissionsMembership(ctx, r.client, membership)
	if err != nil {
		resp.Diagnostics.AddError("failed to update membership", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updateMembership)...)
}

func (r *PermissionsMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state struct {
		ID             types.Int64 `tfsdk:"id"`
		GroupID        types.Int64 `tfsdk:"group_id"`
		UserID         types.Int64 `tfsdk:"user_id"`
		IsGroupManager types.Bool  `tfsdk:"is_group_manager"`
	}

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := metabase.DeletePermissionsMembership(ctx, r.client, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("failed to delete membership", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, nil)...)
}

func (r *PermissionsMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permissions_membership"
}

func (r *PermissionsMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state struct {
		ID             types.Int64 `tfsdk:"id"`
		GroupID        types.Int64 `tfsdk:"group_id"`
		UserID         types.Int64 `tfsdk:"user_id"`
		IsGroupManager types.Bool  `tfsdk:"is_group_manager"`
	}

	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	membership, err := metabase.GetPermissionsMembership(ctx, r.client, int(state.ID.ValueInt64()), int(state.GroupID.ValueInt64()), int(state.UserID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("failed to get membership", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &membership)...)
}

func (r *PermissionsMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
