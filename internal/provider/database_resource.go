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
			"mysql_details": schema.SingleNestedAttribute{
				MarkdownDescription: "Mysql configuration details",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						MarkdownDescription: "Database host",
						Required:            true,
					},
					"port": schema.Int64Attribute{
						MarkdownDescription: "Database port, default 3306",
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
				},
			},
		},
	}
}

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan metabase.Database

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	db := plan

	switch plan.Engine.ValueString() {
	case "postgres":
		if plan.PostgresqlDetails.IsNull() {
			resp.Diagnostics.AddError("missing postgres details", "postgres details are required")
			return
		}

		var postgresDetails metabase.PostgresqlDetails
		diag := plan.PostgresqlDetails.As(ctx, &postgresDetails, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}

		details, objectDiags := types.ObjectValue(metabase.PostgresqlDetailsObjectType.AttrTypes, metabase.TransformPostgresDetails(postgresDetails))

		resp.Diagnostics.Append(objectDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		db.PostgresqlDetails = details
	case "mysql":
		if plan.MysqlDetails.IsNull() {
			resp.Diagnostics.AddError("missing mysql details", "mysql details are required")
			return
		}

		var mysqlDetails metabase.MysqlDetails
		diag := plan.MysqlDetails.As(ctx, &mysqlDetails, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}

		details, objectDiags := types.ObjectValue(metabase.MysqlDetailsObjectType.AttrTypes, metabase.TransformMysqlDetails(mysqlDetails))

		resp.Diagnostics.Append(objectDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		db.MysqlDetails = details
	default:
		resp.Diagnostics.AddError("unsupported database engine", "this database engine is not supported")
		return
	}

	createdDatabase, err := metabase.CreateDatabase(ctx, r.client, db)
	if err != nil {
		resp.Diagnostics.AddError("failed to create database", err.Error())
		return
	}

	plan.ID = createdDatabase.ID

	// Mettre à jour les valeurs non-sensibles si nécessaires
	if createdDatabase.AutoRunQueries.ValueBool() {
		plan.AutoRunQueries = types.BoolValue(true)
	}
	if createdDatabase.IsOnDemand.ValueBool() {
		plan.IsOnDemand = types.BoolValue(true)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state metabase.Database

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	database, err := metabase.GetDatabase(ctx, r.client, state)
	if err != nil {
		resp.Diagnostics.AddError("failed to get database", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &database)...)
}

func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan metabase.Database

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	db := plan

	switch plan.Engine.ValueString() {
	case "postgres":
		var postgresDetails metabase.PostgresqlDetails
		diag := plan.PostgresqlDetails.As(ctx, &postgresDetails, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}

		if postgresDetails.Port.ValueInt64() == 0 {
			postgresDetails.Port = types.Int64Value(5432)
		}

		if postgresDetails.SSL.IsNull() {
			postgresDetails.SSL = types.BoolValue(false)
		}

		details, objectDiags := types.ObjectValue(metabase.PostgresqlDetailsObjectType.AttrTypes, metabase.TransformPostgresDetails(postgresDetails))

		resp.Diagnostics.Append(objectDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		db.PostgresqlDetails = details
	case "mysql":
		var mysqlDetails metabase.MysqlDetails
		diag := plan.MysqlDetails.As(ctx, &mysqlDetails, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}

		if mysqlDetails.Port.ValueInt64() == 0 {
			mysqlDetails.Port = types.Int64Value(3306)
		}

		details, objectDiags := types.ObjectValue(metabase.MysqlDetailsObjectType.AttrTypes, metabase.TransformMysqlDetails(mysqlDetails))

		resp.Diagnostics.Append(objectDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		db.MysqlDetails = details
	default:
		resp.Diagnostics.AddError("unsupported database engine", "this database engine is not supported")
		return
	}

	updatedDatabase, err := metabase.UpdateDatabase(ctx, r.client, db)
	if err != nil {
		resp.Diagnostics.AddError("failed to update database", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedDatabase)...)
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state metabase.Database

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := metabase.DeleteDatabase(ctx, r.client, int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("failed to delete database", err.Error())
		return
	}
}

func (r *DatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
