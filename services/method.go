package services

import (
	"net/http"

	"github.com/dghubble/sling"
)

// ListMetadata is basic metadata for list queries
type ListMetadata struct {
	TotalCount int       `json:"totalCount"`
	Offset     int       `json:"offset"`
	Count      int       `json:"count"`
	Links      ListLinks `json:"links"`
}

// Method is a payment method type
// https://www.mollie.com/nl/docs/reference/methods/get
type Method struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Image       struct {
		Normal string `json:"normal"`
		Bigger string `json:"bigger"`
	} `json:"image"`
	Amount struct {
		Minimum string `json:"minimum"`
		Maximum string `json:"maximum"`
	} `json:"amount"`
}

// MethodList is a list of method objects and list metadata
// https://www.mollie.com/nl/docs/reference/methods/list#response
type MethodList struct {
	Data         []*Method `json:"data"`
	ListMetadata `bson:",inline"`
}

// MethodService provides methods for accessing payment methods.
type MethodService struct {
	sling *sling.Sling
}

// NewMethodService returns a new MethodService.
func NewMethodService(accessToken string) *MethodService {
	// Create mollie api client
	client := NewClient(accessToken)

	return &MethodService{
		sling: client,
	}
}

// List returns the methods available for payments
func (s *MethodService) List() (MethodList, *http.Response, error) {
	methods := new(MethodList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path("methods").Receive(methods, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}

	return *methods, resp, err
}
