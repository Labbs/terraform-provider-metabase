package metabase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type MetabaseInfo struct {
	Info struct {
		Version string `json:"version"`
	} `json:"info"`
}

// getMetabaseVersion retrieves the Metabase API version.
func getMetabaseVersion(endpoint string) (string, error) {
	// Get the Metabase API version
	var metabaseInfo MetabaseInfo

	// Get the Metabase API version.
	resp, err := http.Get(endpoint + "/docs/openapi.json")
	if err != nil {
		return "", fmt.Errorf("failed to get Metabase API version: %w", err)
	}

	// Decode the response.
	if err := json.NewDecoder(resp.Body).Decode(&metabaseInfo); err != nil {
		return "", fmt.Errorf("failed to decode Metabase API version: %w", err)
	}

	version := strings.Join(strings.Split(metabaseInfo.Info.Version, ".")[:2], ".")

	return version, nil
}
