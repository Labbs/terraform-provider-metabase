package metabase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	metabase_v0_50 "github.com/labbs/terraform-provider-metabase/metabase/v0_50"
	metabase_v0_51 "github.com/labbs/terraform-provider-metabase/metabase/v0_51"
)

type PermissionsMembership struct {
	ID             int  `json:"id" tfsdk:"id"`
	GroupID        int  `json:"group_id" tfsdk:"group_id"`
	UserID         int  `json:"user_id" tfsdk:"user_id"`
	IsGroupManager bool `json:"is_group_manager" tfsdk:"is_group_manager"`
}

type PermissionsGroupMembership struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	UserID       int    `json:"user_id"`
	MembershipID int    `json:"membership_id"`
}

type PermissionsMembershipResponse struct {
	MembershipID   int  `json:"membership_id" tfsdk:"id"`
	GroupID        int  `json:"group_id" tfsdk:"group_id"`
	UserID         int  `json:"user_id" tfsdk:"user_id"`
	IsGroupManager bool `json:"is_group_manager" tfsdk:"is_group_manager"`
}

// CreatePermissionsMembership creates a permissions membership based on the API version.
func CreatePermissionsMembership(ctx context.Context, client *Client, permissionsMembership PermissionsMembership) (PermissionsMembership, error) {
	switch client.GetVersion() {
	case "v0.50":
		createdPermissionsMembership, err := client.V0_50.Client.PostPermissionsMembership(ctx, metabase_v0_50.PostPermissionsMembershipJSONRequestBody{
			GroupId:        permissionsMembership.GroupID,
			UserId:         permissionsMembership.UserID,
			IsGroupManager: &permissionsMembership.IsGroupManager,
		})
		if err != nil {
			return PermissionsMembership{}, err
		}

		var permissionsGroupMemebershipResponse []PermissionsGroupMembership

		err = json.NewDecoder(createdPermissionsMembership.Body).Decode(&permissionsGroupMemebershipResponse)
		if err != nil {
			return PermissionsMembership{}, err
		}

		for _, membership := range permissionsGroupMemebershipResponse {
			if membership.UserID == permissionsMembership.UserID {

				return PermissionsMembership{
					ID:             membership.MembershipID,
					GroupID:        permissionsMembership.GroupID,
					UserID:         permissionsMembership.UserID,
					IsGroupManager: permissionsMembership.IsGroupManager,
				}, nil
			}
		}

		return PermissionsMembership{}, fmt.Errorf("could not find created membership list after creation")
	case "v0.51":
		createdPermissionsMembership, err := client.V0_51.Client.PostPermissionsMembership(ctx, metabase_v0_51.PostPermissionsMembershipJSONRequestBody{
			GroupId:        permissionsMembership.GroupID,
			UserId:         permissionsMembership.UserID,
			IsGroupManager: &permissionsMembership.IsGroupManager,
		})
		if err != nil {
			return PermissionsMembership{}, err
		}

		var permissionsGroupMemebershipResponse []PermissionsGroupMembership

		err = json.NewDecoder(createdPermissionsMembership.Body).Decode(&permissionsGroupMemebershipResponse)
		if err != nil {
			return PermissionsMembership{}, err
		}

		for _, membership := range permissionsGroupMemebershipResponse {
			if membership.UserID == permissionsMembership.UserID {

				return PermissionsMembership{
					ID:             membership.MembershipID,
					GroupID:        permissionsMembership.GroupID,
					UserID:         permissionsMembership.UserID,
					IsGroupManager: permissionsMembership.IsGroupManager,
				}, nil
			}
		}

		return PermissionsMembership{}, fmt.Errorf("could not find created membership list after creation")
	default:
		return PermissionsMembership{}, fmt.Errorf("unsupported API version")
	}
}

