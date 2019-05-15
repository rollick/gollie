package services

import (
	"fmt"

	"github.com/dghubble/sling"
)

const (
	baseURL    = "https://api.mollie.nl"
	apiVersion = "v2"
)

// Amount represents payment amount with a currency
// https://docs.mollie.com/guides/common-data-types#amount-object
type Amount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

// Address is an address object
// https://docs.mollie.com/guides/common-data-types#address-object
type Address struct {
	OrganizationName string `json:"organizationName"`
	StreetAndNumber  string `json:"streetAndNumber"`
	PostalCode       string `json:"postalCode"`
	City             string `json:"city"`
	Region           string `json:"region"`
	Country          string `json:"country"`
	Title            string `json:"title"`
	GivenName        string `json:"givenName"`
	FamilyName       string `json:"familyName"`
	Email            string `json:"email"`
	Phone            string `json:"phone,omitemty"`
}

// Link is a link object usually found in the _links response
type Link struct {
	Href string `json:"href"`
	Type string `json:"type"`
}

// Links is a generic links object
type Links struct {
	Links map[string]Link `json:"_links"`
}

// ListParams are the params for any list request
// https://docs.mollie.com/guides/pagination#pagination-in-v2-api-endpoints
type ListParams struct {
	From  *string `url:"from,omitempty"`
	Limit *uint16 `url:"limit,omitempty"`
}

// ListLinks is a standard list links object for a resource list query
// https://docs.mollie.com/guides/pagination#pagination-in-v2-api-endpoints
type ListLinks struct {
	Self          Link `json:"self"`
	Previous      Link `json:"previous"`
	Next          Link `json:"next"`
	Documentation Link `json:"documentation"`
}

// ListMetadata is basic metadata for list queries
// https://docs.mollie.com/guides/pagination#pagination-in-v2-api-endpoints
type ListMetadata struct {
	Count int       `json:"count"`
	Links ListLinks `json:"_links"`
}

// MollieError represents a Mollie API error response
type MollieError struct {
	Status int       `json:"status"`
	Title  string    `json:"title"`
	Detail string    `json:"detail"`
	Links  ListLinks `json:"_links"`
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
	return fmt.Sprintf("Mollie %v error: %v %v", e.Status, e.Title, e.Detail)
}
