package metabase

import (
	"context"
	"encoding/json"
	"fmt"

	metabase_v0_50 "github.com/labbs/terraform-provider-metabase/metabase/v0_50"
	metabase_v0_51 "github.com/labbs/terraform-provider-metabase/metabase/v0_51"
)

type PermissionsGroup struct {
	ID   int    `json:"id" tfsdk:"id"`
	Name string `json:"name" tfsdk:"name"`
}

// CreateGroup creates a permissions group based on the API version.
func CreatePermissionsGroup(ctx context.Context, client *Client, permissionsGroup PermissionsGroup) (PermissionsGroup, error) {
	switch client.GetVersion() {
	case "v0.50":
		createdPermissionsGroup, err := client.V0_50.Client.PostPermissionsGroup(ctx, metabase_v0_50.PostPermissionsGroupJSONRequestBody{
			Name: permissionsGroup.Name,
		})
		if err != nil {
			return PermissionsGroup{}, err
		}

		resp, err := metabase_v0_50.ParsePostPermissionsGroupResponse(createdPermissionsGroup)
		if err != nil {
			return PermissionsGroup{}, err
		}

		var permissionsGroupResponse PermissionsGroup
		err = json.Unmarshal(resp.Body, &permissionsGroupResponse)
		if err != nil {
			return PermissionsGroup{}, err
		}

		if resp.StatusCode() != 200 {
			return PermissionsGroup{}, fmt.Errorf("error creating permissions group")
		}

		return permissionsGroupResponse, nil
	case "v0.51":
		createdPermissionsGroup, err := client.V0_51.Client.PostPermissionsGroup(ctx, metabase_v0_51.PostPermissionsGroupJSONRequestBody{
			Name: permissionsGroup.Name,
		})
		if err != nil {
			return PermissionsGroup{}, err
		}

		resp, err := metabase_v0_51.ParsePostPermissionsGroupResponse(createdPermissionsGroup)
		if err != nil {
			return PermissionsGroup{}, err
		}

		var permissionsGroupResponse PermissionsGroup
		err = json.Unmarshal(resp.Body, &permissionsGroupResponse)
		if err != nil {
			return PermissionsGroup{}, err
		}

		if resp.StatusCode() != 200 {
			return PermissionsGroup{}, fmt.Errorf("error creating permissions group")
		}

		return permissionsGroupResponse, nil
	default:
		return PermissionsGroup{}, fmt.Errorf("unsupported client version")
	}
}

// GetGroup returns a permissions group based on the API version.
func GetPermissionsGroup(ctx context.Context, client *Client, permissionsGroupID int) (PermissionsGroup, error) {
	switch client.GetVersion() {
	case "v0.50":
		permissionsGroup, err := client.V0_50.Client.GetPermissionsGroupId(ctx, permissionsGroupID)
		if err != nil {
			return PermissionsGroup{}, err
		}

		resp, err := metabase_v0_50.ParseGetPermissionsGroupResponse(permissionsGroup)
		if err != nil {
			return PermissionsGroup{}, err
		}

		var permissionsGroupResponse PermissionsGroup
		err = json.Unmarshal(resp.Body, &permissionsGroupResponse)
		if err != nil {
			return PermissionsGroup{}, err
		}

		if resp.StatusCode() != 200 {
			return PermissionsGroup{}, fmt.Errorf("error getting permissions group")
		}

		return permissionsGroupResponse, nil
	case "v0.51":
		permissionsGroup, err := client.V0_51.Client.GetPermissionsGroupId(ctx, permissionsGroupID)
		if err != nil {
			return PermissionsGroup{}, err
		}

		resp, err := metabase_v0_51.ParseGetPermissionsGroupResponse(permissionsGroup)
		if err != nil {
			return PermissionsGroup{}, err
		}

		var permissionsGroupResponse PermissionsGroup
		err = json.Unmarshal(resp.Body, &permissionsGroupResponse)
		if err != nil {
			return PermissionsGroup{}, err
		}

		if resp.StatusCode() != 200 {
			return PermissionsGroup{}, fmt.Errorf("error getting permissions group")
		}

		return permissionsGroupResponse, nil
	default:
		return PermissionsGroup{}, fmt.Errorf("unsupported client version")
	}
}

// UpdateGroup updates a group based on the API version.
func UpdatePermissionsGroup(ctx context.Context, client *Client, permissionsGroup PermissionsGroup) (PermissionsGroup, error) {
	switch client.GetVersion() {
	case "v0.50":
		updatedPermissionsGroup, err := client.V0_50.Client.PutPermissionsGroupGroupId(ctx, permissionsGroup.ID, metabase_v0_50.PutPermissionsGroupGroupIdJSONRequestBody{
			Name: permissionsGroup.Name,
		})
		if err != nil {
			return PermissionsGroup{}, err
		}

		resp, err := metabase_v0_50.ParseDeletePermissionsGroupGroupIdResponse(updatedPermissionsGroup)
		if err != nil {
			return PermissionsGroup{}, err
		}

		var permissionsGroupResponse PermissionsGroup
		err = json.Unmarshal(resp.Body, &permissionsGroupResponse)
		if err != nil {
			return PermissionsGroup{}, err
		}

		if resp.StatusCode() != 200 {
			return PermissionsGroup{}, fmt.Errorf("error updating permissions group")
		}

		return permissionsGroupResponse, nil
	case "v0.51":
		updatedPermissionsGroup, err := client.V0_51.Client.PutPermissionsGroupGroupId(ctx, permissionsGroup.ID, metabase_v0_51.PutPermissionsGroupGroupIdJSONRequestBody{
			Name: permissionsGroup.Name,
		})
		if err != nil {
			return PermissionsGroup{}, err
		}

		resp, err := metabase_v0_51.ParseDeletePermissionsGroupGroupIdResponse(updatedPermissionsGroup)
		if err != nil {
			return PermissionsGroup{}, err
		}

		var permissionsGroupResponse PermissionsGroup
		err = json.Unmarshal(resp.Body, &permissionsGroupResponse)
		if err != nil {
			return PermissionsGroup{}, err
		}

		if resp.StatusCode() != 200 {
			return PermissionsGroup{}, fmt.Errorf("error updating permissions group")
		}

		return permissionsGroupResponse, nil
	default:
		return PermissionsGroup{}, fmt.Errorf("unsupported client version")
	}
}

// DeleteGroup deletes a group based on the API version.
func DeletePermissionsGroup(ctx context.Context, client *Client, permissionsGroupID int) error {
	switch client.GetVersion() {
	case "v0.50":
		_, err := client.V0_50.Client.DeletePermissionsGroupGroupId(ctx, permissionsGroupID)
		if err != nil {
			return err
		}

		return nil
	case "v0.51":
		_, err := client.V0_51.Client.DeletePermissionsGroupGroupId(ctx, permissionsGroupID)
		if err != nil {
			return err
		}

		return nil
	default:
		return fmt.Errorf("unsupported client version")
	}
}
