package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

type applicationFee struct {
	Amount      Amount `json:"amount"`
	Description string `json:"description"`
}

// Payment is a payment object
// https://docs.mollie.com/reference/v2/payments-api/get-payment#response
type Payment struct {
	Resource         string          `json:"resource"`
	ID               string          `json:"id"`
	Mode             string          `json:"mode"`
	CreatedAt        time.Time       `json:"createdAt"`
	Status           string          `json:"status"`
	IsCancelable     bool            `json:"isCancelable"`
	AuthorizedAt     *time.Time      `json:"authorizedAt,omitempty"`
	PaidAt           *time.Time      `json:"paidAt,omitempty"`
	CanceledAt       *time.Time      `json:"canceledAt,omitempty"`
	ExpiresAt        time.Time       `json:"expiresAt"`
	ExpiredAt        *time.Time      `json:"expiredAt,omitempty"`
	FailedAt         time.Time       `json:"failedAt"`
	Amount           Amount          `json:"amount"`
	AmountRefunded   *Amount         `json:"amountRefunded,omitempty"`
	AmountRemaining  *Amount         `json:"amountRemaining,omitempty"`
	AmountCaptured   *Amount         `json:"amountCaptured,omitempty"`
	Description      string          `json:"description"`
	RedirectUrl      *string         `json:"redirectUrl,omitempty"`
	WebhookUrl       *string         `json:"webhookUrl,omitempty"`
	Method           string          `json:"method"`
	Details          interface{}     `json:"details"`
	Metadata         interface{}     `json:"metadata"`
	Locale           string          `json:"locale"`
	CountryCode      *string         `json:"countryCode,omitempty"`
	ProfileID        string          `json:"profileId"`
	SettlementAmount *Amount         `json:"settlementAmount,omitempty"`
	SettlementID     *string         `json:"settlementId,omitempty"`
	CustomerID       *string         `json:"customerId,omitempty"`
	SequenceType     string          `json:"sequenceType"`
	MandateID        *string         `json:"mandateId,omitempty"`
	SubscriptionID   *string         `json:"subscriptionId,omitempty"`
	OrderID          *string         `json:"orderId,omitempty"`
	ApplicationFee   *applicationFee `json:"applicationFee,omitempty"`
	Links            PaymentLinks    `json:",inline"`
}

// PaymentLinks respresents the links object returned in a Payment
// https://docs.mollie.com/reference/v2/payments-api/get-payment#response
type PaymentLinks struct {
	Links map[string]struct {
		Self               Link  `json:"self"`
		Checkout           *Link `json:"checkout,omitempty"`
		ChangePaymentState Link  `json:"changePaymentState"`
		Refunds            *Link `json:"refunds,omitempty"`
		Chargebacks        *Link `json:"chargebacks,omitempty"`
		Captures           *Link `json:"captures,omitempty"`
		Settlement         *Link `json:"settlement,omitempty"`
		Documentation      Link  `json:"documentation"`
		Mandate            *Link `json:"mandate,omitempty"`
		Customer           *Link `json:"customer,omitempty"`
		Order              *Link `json:"order,omitempty"`
		Status             *Link `json:"status,omitempty"`
		PayOnline          *Link `json:"payOnline,omitempty"`
	} `json:"_links"`
}

// PaymentList is a list of payment objects and list metadata
// https://docs.mollie.com/reference/v2/payments-api/list-payments
type PaymentList struct {
	Data         []*Payment `json:"data"`
	ListMetadata `json:",inline"`
}

// PaymentRequest is a payment request
// https://docs.mollie.com/reference/v2/payments-api/create-payment#parameters
type PaymentRequest struct {
	Amount       Amount      `json:"amount"`
	Description  string      `json:"description"`
	RedirectUrl  string      `json:"redirectUrl"`
	WebhookUrl   string      `json:"webhookUrl,omitempty"`
	Locale       string      `json:"locale,omitempty"`
	Method       string      `json:"method,omitempty"`
	Metadata     interface{} `json:"metadata,omitempty"`
	SequenceType string      `json:"sequenceType,omitempty"`
	CustomerID   string      `json:"customerId,omitempty"`
	MandateID    string      `json:"mandateId,omitempty"`
}

// PaymentService provides methods for creating and reading payments
type PaymentService struct {
	sling *sling.Sling
}

// NewPaymentService returns a new PaymentService
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
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}

	return *payments, resp, err
}

// Fetch returns an existing payment
func (s *PaymentService) Fetch(paymentId string) (Payment, *http.Response, error) {
	payment := new(Payment)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Get(fmt.Sprintf("payments/%s", paymentId)).Receive(payment, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *payment, resp, err
}

// Create creates a new payment
func (s *PaymentService) Create(paymentBody *PaymentRequest) (Payment, *http.Response, error) {
	payment := new(Payment)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post("payments").BodyJSON(paymentBody).Receive(payment, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *payment, resp, err
}

// Cancel will cancel a payment if possible
func (s *PaymentService) Cancel(paymentId string) (Payment, *http.Response, error) {
	payment := new(Payment)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Delete(fmt.Sprintf("payments/%s", paymentId)).Receive(payment, mollieError)
	log.Printf("+%v", mollieError.Links)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *payment, resp, err
}
