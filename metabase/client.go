package metabase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	metabase_v0_50 "github.com/labbs/terraform-provider-metabase/metabase/v0_50"
	metabase_v0_51 "github.com/labbs/terraform-provider-metabase/metabase/v0_51"
)

// ClientConfig contains the common configuration for all versions.
type ClientConfig struct {
	BaseURL  string
	Username string
	Password string
	APIKey   string
}

// MetabaseVersionInfo represents the API response for the version.
type MetabaseVersionInfo struct {
	Version string `json:"version"`
}

// VersionedClient is a generic client that works with all versions.
type VersionedClient[T any, O any, R any] struct {
	Client  T
	Version string
	Premium bool
}

// Concrete client that implements MetabaseAPI.
type Client struct {
	V0_50   *VersionedClient[metabase_v0_50.Client, metabase_v0_50.ClientOption, metabase_v0_50.RequestEditorFn]
	V0_51   *VersionedClient[metabase_v0_51.Client, metabase_v0_51.ClientOption, metabase_v0_51.RequestEditorFn]
	Version string
	Premium bool
}

// NewVersionedClient creates a new typed client for a specific version.
func NewVersionedClient[T any, O any, R any](
	config ClientConfig,
	version string,
	newClientFn func(string, ...O) (*T, error),
	withRequestEditorFn func(R) O,
) (*VersionedClient[T, O, R], error) {
	var typedAuthFn R

	// Basic authentication function.
	baseAuthFn := getAuthFunction(config)

	// Create a function with the correct type.
	switch any(typedAuthFn).(type) {
	case metabase_v0_50.RequestEditorFn:
		converted, ok := any(metabase_v0_50.RequestEditorFn(baseAuthFn)).(R)
		if !ok {
			return nil, fmt.Errorf("failed to convert to v0.50 RequestEditorFn")
		}
		typedAuthFn = converted
	case metabase_v0_51.RequestEditorFn:
		converted, ok := any(metabase_v0_51.RequestEditorFn(baseAuthFn)).(R)
		if !ok {
			return nil, fmt.Errorf("failed to convert to v0.51 RequestEditorFn")
		}
		typedAuthFn = converted
	default:
		return nil, fmt.Errorf("unsupported RequestEditorFn type")
	}

	// Create the client with the typed function.
	client, err := newClientFn(
		config.BaseURL,
		withRequestEditorFn(typedAuthFn),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return &VersionedClient[T, O, R]{
		Client:  *client,
		Version: version,
		Premium: false,
	}, nil
}

// getAuthFunction creates the appropriate authentication function.
func getAuthFunction(config ClientConfig) func(context.Context, *http.Request) error {
	if config.APIKey != "" {
		return func(ctx context.Context, req *http.Request) error {
			req.Header.Set("x-api-key", config.APIKey)
			return nil
		}
	}

	return func(ctx context.Context, req *http.Request) error {
		sessionReq := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{
			Username: config.Username,
			Password: config.Password,
		}

		jsonData, err := json.Marshal(sessionReq)
		if err != nil {
			return fmt.Errorf("error marshalling session request: %w", err)
		}

		baseURL := strings.TrimSuffix(req.URL.String(), req.URL.Path)
		sessionURL := fmt.Sprintf("%s/session", baseURL)
		sessReq, err := http.NewRequestWithContext(ctx, "POST", sessionURL, strings.NewReader(string(jsonData)))
		if err != nil {
			return fmt.Errorf("error creating session request: %w", err)
		}
		sessReq.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(sessReq)
		if err != nil {
			return fmt.Errorf("error during session request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("authentication error: status code %d", resp.StatusCode)
		}

		var sessionResp struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
			return fmt.Errorf("error decoding session response: %w", err)
		}

		req.Header.Set("X-Metabase-Session", sessionResp.ID)
		return nil
	}
}

// NewAutoVersionedClient automatically creates the correct client based on the API version.
func NewAutoVersionedClient(config ClientConfig) (*Client, error) {
	version, err := getMetabaseVersion(config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("error detecting version: %w", err)
	}

	client := &Client{
		Version: version,
	}

	switch version {
	case "v0.50":
		versionedClient, err := NewVersionedClient[
			metabase_v0_50.Client,
			metabase_v0_50.ClientOption,
			metabase_v0_50.RequestEditorFn,
		](config, version, metabase_v0_50.NewClient, metabase_v0_50.WithRequestEditorFn)
		if err != nil {
			return nil, err
		}
		client.V0_50 = versionedClient

	case "v0.51":
		versionedClient, err := NewVersionedClient[
			metabase_v0_51.Client,
			metabase_v0_51.ClientOption,
			metabase_v0_51.RequestEditorFn,
		](config, version, metabase_v0_51.NewClient, metabase_v0_51.WithRequestEditorFn)
		if err != nil {
			return nil, err
		}
		client.V0_51 = versionedClient

	default:
		return nil, fmt.Errorf("unsupported version: %s", version)
	}

	return client, nil
}

// GetClient returns the appropriate underlying client.
func (c *Client) GetClient() interface{} {
	switch c.Version {
	case "v0.50":
		return c.V0_50.Client
	case "v0.51":
		return c.V0_51.Client
	default:
		return nil
	}
}

// GetVersion returns the client version.
func (c *Client) GetVersion() string {
	return c.Version
}
