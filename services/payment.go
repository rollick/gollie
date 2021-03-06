package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"
	"github.com/rollick/decimal"
)

// Payment is a payment object
// https://www.mollie.com/nl/docs/reference/payments/get#response
type Payment struct {
	ID                string          `json:"id"`
	Resource          string          `json:"resource"`
	Description       string          `json:"description"`
	CreatedDatetime   *time.Time      `json:"createdDatetime"`
	PaidDatetime      *time.Time      `json:"paidDatetime"`
	CancelledDatetime *time.Time      `json:"cancelledDatetime"`
	ExpiredDatetime   *time.Time      `json:"expiredDatetime"`
	ExpiryPeriod      string          `json:"expiryPeriod"`
	FailedDatetime    *time.Time      `json:"failedDatetime"`
	Amount            decimal.Decimal `json:"amount"`
	AmountRemaining   decimal.Decimal `json:"amountRemaining"`
	AmountRefunded    decimal.Decimal `json:"amountRefunded"`
	Mode              string          `json:"mode"`
	Method            string          `json:"method"`
	Status            string          `json:"status"`
	Locale            string          `json:"locale"`
	CountryCode       string          `json:"countryCode"`
	ProfileID         string          `json:"profileId"`
	CustomerID        string          `json:"customerId"`
	MandateID         string          `json:"mandateId"`
	SubscriptionID    string          `json:"subscriptionId"`
	SettlementID      string          `json:"settlementId"`
	RecurringType     string          `json:"recurringType"`
	FailureReason     string          `json:"failureReason"`
	ApplicationFee    ApplicationFee  `json:"applicationFee"`
	Issuer            string          `json:"issuer"`
	Metadata          interface{}     `json:"metadata"`
	Details           interface{}     `json:"details"`
	Links             PaymentLinks    `json:"links"`
}

// ApplicationFee is the application fee, if the payment was created with one.
type ApplicationFee struct {
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}

// PaymentLinks respresents the links object returned in a Payment
// https://www.mollie.com/en/docs/reference/payments/get#response
type PaymentLinks struct {
	PaymentUrl  string `json:"paymentUrl"`
	WebhookUrl  string `json:"webhookUrl"`
	RedirectUrl string `json:"redirectUrl"`
	Settlement  string `json:"settlement"`
	Refunds     string `json:"refunds"`
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
	CustomerID    string          `json:"customerId,omitempty"`
	MandateID     string          `json:"mandateId,omitempty"`
	Metadata      interface{}     `json:"metadata,omitempty"`
}

// PaymentRefund is a payment refund response
// https://www.mollie.com/en/docs/reference/refunds/get#response
type PaymentRefund struct {
	ID             string          `json:"id"`
	Payment        Payment         `json:"payment"`
	Amount         decimal.Decimal `json:"amount"`
	Status         string          `json:"status"`
	RefundDatetime *time.Time      `json:"refundDatetime"`
}

// PaymentRefundRequest is a payment refund request
// https://www.mollie.com/en/docs/reference/refunds/create
type PaymentRefundRequest struct {
	Amount      decimal.Decimal `json:"amount,omitempty"`
	Description string          `json:"description,omitempty"`
}

// PaymentRefundList is a list of payment refund objects and list metadata
// https://www.mollie.com/en/docs/reference/refunds/list#response
type PaymentRefundList struct {
	Data         []*PaymentRefund `json:"data"`
	ListMetadata `bson:",inline"`
}

// PaymentChargeback is a payment chargeback response
// https://www.mollie.com/en/docs/reference/chargebacks/get#response
type PaymentChargeback struct {
	ID                 string          `json:"id"`
	Payment            Payment         `json:"payment"`
	Amount             decimal.Decimal `json:"amount"`
	Status             string          `json:"status"`
	ChargebackDatetime *time.Time      `json:"chargebackDatetime"`
	ReversedDatetime   *time.Time      `json:"reversedDatetime"`
}

// PaymentChargebackList is a list of payment chargeback objects and list metadata
// https://www.mollie.com/en/docs/reference/chargebacks/list#response
type PaymentChargebackList struct {
	Data         []*PaymentChargeback `json:"data"`
	ListMetadata `bson:",inline"`
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
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}

	return *payments, resp, err
}

// Fetch returns an existing payment
func (s *PaymentService) Fetch(paymentId string) (Payment, *http.Response, error) {
	payment := new(Payment)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Get(fmt.Sprintf("payments/%s", paymentId)).Receive(payment, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}
	return *payment, resp, err
}

// Create creates a new payment
func (s *PaymentService) Create(paymentBody *PaymentRequest) (Payment, *http.Response, error) {
	payment := new(Payment)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post("payments").BodyJSON(paymentBody).Receive(payment, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}
	return *payment, resp, err
}

// CreateRefund creates a new payment refund
func (s *PaymentService) CreateRefund(paymentId string, refundBody *PaymentRefundRequest) (PaymentRefund, *http.Response, error) {
	refund := new(PaymentRefund)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post(fmt.Sprintf("payments/%s/refunds", paymentId)).BodyJSON(refundBody).Receive(refund, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}
	return *refund, resp, err
}

// FetchRefund returns a payment refund
func (s *PaymentService) FetchRefund(paymentId string, refundId string) (PaymentRefund, *http.Response, error) {
	refund := new(PaymentRefund)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Get(fmt.Sprintf("payments/%s/refunds/%s", paymentId, refundId)).Receive(refund, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}
	return *refund, resp, err
}

// RefundList returns all payment refunds created
func (s *PaymentService) RefundList(paymentId string, params *ListParams) (PaymentRefundList, *http.Response, error) {
	refunds := new(PaymentRefundList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path(fmt.Sprintf("payments/%s/refunds", paymentId)).QueryStruct(params).Receive(refunds, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}

	return *refunds, resp, err
}

// FetchChargeback returns a payment chargeback
func (s *PaymentService) FetchChargeback(paymentId string, chargebackId string) (PaymentChargeback, *http.Response, error) {
	chargeback := new(PaymentChargeback)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Get(fmt.Sprintf("payments/%s/chargebacks/%s", paymentId, chargebackId)).Receive(chargeback, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}
	return *chargeback, resp, err
}

// ChargebackList returns all payment chargebacks created
func (s *PaymentService) ChargebackList(paymentId string, params *ListParams) (PaymentChargebackList, *http.Response, error) {
	chargebacks := new(PaymentChargebackList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path(fmt.Sprintf("payments/%s/chargebacks", paymentId)).QueryStruct(params).Receive(chargebacks, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}

	return *chargebacks, resp, err
}
