package metabase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	metabase_v0_50 "github.com/labbs/terraform-provider-metabase/metabase/v0_50"
	metabase_v0_51 "github.com/labbs/terraform-provider-metabase/metabase/v0_51"
)

type Database struct {
	ID                types.Int64            `json:"id" tfsdk:"id"`
	Name              types.String           `json:"name" tfsdk:"name"`
	Engine            types.String           `json:"engine" tfsdk:"engine"`
	AutoRunQueries    types.Bool             `json:"auto_run_queries" tfsdk:"auto_run_queries"`
	IsOnDemand        types.Bool             `json:"is_on_demand" tfsdk:"is_on_demand"`
	PostgresqlDetails types.Object           `json:"-" tfsdk:"postgresql_details"`
	Details           map[string]interface{} `json:"details" tfsdk:"-"`
}

type PostgresqlDetails struct {
	Host             types.String `json:"host" tfsdk:"host"`
	Port             types.Int64  `json:"port" tfsdk:"port"`
	Database         types.String `json:"database" tfsdk:"database"`
	User             types.String `json:"user" tfsdk:"user"`
	Password         types.String `json:"password" tfsdk:"password"`
	SchemaFilter     types.String `json:"schema_filter" tfsdk:"schema_filter"`
	SSL              types.Bool   `json:"ssl" tfsdk:"ssl"`
	SSLMode          types.String `json:"ssl_mode" tfsdk:"ssl_mode"`
	SSLUseClientMode types.Bool   `json:"ssl_use_client_mode" tfsdk:"ssl_use_client_mode"`
}

var PostgresqlDetailsObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"host":                types.StringType,
		"port":                types.Int64Type,
		"database":            types.StringType,
		"user":                types.StringType,
		"password":            types.StringType,
		"schema_filter":       types.StringType,
		"ssl":                 types.BoolType,
		"ssl_mode":            types.StringType,
		"ssl_use_client_mode": types.BoolType,
	},
}

// CreateDatabase creates a database based on the API version.
func CreateDatabase(ctx context.Context, client *Client, database Database) (Database, error) {
	var isOnDemand interface{} = database.IsOnDemand
	var details map[string]interface{}

	autoRunQueries := database.AutoRunQueries.ValueBool()

	switch database.Engine.ValueString() {
	case "postgres":
		attrMap := database.PostgresqlDetails.Attributes()

		details = map[string]interface{}{
			"host":                attrMap["host"].(types.String).ValueString(),
			"port":                attrMap["port"].(types.Int64).ValueInt64(),
			"db":                  attrMap["database"].(types.String).ValueString(),
			"user":                attrMap["user"].(types.String).ValueString(),
			"password":            attrMap["password"].(types.String).ValueString(),
			"schema-filter":       attrMap["schema_filter"].(types.String).ValueString(),
			"ssl":                 attrMap["ssl"].(types.Bool).ValueBool(),
			"ssl-mode":            attrMap["ssl_mode"].(types.String).ValueString(),
			"ssl-use-client-mode": attrMap["ssl_use_client_mode"].(types.Bool).ValueBool(),
		}
	default:
		return Database{}, fmt.Errorf("unsupported database engine")
	}

	switch client.GetVersion() {
	case "v0.50":
		createDatabase, err := client.V0_50.Client.PostDatabase(ctx, metabase_v0_50.PostDatabaseJSONRequestBody{
			Name:           database.Name.ValueString(),
			Engine:         database.Engine.ValueString(),
			AutoRunQueries: &autoRunQueries,
			IsOnDemand:     &isOnDemand,
			Details:        details,
		})
		if err != nil {
			return Database{}, err
		}

		var databaseResponse map[string]interface{}
		err = json.NewDecoder(createDatabase.Body).Decode(&databaseResponse)
		if err != nil {
			return Database{}, err
		}

		return Database{
			ID: types.Int64Value(int64(databaseResponse["id"].(float64))),
		}, nil
	case "v0.51":
		createDatabase, err := client.V0_51.Client.PostDatabase(ctx, metabase_v0_51.PostDatabaseJSONRequestBody{
			Name:           database.Name.ValueString(),
			Engine:         database.Engine.ValueString(),
			AutoRunQueries: &autoRunQueries,
			IsOnDemand:     &isOnDemand,
			Details:        details,
		})
		if err != nil {
			return Database{}, err
		}

		var databaseResponse map[string]interface{}
		err = json.NewDecoder(createDatabase.Body).Decode(&databaseResponse)
		if err != nil {
			return Database{}, err
		}

		return Database{
			ID: types.Int64Value(int64(databaseResponse["id"].(float64))),
		}, nil
	default:
		return Database{}, fmt.Errorf("unsupported client version")
	}
}

