package services

import (
	"fmt"

	"github.com/dghubble/sling"
)

const (
	baseURL    = "https://api.mollie.nl"
	apiVersion = "v1"
)

// MollieError represents a Mollie API error response
type MollieError struct {
	Err struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Field   string `json:"field"`
	} `json:"error"`
}

// ListParams are the params for any list request
// https://www.mollie.com/nl/docs/reference/payments/list#parameters
type ListParams struct {
	Offset int `url:"offset,omitempty"`
	Count  int `url:"count,omitempty"`
}

// ListLinks is a standard list links object for a resource list query
type ListLinks struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
	First    string `json:"first"`
	Last     string `json:"last"`
}

// ListMetadata is basic metadata for list queries
type ListMetadata struct {
	TotalCount int       `json:"totalCount"`
	Offset     int       `json:"offset"`
	Count      int       `json:"count"`
	Links      ListLinks `json:"links"`
}

// NewClient returns a new Mollie client
func NewClient(accessToken string) *sling.Sling {
	// Create mollie api client
	client := sling.New().Client(nil).Base(fmt.Sprintf("%s/%s/", baseURL, apiVersion))

	// Add request headers
	client.Set("authorization", fmt.Sprintf("Bearer %s", accessToken))
	client.Set("user-agent", "Mollie/1.1.8 Go/1.4 OpenSSL/1.0.2d")

	return client
}

// Error is a formatted Mollie error
func (e MollieError) Error() string {
	return fmt.Sprintf("Mollie %v error: %v %v", e.Err.Type, e.Err.Message, e.Err.Field)
}
