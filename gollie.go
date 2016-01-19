// Mollie API access (partial) using token authentication

package gollie

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/shopspring/decimal"
)

const (
	baseURL    = "https://api.mollie.nl"
	apiVersion = "v1"
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
	ID          int    `json:"id"`
	Description string `json:"description"`
	Image       struct {
		Normal string `json:"normal"`
		Bigger string `json:"bigger"`
	} `json:"image"`
	Amount struct {
		Minimum int `json:"minimum"`
		Maximum int `json:"maximum"`
	} `json:"amount"`
}

// PaymentList is a list of payment objects and list metadata
// https://www.mollie.com/nl/docs/reference/payments/list#response
type MethodList struct {
	Data         []*Payment `json:"data"`
	ListMetadata `bson:",inline"`
}

// Payment is a payment object
// https://www.mollie.com/nl/docs/reference/payments/get#response
type Payment struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	Mode        string          `json:"mode"`
	Method      string          `json:"method"`
	Status      string          `json:"status"`
	Locale      string          `json:"locale"`
	ProfileID   string          `json:"profileId"`
	Metadata    interface{}     `json:"metadata"`
	Links       PaymentLinks    `json:"links"`
}

type PaymentLinks struct {
	PaymentUrl  string `json:"paymentUrl"`
	WebhookUrl  string `json:"webhookUrl"`
	RedirectUrl string `json:"redirectUrl"`
	Settlement  string `json:"settlement"`
}

// PaymentList is a list of payment objects and list metadata
// https://www.mollie.com/nl/docs/reference/payments/list#response
type PaymentList struct {
	Data         []*Payment `json:"data"`
	ListMetadata `bson:",inline"`
}

// MollieError represents a Mollie API error response
type MollieError struct {
	Err struct {
		Message string `json:"message"`
		Field   string `json:"field"`
	} `json:"error"`
}

// Errorr is a Formatted MollieError
func (e MollieError) Error() string {
	return fmt.Sprintf("Mollie Error: %v %v", e.Err.Message, e.Err.Field)
}

// PaymentRequest is a payment request
// https://www.mollie.com/nl/docs/reference/payments/create
type PaymentRequest struct {
	Amount      decimal.Decimal `json:"amount,omitempty"`
	Description string          `json:"description,omitempty"`
	RedirectUrl string          `json:"redirectUrl,omitempty"`
	WebhookUrl  string          `json:"webhookUrl,omitempty"`
	Method      string          `json:"method,omitempty"`
	Locale      string          `json:"locale,omitempty"`
	Metadata    interface{}     `json:"metadata,omitempty"`
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

//
// Services
//

// MethodService provides methods for creating and reading issues.
type MethodService struct {
	sling *sling.Sling
}

// NewMethodService returns a new MethodService.
func NewMethodService(accessToken string) *MethodService {
	// Create mollie api client
	client := sling.New().Client(nil).Base(fmt.Sprintf("%s/%s/", baseURL, apiVersion))

	// Add request headers
	client.Set("authorization", fmt.Sprintf("Bearer %s", accessToken))
	client.Set("user-agent", "Mollie/1.1.8 Go/1.4 OpenSSL/1.0.2d")

	return &MethodService{
		sling: client,
	}
}

// List returns the methods available for payments
func (s *MethodService) List() (MethodList, *http.Response, error) {
	methods := new(MethodList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path("methods").Receive(methods, mollieError)

	if err == nil {
		err = mollieError
	}

	return *methods, resp, err
}

// PaymentService provides methods for creating and reading issues.
type PaymentService struct {
	sling *sling.Sling
}

// NewPaymentService returns a new PaymentService.
func NewPaymentService(accessToken string) *PaymentService {
	// Create mollie api client
	client := sling.New().Client(nil).Base(fmt.Sprintf("%s/%s/", baseURL, apiVersion))

	// Add request headers
	client.Set("authorization", fmt.Sprintf("Bearer %s", accessToken))
	client.Set("user_agent", "Mollie/1.1.8 Go/1.4 OpenSSL/1.0.2d")

	return &PaymentService{
		sling: client,
	}
}

// List returns the accessible payments
func (s *PaymentService) List(params *ListParams) (PaymentList, *http.Response, error) {
	payments := new(PaymentList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path("payments").QueryStruct(params).Receive(payments, mollieError)
	if err == nil {
		err = mollieError
	}

	return *payments, resp, err
}

// Get an existing payment
func (s *PaymentService) Fetch(paymentId string) (Payment, *http.Response, error) {
	payment := new(Payment)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Get(fmt.Sprintf("payments/%s", paymentId)).Receive(payment, mollieError)

	fmt.Println(resp.Status)
	if err == nil {
		err = mollieError
	}
	return *payment, resp, err
}

// Creates a new payment
func (s *PaymentService) Create(paymentBody *PaymentRequest) (Payment, *http.Response, error) {
	payment := new(Payment)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post("payments").BodyJSON(paymentBody).Receive(payment, mollieError)

	fmt.Println(resp.Status)
	if err == nil {
		err = mollieError
	}
	return *payment, resp, err
}

//
// Client to wrap services
//

// Client is a tiny Mollie API client
type Client struct {
	MethodService  *MethodService
	PaymentService *PaymentService
	// TODO: Other service endpoints to be added
}

// NewClient returns a new Client
func NewClient(accessToken string) *Client {
	return &Client{
		MethodService:  NewMethodService(accessToken),
		PaymentService: NewPaymentService(accessToken),
	}
}
