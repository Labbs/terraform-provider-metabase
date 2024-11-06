package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/labbs/terraform-provider-metabase/metabase"
)

// Ensure MetabaseProvider satisfies various provider interfaces.
var _ provider.Provider = &MetabaseProvider{}

// MetabaseProvider defines the provider implementation.
type MetabaseProvider struct {
	version string
	client  *metabase.Client
}

// MetabaseProviderModel describes the provider data model.
type MetabaseProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	ApiKey   types.String `tfsdk:"api_key"`
}

func (p *MetabaseProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "metabase"
	resp.Version = p.version
}

func (p *MetabaseProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Endpoint of the Metabase instance",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for the Metabase instance. Can also be set via the METABASE_USERNAME environment variable",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for the Metabase instance. Can also be set via the METABASE_PASSWORD environment variable",
				Optional:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for the Metabase instance. Can also be set via the METABASE_API_KEY environment variable",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *MetabaseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MetabaseProviderModel
	var clientConfig metabase.ClientConfig = metabase.ClientConfig{}

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check if the configuration values are set in the environment variables
	if data.Username.IsNull() {
		if username := os.Getenv("METABASE_USERNAME"); username != "" {
			data.Username = types.StringValue(username)
		}
	} else if data.Password.IsNull() {
		if password := os.Getenv("METABASE_PASSWORD"); password != "" {
			data.Password = types.StringValue(password)
		}
	}
	if data.ApiKey.IsNull() {
		if apiKey := os.Getenv("METABASE_API_KEY"); apiKey != "" {
			data.ApiKey = types.StringValue(apiKey)
		}
	}

	// Check if the configuration values are set
	if !data.Endpoint.IsNull() {
		clientConfig.BaseURL = data.Endpoint.ValueString()
	} else {
		resp.Diagnostics.AddError("endpoint is required", "")
		return
	}

	// Configuration values are now available.
	if !data.Username.IsNull() && !data.Password.IsNull() && !data.ApiKey.IsNull() {
		resp.Diagnostics.AddError("only one of username, password, or api_key can be set", "")
		return
	} else if data.Username.IsNull() && data.Password.IsNull() && data.ApiKey.IsNull() {
		resp.Diagnostics.AddError("one of username, password, or api_key must be set", "")
		return
	} else if !data.Username.IsNull() && !data.Password.IsNull() {
		clientConfig.Username = data.Username.ValueString()
		clientConfig.Password = data.Password.ValueString()
	} else if !data.ApiKey.IsNull() {
		clientConfig.APIKey = data.ApiKey.ValueString()
	} else {
		resp.Diagnostics.AddError("unexpected error", "An unexpected error occurred")
		return
	}

	// The client is automatically created with the correct version
	client, err := metabase.NewAutoVersionedClient(clientConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client creation failed",
			fmt.Sprintf("Impossible to connect to Metabase: %v", err),
		)
		return
	}

	p.client = client
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MetabaseProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserResource,
		NewPermissionsGroupResource,
		NewPermissionsMembershipResource,
	}
}

func (p *MetabaseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MetabaseProvider{
			version: version,
		}
	}
}
