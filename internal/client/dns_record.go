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

func (dns DNSRecord) Merge(data DNSRecord) DNSRecord {
	if data.Name != "" {
		dns.Name = data.Name
	}
	if data.Type != "" {
		dns.Type = data.Type
	}
	if data.Content != "" {
		dns.Content = data.Content
	}
	if data.TTL != "" {
		dns.TTL = data.TTL
	}
	if data.Priority != "" {
		dns.Priority = data.Priority
	}
	if data.Notes != "" {
		dns.Notes = data.Notes
	}
	return dns
}

type CreateDNSRecordResponse struct {
	Status
	ID string `json:"id"`
}

func (c *Client) CreateDNSRecord(ctx context.Context, domain string, record DNSRecord) (string, error) {
	url := c.baseURL.JoinPath("dns", "create", domain)

	var response CreateDNSRecordResponse
	err := c.Do(ctx, url, record, &response)

	if err != nil {
		return "", err
	}

	if response.HasFailed() {
		return "", err
	}

	return response.ID, nil
}

type RetrieveDNSRecordResponse struct {
	Status
	Records []DNSRecord `json:"records"`
}

func (c *Client) RetrieveDNSRecord(ctx context.Context, domain, id string) (DNSRecord, error) {
	url := c.baseURL.JoinPath("dns", "retrieve", domain, id)

	var response RetrieveDNSRecordResponse
	err := c.Do(ctx, url, nil, &response)

	if err != nil {
		return DNSRecord{}, err
	}

	if response.HasFailed() {
		return DNSRecord{}, response.Status
	}

	if len(response.Records) < 1 {
		return DNSRecord{}, fmt.Errorf("DNS record not found")
	}

	return response.Records[0], nil
}

func (c *Client) EditDNSRecord(ctx context.Context, domain, id string, record DNSRecord) error {
	url := c.baseURL.JoinPath("dns", "edit", domain, id)

	response := Status{}
	err := c.Do(ctx, url, record, &response)

	if err != nil {
		return err
	}

	if response.HasFailed() {
		return response
	}

	return nil
}

func (c *Client) DeleteDNSRecord(ctx context.Context, domain, id string) error {
	url := c.baseURL.JoinPath("dns", "delete", domain, id)

	response := Status{}
	err := c.Do(ctx, url, nil, &response)

	if err != nil {
		return err
	}

	if response.HasFailed() {
		return response
	}

	return nil
}