// UpdatePermissionsMembership updates a permissions membership based on the API version.
func UpdatePermissionsMembership(ctx context.Context, client *Client, permissionsMembership PermissionsMembership) (PermissionsMembership, error) {
	switch client.GetVersion() {
	case "v0.50":
		updatedPermissionsMembership, err := client.V0_50.Client.PutPermissionsMembershipId(ctx, permissionsMembership.ID, metabase_v0_50.PutPermissionsMembershipIdJSONRequestBody{
			IsGroupManager: permissionsMembership.IsGroupManager,
		})
		if err != nil && updatedPermissionsMembership.StatusCode != 402 {
			return PermissionsMembership{}, err
		}

		if updatedPermissionsMembership.StatusCode == 402 {
			return PermissionsMembership{}, fmt.Errorf("please enable the Metabase Pro license to use this feature")
		}

		resp, err := metabase_v0_50.ParsePutPermissionsMembershipIdResponse(updatedPermissionsMembership)
		if err != nil {
			return PermissionsMembership{}, err
		}

		var permissionsMembershipResponse PermissionsMembership
		err = json.Unmarshal(resp.Body, &permissionsMembershipResponse)
		if err != nil {
			return PermissionsMembership{}, err
		}

		return permissionsMembershipResponse, nil
	case "v0.51":
		updatedPermissionsMembership, err := client.V0_51.Client.PutPermissionsMembershipId(ctx, permissionsMembership.ID, metabase_v0_51.PutPermissionsMembershipIdJSONRequestBody{
			IsGroupManager: permissionsMembership.IsGroupManager,
		})
		if err != nil && updatedPermissionsMembership.StatusCode != 402 {
			return PermissionsMembership{}, err
		}

		if updatedPermissionsMembership.StatusCode == 402 {
			return PermissionsMembership{}, fmt.Errorf("please enable the Metabase Pro license to use this feature")
		}

		resp, err := metabase_v0_51.ParsePutPermissionsMembershipIdResponse(updatedPermissionsMembership)
		if err != nil {
			return PermissionsMembership{}, err
		}

		var permissionsMembershipResponse PermissionsMembership
		err = json.Unmarshal(resp.Body, &permissionsMembershipResponse)
		if err != nil {
			return PermissionsMembership{}, err
		}

		return permissionsMembershipResponse, nil
	default:
		return PermissionsMembership{}, fmt.Errorf("unsupported API version")
	}
}

// DeletePermissionsMembership deletes a permissions membership based on the API version.
func DeletePermissionsMembership(ctx context.Context, client *Client, permissionsMembershipID int) error {
	switch client.GetVersion() {
	case "v0.50":
		_, err := client.V0_50.Client.DeletePermissionsMembershipId(ctx, permissionsMembershipID)
		if err != nil {
			return err
		}

		return nil
	case "v0.51":
		_, err := client.V0_51.Client.DeletePermissionsMembershipId(ctx, permissionsMembershipID)
		if err != nil {
			return err
		}

		return nil
	default:
		return fmt.Errorf("unsupported API version")
	}
}

// GetPermissionsMembership retrieves a permissions membership based on the API version.
func GetPermissionsMembership(ctx context.Context, client *Client, membershipID, groupID, userID int) (PermissionsMembership, error) {
	switch client.GetVersion() {
	case "v0.50":
		permissionsMembership, err := client.V0_50.Client.GetPermissionsMembership(ctx)
		if err != nil {
			return PermissionsMembership{}, err
		}

		var permissionsMemebershipResponse map[string][]PermissionsMembershipResponse

		err = json.NewDecoder(permissionsMembership.Body).Decode(&permissionsMemebershipResponse)
		if err != nil {
			return PermissionsMembership{}, err
		}

		for _, membership := range permissionsMemebershipResponse[strconv.Itoa(userID)] {
			if membership.MembershipID == membershipID {
				return PermissionsMembership{
					ID:             membership.MembershipID,
					GroupID:        membership.GroupID,
					UserID:         membership.UserID,
					IsGroupManager: membership.IsGroupManager,
				}, nil
			}
		}

		return PermissionsMembership{}, fmt.Errorf("could not find membership with ID %d", membershipID)
	case "v0.51":
		permissionsMembership, err := client.V0_51.Client.GetPermissionsMembership(ctx)
		if err != nil {
			return PermissionsMembership{}, err
		}

		var permissionsMemebershipResponse map[string][]PermissionsMembershipResponse

		err = json.NewDecoder(permissionsMembership.Body).Decode(&permissionsMemebershipResponse)
		if err != nil {
			return PermissionsMembership{}, err
		}

		for _, membership := range permissionsMemebershipResponse[strconv.Itoa(userID)] {
			if membership.MembershipID == membershipID {
				return PermissionsMembership{
					ID:             membership.MembershipID,
					GroupID:        membership.GroupID,
					UserID:         membership.UserID,
					IsGroupManager: membership.IsGroupManager,
				}, nil
			}
		}

		return PermissionsMembership{}, fmt.Errorf("could not find membership with ID %d", membershipID)
	default:
		return PermissionsMembership{}, fmt.Errorf("unsupported API version")
	}
}
