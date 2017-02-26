package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"
	"github.com/rollick/decimal"
)

// SubscriptionService provides methods for accessing subscription records.
type SubscriptionService struct {
	sling *sling.Sling
}

// Subscription is a subscription object
// https://www.mollie.com/nl/docs/reference/subscriptions/get#response
type Subscription struct {
	Resource    string          `json:"resource"`
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	Interval    string          `json:"interval"`
	Times       int             `json:"times"`
	Mode        string          `json:"mode"`
	Method      string          `json:"method"`
	Status      string          `json:"status"`
	Locale      string          `json:"locale"`
	ProfileID   string          `json:"profileId"`
	CustomerID  string          `json:"customerId"`
	CancelledAt *time.Time      `json:"cancelledDatetime"`
	CreatedAt   *time.Time      `json:"createdDatetime"`
	StartDate   string          `json:"startDate"`
	Links       PaymentLinks    `json:"links"`
}

// SubscriptionList is a list of subscription objects and list metadata
// https://www.mollie.com/nl/docs/reference/subscriptions/list#response
type SubscriptionList struct {
	Data         []*Subscription `json:"data"`
	ListMetadata `bson:",inline"`
}

// SubscriptionRequest is a subscription create request
// https://www.mollie.com/nl/docs/reference/subscriptions/create#parameters
type SubscriptionRequest struct {
	Amount      decimal.Decimal `json:"amount,omitempty"`
	Times       int             `json:"times,omitempty"`
	Interval    string          `json:"interval,omitempty"`
	StartDate   string          `json:"startDate,omitempty"`
	Description string          `json:"description,omitempty"`
	Method      string          `json:"method,omitempty"`
	WebhookUrl  string          `json:"webhookUrl,omitempty"`
}

// NewSubscriptionService returns a new SubscriptionService.
func NewSubscriptionService(accessToken string) *SubscriptionService {
	client := NewClient(accessToken)

	return &SubscriptionService{
		sling: client,
	}
}

// List returns all subscriptions created.
func (s *SubscriptionService) List(customerId string, params *ListParams) (SubscriptionList, *http.Response, error) {
	subscriptions := new(SubscriptionList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path(fmt.Sprintf("customers/%s/subscriptions", customerId)).QueryStruct(params).Receive(subscriptions, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}

	return *subscriptions, resp, err
}

// Fetch returns a created subscription
func (s *SubscriptionService) Fetch(customerId string, subscriptionId string) (Subscription, *http.Response, error) {
	subscription := new(Subscription)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Get(fmt.Sprintf("customers/%s/subscriptions/%s", customerId, subscriptionId)).Receive(subscription, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}
	return *subscription, resp, err
}

// Create creates a new subscription
func (s *SubscriptionService) Create(customerId string, subscriptionBody *SubscriptionRequest) (Subscription, *http.Response, error) {
	subscription := new(Subscription)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post(fmt.Sprintf("customers/%s/subscriptions", customerId)).BodyJSON(subscriptionBody).Receive(subscription, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}
	return *subscription, resp, err
}
