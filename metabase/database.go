package metabase

import (
	"context"
	"encoding/json"
	"fmt"

	metabase_v0_50 "github.com/labbs/terraform-provider-metabase/metabase/v0_50"
	metabase_v0_51 "github.com/labbs/terraform-provider-metabase/metabase/v0_51"
)

type Database struct {
	ID             int    `json:"id" tfsdk:"id"`
	Name           string `json:"name" tfsdk:"name"`
	Engine         string `json:"engine" tfsdk:"engine"`
	AutoRunQueries bool   `json:"auto_run_queries" tfsdk:"auto_run_queries"`
	IsOnDemand     bool   `json:"is_on_demand" tfsdk:"is_on_demand"`

	// PostgresqlDetails is the details for a Postgresql database.
	PostgresqlDetails PostgresqlDetails `json:"postgresql_details" tfsdk:"postgresql_details"`
}

type PostgresqlDetails struct {
	Host             string `json:"host" tfsdk:"host"`
	Port             int    `json:"port" tfsdk:"port"`
	Database         string `json:"database" tfsdk:"database"`
	User             string `json:"user" tfsdk:"user"`
	Password         string `json:"password" tfsdk:"password"`
	SchemaFilter     string `json:"schema_filter" tfsdk:"schema_filter"`
	SSL              bool   `json:"ssl" tfsdk:"ssl"`
	SSLMode          string `json:"ssl_mode" tfsdk:"ssl_mode"`
	SSLUseClientMode bool   `json:"ssl_use_client_mode" tfsdk:"ssl_use_client_mode"`
}

// CreateDatabase creates a database based on the API version.
func CreateDatabase(ctx context.Context, client *Client, database Database) (Database, error) {
	var isOnDemand interface{} = database.IsOnDemand
	var details map[string]interface{}
	data, _ := json.Marshal(database.PostgresqlDetails)
	json.Unmarshal(data, &details)

	switch client.GetVersion() {
	case "v0.50":
		createDatabase, err := client.V0_50.Client.PostDatabase(ctx, metabase_v0_50.PostDatabaseJSONRequestBody{
			Name:           database.Name,
			Engine:         database.Engine,
			AutoRunQueries: &database.AutoRunQueries,
			IsOnDemand:     &isOnDemand,
			Details:        details,
		})
		if err != nil {
			return Database{}, err
		}

		_, err = metabase_v0_50.ParsePostDatabaseResponse(createDatabase)
		if err != nil {
			return Database{}, err
		}

		return database, nil
	case "v0.51":
		createDatabase, err := client.V0_51.Client.PostDatabase(ctx, metabase_v0_51.PostDatabaseJSONRequestBody{
			Name:           database.Name,
			Engine:         database.Engine,
			AutoRunQueries: &database.AutoRunQueries,
			IsOnDemand:     &isOnDemand,
			Details:        details,
		})
		if err != nil {
			return Database{}, err
		}

		_, err = metabase_v0_51.ParsePostDatabaseResponse(createDatabase)
		if err != nil {
			return Database{}, err
		}

		return database, nil
	default:
		return Database{}, fmt.Errorf("unsupported client version")
	}
}
