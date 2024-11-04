package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/labbs/terraform-provider-metabase/metabase"
)

var _ resource.ResourceWithImportState = &DatabaseResource{}

func NewDatabaseResource() resource.Resource {
	return &DatabaseResource{
		name: "metabase_database",
	}
}

type DatabaseResource struct {
	name   string
	client *metabase.Client
}

func (r *DatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Metabase Database",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Database Id",
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Database name",
				Required:            true,
			},
			"engine": schema.StringAttribute{
				MarkdownDescription: "Database engine",
				Required:            true,
			},
			"auto_run_queries": schema.BoolAttribute{
				MarkdownDescription: "Auto run queries",
				Optional:            true,
			},
			"is_on_demand": schema.BoolAttribute{
				MarkdownDescription: "Is on demand",
				Optional:            true,
			},
			"postgresql_details": schema.SingleNestedAttribute{
				MarkdownDescription: "Postgresql configuration details",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						MarkdownDescription: "Database host",
						Required:            true,
					},
					"port": schema.Int64Attribute{
						MarkdownDescription: "Database port",
						Optional:            true,
					},
					"database": schema.StringAttribute{
						MarkdownDescription: "Database name",
						Required:            true,
					},
					"user": schema.StringAttribute{
						MarkdownDescription: "Database user",
						Required:            true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "Database password",
						Required:            true,
						Sensitive:           true,
					},
					"schema_filter": schema.StringAttribute{
						MarkdownDescription: "Database schema filter",
						Optional:            true,
					},
					"ssl": schema.BoolAttribute{
						MarkdownDescription: "Database ssl",
						Optional:            true,
					},
					"ssl_mode": schema.StringAttribute{
						MarkdownDescription: "Database ssl mode",
						Optional:            true,
					},
					"ssl_use_client_mode": schema.BoolAttribute{
						MarkdownDescription: "Database ssl use client mode",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan struct {
		ID                types.Int64  `tfsdk:"id"`
		Name              types.String `tfsdk:"name"`
		Engine            types.String `tfsdk:"engine"`
		AutoRunQueries    types.Bool   `tfsdk:"auto_run_queries"`
		IsOnDemand        types.Bool   `tfsdk:"is_on_demand"`
		PostgresqlDetails types.Object `tfsdk:"postgresql_details"`
	}

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	database := metabase.Database{
		Name:           plan.Name.ValueString(),
		Engine:         plan.Engine.ValueString(),
		AutoRunQueries: plan.AutoRunQueries.ValueBool(),
		IsOnDemand:     plan.IsOnDemand.ValueBool(),
	}

	switch plan.Engine.ValueString() {
	case "postgres":
		var postgresDetails metabase.PostgresqlDetails
		diag := plan.PostgresqlDetails.As(ctx, &postgresDetails, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}

		database.PostgresqlDetails = postgresDetails
		if database.PostgresqlDetails.Port == 0 {
			database.PostgresqlDetails.Port = 5432
		}
	default:
		resp.Diagnostics.AddError("unsupported database engine", "this database engine is not supported")
		return
	}

	createdDatabase, err := metabase.CreateDatabase(ctx, r.client, database)
	if err != nil {
		resp.Diagnostics.AddError("failed to create database", err.Error())
		return
	}

	plan.ID = types.Int64Value(int64(createdDatabase.ID))
	if createdDatabase.AutoRunQueries {
		plan.AutoRunQueries = types.BoolValue(true)
	}
	if createdDatabase.IsOnDemand {
		plan.IsOnDemand = types.BoolValue(true)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *DatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}

func (r *DatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
