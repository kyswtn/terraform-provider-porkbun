package client

import "context"

type PingResponse struct {
	Status
	YourIp string `json:"yourIp"`
}

func (c *Client) Ping(ctx context.Context) (string, error) {
	url := c.baseURL.JoinPath("ping")

	var response PingResponse
	err := c.Do(ctx, url, nil, &response)
	if err != nil {
		return "", err
	}

	if response.Status.HasFailed() {
		return "", response.Status
	}

	return response.YourIp, nil
}
