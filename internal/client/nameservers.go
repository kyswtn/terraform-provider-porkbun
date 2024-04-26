package client

import "context"

type getNameserversResponse struct {
	status
	Ns []string `json:"ns"`
}

func (c *Client) GetNameservers(ctx context.Context, domain string) ([]string, error) {
	url := c.baseURL.JoinPath("domain", "getNs", domain)

	var response getNameserversResponse
	err := c.do(ctx, url, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.failed() {
		return nil, response.status
	}

	return response.Ns, nil
}

type updateNameserversPayload struct {
	Ns []string `json:"ns"`
}

func (c *Client) UpdateNameservers(ctx context.Context, domain string, Ns []string) error {
	url := c.baseURL.JoinPath("domain", "updateNs", domain)

	response := status{}
	err := c.do(ctx, url, updateNameserversPayload{Ns}, &response)

	if err != nil {
		return err
	}

	if response.failed() {
		return response
	}

	return nil
}
