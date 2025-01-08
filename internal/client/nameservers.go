package client

import "context"

type GetNameserversResponse struct {
	Status
	Ns []string `json:"ns"`
}

func (c *Client) GetNameservers(ctx context.Context, domain string) ([]string, error) {
	url := c.baseURL.JoinPath("domain", "getNs", domain)

	var response GetNameserversResponse
	err := c.Do(ctx, url, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.HasFailed() {
		return nil, response.Status
	}

	return response.Ns, nil
}

type UpdateNameserversPayload struct {
	Ns []string `json:"ns"`
}

func (c *Client) UpdateNameservers(ctx context.Context, domain string, Ns []string) error {
	url := c.baseURL.JoinPath("domain", "updateNs", domain)

	response := Status{}
	err := c.Do(ctx, url, UpdateNameserversPayload{Ns}, &response)

	if err != nil {
		return err
	}

	if response.HasFailed() {
		return response
	}

	return nil
}
