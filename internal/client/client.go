package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/fatih/structs"
)

type Client struct {
	APIKey       string
	SecretAPIKey string

	baseURL    *url.URL
	httpClient *http.Client
}

func New(APIKey, SecretAPIKey string) Client {
	return Client{
		APIKey:       APIKey,
		SecretAPIKey: SecretAPIKey,
		baseURL: &url.URL{
			Scheme: "https",
			Host:   "api.porkbun.com",
			Path:   "/api/json/v3",
		},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SetCustomBaseURL(customURL *url.URL) {
	c.baseURL = customURL
}

func (c *Client) SetCustomHTTPClient(customHTTPClient *http.Client) {
	c.httpClient = customHTTPClient
}

func (c *Client) Do(ctx context.Context, url *url.URL, requestBody, responseBuffer any) error {
	// turn the request body into a map[string]any so that we can easily add the api credentials
	body := make(map[string]any)
	if requestBody != nil {
		body = structs.Map(requestBody)
	}
	body["apikey"] = c.APIKey
	body["secretapikey"] = c.SecretAPIKey

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url.String(),
		bytes.NewReader(data),
	)
	if err != nil {
		return fmt.Errorf("porkbun: could not create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("porkbun: http request failed: %w", err)
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(responseBuffer)
	if err != nil {
		return fmt.Errorf("porkbun: decoding response body failed: %w", err)
	}

	return err
}
