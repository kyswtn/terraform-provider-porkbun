package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type apiKeys struct {
	APIKey       string `json:"apikey"`
	SecretAPIKey string `json:"secretapikey"`
}

type Client struct {
	apiKeys    *apiKeys
	baseURL    *url.URL
	httpClient *http.Client
}

func New(APIKey, SecretAPIKey string) Client {
	return Client{
		apiKeys: &apiKeys{
			APIKey,
			SecretAPIKey,
		},
		baseURL: &url.URL{
			Scheme: "https",
			Host:   "porkbun.com",
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

func (c *Client) do(ctx context.Context, url *url.URL, requestBody, responseBuffer interface{}) error {
	bodyMarshaled, err := marshalAndJoin(c.apiKeys, requestBody)
	if err != nil {
		return fmt.Errorf("marshaling request body failed: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url.String(),
		bytes.NewReader(bodyMarshaled),
	)
	if err != nil {
		return fmt.Errorf("creating request object failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	defer func() { err = resp.Body.Close() }()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body failed: %w", err)
	}

	err = json.Unmarshal(responseBody, responseBuffer)
	if err != nil {
		return fmt.Errorf("unmarshaling response body failed: %w", err)
	}

	return err
}
