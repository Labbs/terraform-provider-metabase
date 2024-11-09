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

func TransformMysqlDetails(details MysqlDetails) map[string]attr.Value {
	var result = make(map[string]attr.Value)
	result["host"] = types.StringValue(details.Host.ValueString())

	if details.Port.ValueInt64() == 0 {
		result["port"] = types.Int64Value(3306)
	} else {
		result["port"] = types.Int64Value(details.Port.ValueInt64())
	}

	result["database"] = types.StringValue(details.Database.ValueString())

	result["user"] = types.StringValue(details.User.ValueString())

	result["password"] = types.StringValue(details.Password.ValueString())

	return result
}

func detailsPostgresAttribute(attrMap map[string]attr.Value) map[string]interface{} {
	details := make(map[string]interface{})

	if host, ok := attrMap["host"].(types.String); ok {
		details["host"] = host.ValueString()
	}
	if port, ok := attrMap["port"].(types.Int64); ok {
		details["port"] = port.ValueInt64()
	}
	if db, ok := attrMap["database"].(types.String); ok {
		details["db"] = db.ValueString()
	}
	if user, ok := attrMap["user"].(types.String); ok {
		details["user"] = user.ValueString()
	}
	if password, ok := attrMap["password"].(types.String); ok {
		details["password"] = password.ValueString()
	}
	if schemaFilter, ok := attrMap["schema_filter"].(types.String); ok {
		details["schema-filter"] = schemaFilter.ValueString()
	}
	if ssl, ok := attrMap["ssl"].(types.Bool); ok {
		details["ssl"] = ssl.ValueBool()
	}
	if sslMode, ok := attrMap["ssl_mode"].(types.String); ok {
		details["ssl-mode"] = sslMode.ValueString()
	}
	if sslUseClientMode, ok := attrMap["ssl_use_client_mode"].(types.Bool); ok {
		details["ssl-use-client-mode"] = sslUseClientMode.ValueBool()
	}

	return details
}

func detailsMysqlAttribute(attrMap map[string]attr.Value) map[string]interface{} {
	details := make(map[string]interface{})

	if host, ok := attrMap["host"].(types.String); ok {
		details["host"] = host.ValueString()
	}
	if port, ok := attrMap["port"].(types.Int64); ok {
		details["port"] = port.ValueInt64()
	}
	if db, ok := attrMap["database"].(types.String); ok {
		details["db"] = db.ValueString()
	}
	if user, ok := attrMap["user"].(types.String); ok {
		details["user"] = user.ValueString()
	}
	if password, ok := attrMap["password"].(types.String); ok {
		details["password"] = password.ValueString()
	}

	return details
}

func PostgresqlDetailsVerifyType(details map[string]interface{}) PostgresqlDetails {
	var result PostgresqlDetails

	if host, ok := details["host"].(string); ok {
		result.Host = types.StringValue(host)
	}
	if port, ok := details["port"].(int64); ok {
		result.Port = types.Int64Value(port)
	}
	if db, ok := details["db"].(string); ok {
		result.Database = types.StringValue(db)
	}
	if user, ok := details["user"].(string); ok {
		result.User = types.StringValue(user)
	}
	if password, ok := details["password"].(string); ok {
		result.Password = types.StringValue(password)
	}
	if schemaFilter, ok := details["schema-filter"].(string); ok {
		result.SchemaFilter = types.StringValue(schemaFilter)
	}
	if ssl, ok := details["ssl"].(bool); ok {
		result.SSL = types.BoolValue(ssl)
	}
	if sslMode, ok := details["ssl-mode"].(string); ok {
		result.SSLMode = types.StringValue(sslMode)
	}
	if sslUseClientMode, ok := details["ssl-use-client-mode"].(bool); ok {
		result.SSLUseClientMode = types.BoolValue(sslUseClientMode)
	}

	return result
}

func MysqlDetailsVerifyType(details map[string]interface{}) MysqlDetails {
	var result MysqlDetails

	if host, ok := details["host"].(string); ok {
		result.Host = types.StringValue(host)
	}
	if port, ok := details["port"].(int64); ok {
		result.Port = types.Int64Value(port)
	}
	if db, ok := details["db"].(string); ok {
		result.Database = types.StringValue(db)
	}
	if user, ok := details["user"].(string); ok {
		result.User = types.StringValue(user)
	}
	if password, ok := details["password"].(string); ok {
		result.Password = types.StringValue(password)
	}

	return result
}
