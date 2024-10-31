package metabase

import (
	"context"
	"encoding/json"
	"fmt"

	metabase_v0_50 "github.com/labbs/terraform-provider-metabase/metabase/v0_50"
	metabase_v0_51 "github.com/labbs/terraform-provider-metabase/metabase/v0_51"
)

type User struct {
	ID        int    `json:"id" tfsdk:"id"`
	Email     string `json:"email" tfsdk:"email"`
	FirstName string `json:"first_name" tfsdk:"first_name"`
	LastName  string `json:"last_name" tfsdk:"last_name"`
}

// CreateUser creates a user based on the API version.
func CreateUser(ctx context.Context, client *Client, user User) (User, error) {
	switch client.GetVersion() {
	case "v0.50":
		createdUser, err := client.V0_50.Client.PostUser(ctx, metabase_v0_50.PostUserJSONRequestBody{
			Email:     user.Email,
			FirstName: &user.FirstName,
			LastName:  &user.LastName,
		})
		if err != nil {
			return User{}, err
		}

		resp, err := metabase_v0_50.ParsePostUserResponse(createdUser)
		if err != nil {
			return User{}, err
		}

		var userResponse User
		err = json.Unmarshal(resp.Body, &userResponse)
		if err != nil {
			return User{}, err
		}
		return userResponse, nil
	case "v0.51":
		createdUser, err := client.V0_51.Client.PostUser(ctx, metabase_v0_51.PostUserJSONRequestBody{
			Email:     user.Email,
			FirstName: &user.FirstName,
			LastName:  &user.LastName,
		})
		if err != nil {
			return User{}, err
		}

		resp, err := metabase_v0_51.ParsePostUserResponse(createdUser)
		if err != nil {
			return User{}, err
		}

		var userResponse User
		err = json.Unmarshal(resp.Body, &userResponse)
		if err != nil {
			return User{}, err
		}

		fmt.Println(userResponse)

		return userResponse, nil
	default:
		return User{}, fmt.Errorf("unsupported client version")
	}
}

// GetUser retrieves a user based on the API version.
func GetUser(ctx context.Context, client *Client, id int) (User, error) {
	switch client.GetVersion() {
	case "v0.50":
		user, err := client.V0_50.Client.GetUserId(ctx, id)
		if err != nil {
			return User{}, err
		}

		resp, err := metabase_v0_50.ParseGetUserIdResponse(user)
		if err != nil {
			return User{}, err
		}

		var userResponse User
		err = json.Unmarshal(resp.Body, &userResponse)
		if err != nil {
			return User{}, err
		}

		return userResponse, nil
	case "v0.51":
		user, err := client.V0_51.Client.GetUserId(ctx, id)
		if err != nil {
			return User{}, err
		}

		resp, err := metabase_v0_51.ParseGetUserIdResponse(user)
		if err != nil {
			return User{}, err
		}

		var userResponse User
		err = json.Unmarshal(resp.Body, &userResponse)
		if err != nil {
			return User{}, err
		}

		return userResponse, nil
	default:
		return User{}, fmt.Errorf("unsupported client version")
	}
}

// UpdateUser updates a user based on the API version.
func UpdateUser(ctx context.Context, client *Client, user User) (User, error) {
	switch client.GetVersion() {
	case "v0.50":
		updatedUser, err := client.V0_50.Client.PutUserId(ctx, user.ID, metabase_v0_50.PutUserIdJSONRequestBody{
			Email:     &user.Email,
			FirstName: &user.FirstName,
			LastName:  &user.LastName,
		})
		if err != nil {
			return User{}, err
		}

		resp, err := metabase_v0_50.ParsePutUserIdResponse(updatedUser)
		if err != nil {
			return User{}, err
		}

		var userResponse User
		err = json.Unmarshal(resp.Body, &userResponse)
		if err != nil {
			return User{}, err
		}

		return userResponse, nil
	case "v0.51":
		updatedUser, err := client.V0_51.Client.PutUserId(ctx, user.ID, metabase_v0_51.PutUserIdJSONRequestBody{
			Email:     &user.Email,
			FirstName: &user.FirstName,
			LastName:  &user.LastName,
		})
		if err != nil {
			return User{}, err
		}

		resp, err := metabase_v0_51.ParsePutUserIdResponse(updatedUser)
		if err != nil {
			return User{}, err
		}

		var userResponse User
		err = json.Unmarshal(resp.Body, &userResponse)
		if err != nil {
			return User{}, err
		}

		return userResponse, nil
	default:
		return User{}, fmt.Errorf("unsupported client version")
	}
}

// DeleteUser deletes a user based on the API version.
func DeleteUser(ctx context.Context, client *Client, id int) error {
	switch client.GetVersion() {
	case "v0.50":
		_, err := client.V0_50.Client.DeleteUserId(ctx, id)
		if err != nil {
			return err
		}

		return nil
	case "v0.51":
		_, err := client.V0_51.Client.DeleteUserId(ctx, id)
		if err != nil {
			return err
		}

		return nil
	default:
		return fmt.Errorf("unsupported client version")
	}
}
