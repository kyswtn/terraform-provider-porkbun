package client

import (
	"context"
	"fmt"
)

// All the fields are `omitempty` so that the same struct can be used both as input type for creating records
// and as return type from reading records.
type DNSRecord struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Content  string `json:"content,omitempty"`
	TTL      string `json:"ttl,omitempty"`
	Priority string `json:"prio,omitempty"`
	Notes    string `json:"notes,omitempty"`
}

type createDNSRecordResponse struct {
	status
	// Porkbun's API will return an `int` upon creation, but it'll be a string when retrieved.
	ID int `json:"id"`
}

func (c *Client) CreateDNSRecord(ctx context.Context, domain string, record DNSRecord) (int, error) {
	url := c.baseURL.JoinPath("dns", "create", domain)

	var response createDNSRecordResponse
	err := c.do(ctx, url, record, &response)

	if err != nil {
		return 0, err
	}

	if response.failed() {
		return 0, err
	}

	return response.ID, nil
}

type retrieveDNSRecordResponse struct {
	status
	Records []DNSRecord `json:"records"`
}

func (c *Client) RetrieveDNSRecord(ctx context.Context, domain, id string) (DNSRecord, error) {
	url := c.baseURL.JoinPath("dns", "retrieve", domain, id)

	var response retrieveDNSRecordResponse
	err := c.do(ctx, url, nil, &response)

	if err != nil {
		return DNSRecord{}, err
	}

	if response.failed() {
		return DNSRecord{}, response.status
	}

	if len(response.Records) < 1 {
		return DNSRecord{}, fmt.Errorf("DNS record not found")
	}

	return response.Records[0], nil
}

func (c *Client) EditDNSRecord(ctx context.Context, domain, id string, record DNSRecord) error {
	url := c.baseURL.JoinPath("dns", "edit", domain, id)

	response := status{}
	err := c.do(ctx, url, record, &response)

	if err != nil {
		return err
	}

	if response.failed() {
		return response
	}

	return nil
}

func (c *Client) DeleteDNSRecord(ctx context.Context, domain, id string) error {
	url := c.baseURL.JoinPath("dns", "delete", domain, id)

	response := status{}
	err := c.do(ctx, url, nil, &response)

	if err != nil {
		return err
	}

	if response.failed() {
		return response
	}

	return nil
}
