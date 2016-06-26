package services

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
	"github.com/shopspring/decimal"
)

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

// PaymentRequest is a payment request
// https://www.mollie.com/nl/docs/reference/payments/create
type PaymentRequest struct {
	Amount        decimal.Decimal `json:"amount,omitempty"`
	Description   string          `json:"description,omitempty"`
	RedirectUrl   string          `json:"redirectUrl,omitempty"`
	WebhookUrl    string          `json:"webhookUrl,omitempty"`
	Method        string          `json:"method,omitempty"`
	Locale        string          `json:"locale,omitempty"`
	RecurringType string          `json:"recurringType,omitempty"`
	Metadata      interface{}     `json:"metadata,omitempty"`
}

// PaymentService provides methods for creating and reading payments.
type PaymentService struct {
	sling *sling.Sling
}

// NewPaymentService returns a new PaymentService.
func NewPaymentService(accessToken string) *PaymentService {
	// Create mollie api client
	client := NewClient(accessToken)

	return &PaymentService{
		sling: client,
	}
}

// List returns the accessible payments
func (s *PaymentService) List(params *ListParams) (PaymentList, *http.Response, error) {
	payments := new(PaymentList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path("payments").QueryStruct(params).Receive(payments, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}

	return *payments, resp, err
}

// Get an existing payment
func (s *PaymentService) Fetch(paymentId string) (Payment, *http.Response, error) {
	payment := new(Payment)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Get(fmt.Sprintf("payments/%s", paymentId)).Receive(payment, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}
	return *payment, resp, err
}

// Creates a new payment
func (s *PaymentService) Create(paymentBody *PaymentRequest) (Payment, *http.Response, error) {
	payment := new(Payment)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post("payments").BodyJSON(paymentBody).Receive(payment, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}
	return *payment, resp, err
}