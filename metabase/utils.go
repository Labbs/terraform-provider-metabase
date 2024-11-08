package metabase

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TransformPostgresDetails(details PostgresqlDetails) map[string]attr.Value {
	var result = make(map[string]attr.Value)
	result["host"] = types.StringValue(details.Host.ValueString())

	if details.Port.ValueInt64() == 0 {
		result["port"] = types.Int64Value(5432)
	} else {
		result["port"] = types.Int64Value(details.Port.ValueInt64())
	}

	result["database"] = types.StringValue(details.Database.ValueString())

	result["user"] = types.StringValue(details.User.ValueString())

	result["password"] = types.StringValue(details.Password.ValueString())

	if details.SchemaFilter.IsNull() || details.SchemaFilter.ValueString() == "" {
		result["schema_filter"] = types.StringNull()
	} else {
		result["schema_filter"] = types.StringValue(details.SchemaFilter.ValueString())
	}

	if details.SSL.IsNull() || !details.SSL.ValueBool() {
		result["ssl"] = types.BoolNull()
	} else {
		result["ssl"] = types.BoolValue(details.SSL.ValueBool())
	}
	if details.SSLMode.IsNull() || details.SSLMode.ValueString() == "" {
		result["ssl_mode"] = types.StringNull()
	} else {
		result["ssl_mode"] = types.StringValue(details.SSLMode.ValueString())
	}

	if details.SSLUseClientMode.IsNull() || !details.SSLUseClientMode.ValueBool() {
		result["ssl_use_client_mode"] = types.BoolNull()
	} else {
		result["ssl_use_client_mode"] = types.BoolValue(details.SSLUseClientMode.ValueBool())
	}

	return result
}
