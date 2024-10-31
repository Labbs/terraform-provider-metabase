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

var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{
		name: "metabase_user",
	}
}

type UserResource struct {
	name   string
	client *metabase.Client
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Metabase User",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "User Id",
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "User email",
				Required:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "User first name",
				Optional:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "User last name",
				Optional:            true,
			},
		},
	}
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan struct {
		ID        types.Int64  `tfsdk:"id"`
		Email     types.String `tfsdk:"email"`
		FirstName types.String `tfsdk:"first_name"`
		LastName  types.String `tfsdk:"last_name"`
	}

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := metabase.User{
		Email:     plan.Email.ValueString(),
		FirstName: plan.FirstName.ValueString(),
		LastName:  plan.LastName.ValueString(),
	}

	createdUser, err := metabase.CreateUser(ctx, r.client, user)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &createdUser)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state struct {
		ID        types.Int64  `tfsdk:"id"`
		Email     types.String `tfsdk:"email"`
		FirstName types.String `tfsdk:"first_name"`
		LastName  types.String `tfsdk:"last_name"`
	}

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := metabase.GetUser(ctx, r.client, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Read Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &user)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan struct {
		ID        types.Int64  `tfsdk:"id"`
		Email     types.String `tfsdk:"email"`
		FirstName types.String `tfsdk:"first_name"`
		LastName  types.String `tfsdk:"last_name"`
	}

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := metabase.User{
		ID:        int(plan.ID.ValueInt64()),
		Email:     plan.Email.ValueString(),
		FirstName: plan.FirstName.ValueString(),
		LastName:  plan.LastName.ValueString(),
	}

	updatedUser, err := metabase.UpdateUser(ctx, r.client, user)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedUser)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state struct {
		ID        types.Int64  `tfsdk:"id"`
		Email     types.String `tfsdk:"email"`
		FirstName types.String `tfsdk:"first_name"`
		LastName  types.String `tfsdk:"last_name"`
	}

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := metabase.DeleteUser(ctx, r.client, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", err.Error())
		return
	}
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
