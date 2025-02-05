package metabase

import (
	"context"
	"encoding/json"
	"fmt"

	metabase_v0_50 "github.com/labbs/terraform-provider-metabase/metabase/v0_50"
	metabase_v0_51 "github.com/labbs/terraform-provider-metabase/metabase/v0_51"
)

type PermissionsGraph struct {
	Revision int                                   `json:"revision" tfsdk:"revision"`
	Groups   map[int]map[int]PermissionsGraphGroup `json:"groups" tfsdk:"groups"`
}

type PermissionsGraphGroup struct {
	ViewData      string   `json:"view_data" tfsdk:"view_data"`
	CreateQueries string   `json:"create_queries" tfsdk:"create_queries"`
	Download      Download `json:"download" tfsdk:"download"`
}

type Download struct {
	Schemas string `json:"schemas" tfsdk:"schemas"`
}

func GetPermissionsGraph(ctx context.Context, client *Client, groupId int) (PermissionsGraph, error) {
	switch client.GetVersion() {
	case "v0.50":
		permissionsGraphGroup, err := client.V0_50.Client.GetPermissionsGraphGroupGroupId(ctx, groupId)
		if err != nil {
			return PermissionsGraph{}, err
		}

		resp, err := metabase_v0_50.ParseGetPermissionsGraphGroupGroupIdResponse(permissionsGraphGroup)
		if err != nil {
			return PermissionsGraph{}, err
		}

		var permissionsGraph PermissionsGraph
		err = json.Unmarshal(resp.Body, &permissionsGraph)
		if err != nil {
			return PermissionsGraph{}, err
		}

		return permissionsGraph, nil
	case "v0.51":
		permissionsGraphGroup, err := client.V0_50.Client.GetPermissionsGraphGroupGroupId(ctx, groupId)
		if err != nil {
			return PermissionsGraph{}, err
		}

		resp, err := metabase_v0_51.ParseGetPermissionsGraphGroupGroupIdResponse(permissionsGraphGroup)
		var permissionsGraph PermissionsGraph
		err = json.Unmarshal(resp.Body, &permissionsGraph)
		if err != nil {
			return PermissionsGraph{}, err
		}

		return permissionsGraph, nil
	default:
		return PermissionsGraph{}, fmt.Errorf("unsupported client version")
	}
}