// GetDatabase returns a database based on the API version.
func GetDatabase(ctx context.Context, client *Client, state Database) (Database, error) {
	var databaseResponse map[string]interface{}
	switch client.GetVersion() {
	case "v0.50":
		database, err := client.V0_50.Client.GetDatabaseId(ctx, int(state.ID.ValueInt64()), nil)
		if err != nil {
			return Database{}, err
		}

		resp, err := metabase_v0_50.ParseGetDatabaseIdResponse(database)
		if err != nil {
			return Database{}, err
		}
		err = json.Unmarshal(resp.Body, &databaseResponse)
		if err != nil {
			return Database{}, err
		}

	case "v0.51":
		database, err := client.V0_51.Client.GetDatabaseId(ctx, int(state.ID.ValueInt64()), nil)
		if err != nil {
			return Database{}, err
		}

		resp, err := metabase_v0_51.ParseGetDatabaseIdResponse(database)
		if err != nil {
			return Database{}, err
		}

		var databaseResponse Database
		err = json.Unmarshal(resp.Body, &databaseResponse)
		if err != nil {
			return Database{}, err
		}
	default:
		return Database{}, fmt.Errorf("unsupported client version")
	}

	switch state.Engine.ValueString() {
	case "postgres":
		var respDetails map[string]interface{} = databaseResponse["details"].(map[string]interface{})
		var postgresqlDetails PostgresqlDetails
		diag := state.PostgresqlDetails.As(ctx, &postgresqlDetails, basetypes.ObjectAsOptions{})
		if diag.HasError() {
			return Database{}, fmt.Errorf("failed to convert PostgresqlDetails")
		}

		t := TransformPostgresDetails(PostgresqlDetails{
			Host:             types.StringValue(respDetails["host"].(string)),
			Port:             types.Int64Value(int64(respDetails["port"].(float64))),
			Database:         types.StringValue(respDetails["db"].(string)),
			User:             types.StringValue(respDetails["user"].(string)),
			Password:         types.StringValue(respDetails["password"].(string)),
			SchemaFilter:     types.StringValue(respDetails["schema-filter"].(string)),
			SSL:              types.BoolValue(respDetails["ssl"].(bool)),
			SSLMode:          types.StringValue(respDetails["ssl-mode"].(string)),
			SSLUseClientMode: types.BoolValue(respDetails["ssl-use-client-mode"].(bool)),
		})

		// This test is necessary because the password is not returned in the response
		if postgresqlDetails.Password.ValueString() != respDetails["password"].(string) {
			t["password"] = postgresqlDetails.Password
		}

		details, objectDiags := types.ObjectValue(PostgresqlDetailsObjectType.AttrTypes, t)

		if objectDiags.HasError() {
			return Database{}, fmt.Errorf("failed to convert PostgresqlDetails")
		}

		state.PostgresqlDetails = details

		return state, nil
	default:
		return Database{}, fmt.Errorf("unsupported database engine")
	}
}

// UpdateDatabase updates a database based on the API version.
func UpdateDatabase(ctx context.Context, client *Client, database Database) (Database, error) {
	var engine interface{} = database.Engine.ValueString()
	var details map[string]interface{}

	var name string = database.Name.ValueString()
	var autoRunQueries bool = database.AutoRunQueries.ValueBool()
	var id int = int(database.ID.ValueInt64())

	switch database.Engine.ValueString() {
	case "postgres":
		attrMap := database.PostgresqlDetails.Attributes()

		log.Println(attrMap["password"].(types.String).ValueString())

		details = map[string]interface{}{
			"host":                attrMap["host"].(types.String).ValueString(),
			"port":                attrMap["port"].(types.Int64).ValueInt64(),
			"db":                  attrMap["database"].(types.String).ValueString(),
			"user":                attrMap["user"].(types.String).ValueString(),
			"password":            attrMap["password"].(types.String).ValueString(),
			"schema-filter":       attrMap["schema_filter"].(types.String).ValueString(),
			"ssl":                 attrMap["ssl"].(types.Bool).ValueBool(),
			"ssl-mode":            attrMap["ssl_mode"].(types.String).ValueString(),
			"ssl-use-client-mode": attrMap["ssl_use_client_mode"].(types.Bool).ValueBool(),
		}

	default:
		return Database{}, fmt.Errorf("unsupported database engine")
	}

	switch client.GetVersion() {
	case "v0.50":
		updatedDatabase, err := client.V0_50.Client.PutDatabaseId(ctx, id, metabase_v0_50.PutDatabaseIdJSONRequestBody{
			Name:           &name,
			Engine:         &engine,
			AutoRunQueries: &autoRunQueries,
			Details:        &details,
		})
		if err != nil {
			return Database{}, err
		}

		_, err = metabase_v0_50.ParsePutDatabaseIdResponse(updatedDatabase)
		if err != nil {
			return Database{}, err
		}

		return database, nil
	case "v0.51":
		updatedDatabase, err := client.V0_51.Client.PutDatabaseId(ctx, id, metabase_v0_51.PutDatabaseIdJSONRequestBody{
			Name:           &name,
			Engine:         &engine,
			AutoRunQueries: &autoRunQueries,
			Details:        &details,
		})
		if err != nil {
			return Database{}, err
		}

		_, err = metabase_v0_51.ParsePutDatabaseIdResponse(updatedDatabase)
		if err != nil {
			return Database{}, err
		}

		return database, nil
	default:
		return Database{}, fmt.Errorf("unsupported client version")
	}
}

// DeleteDatabase deletes a database based on the API version.
func DeleteDatabase(ctx context.Context, client *Client, databaseID int) error {
	switch client.GetVersion() {
	case "v0.50":
		_, err := client.V0_50.Client.DeleteDatabaseId(ctx, databaseID)
		if err != nil {
			return err
		}

		return nil
	case "v0.51":
		_, err := client.V0_51.Client.DeleteDatabaseId(ctx, databaseID)
		if err != nil {
			return err
		}

		return nil
	default:
		return fmt.Errorf("unsupported client version")
	}
}
