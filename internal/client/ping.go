package client

import "context"

type pingResponse struct {
	status
	YourIp string `json:"yourIp"`
}

func (c *Client) Ping(ctx context.Context) (string, error) {
	url := c.baseURL.JoinPath("ping")

	var response pingResponse
	err := c.do(ctx, url, nil, &response)
	if err != nil {
		return "", err
	}

	if response.status.failed() {
		return "", response.status
	}

	return response.YourIp, nil
}
