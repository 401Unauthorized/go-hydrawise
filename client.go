package hydrawise

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client is the main struct for interacting with Hydrawise API.
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Hydrawise API client.
// If httpClient is nil, http.DefaultClient is used.
func NewClient(apiKey string, httpClient *http.Client) *Client {
	baseURL := "https://api.hydrawise.com/api/v1/"
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		APIKey:     apiKey,
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}
}

// doGet sends a GET request to Hydrawise API.
func (c *Client) doGet(endpoint string, params url.Values, resp interface{}) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("api_key", c.APIKey)
	fullURL := c.BaseURL + endpoint + "?" + params.Encode()
	httpResp, err := c.HTTPClient.Get(fullURL)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != 200 {
		return fmt.Errorf("api error: %s", body)
	}
	return json.Unmarshal(body, resp)
}
