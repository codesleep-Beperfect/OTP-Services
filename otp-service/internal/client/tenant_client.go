package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TenantClient struct {
	BaseURL string
}

func NewTenantClient(baseURL string) *TenantClient {
	return &TenantClient{BaseURL: baseURL}
}

func (c *TenantClient) Validate(apiKey string) (string, error) {
	url := fmt.Sprintf("%s/v1/tenant/validate?api_key=%s", c.BaseURL, apiKey)
	resp, err := http.Get(url)
	if err != nil {
    return "", fmt.Errorf("tenant service unavailable")
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid tenant")
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	id, ok := data["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response from tenant service")
	}
	return id, nil
}