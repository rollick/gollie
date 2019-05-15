package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

// Chargeback is a payment chargeback response
// https://docs.mollie.com/reference/v2/chargebacks-api/get-chargeback#parameters
type Chargeback struct {
	ID               string          `json:"id"`
	Amount           Amount          `json:"amount"`
	SettlementAmount *Amount         `json:"settlementAmount,omitempty"`
	CreatedAt        time.Time       `json:"createdAt"`
	ReversedAt       time.Time       `json:"reversedAt"`
	Links            ChargebackLinks `json:"_links"`
}

// ChargebackLinks is payment chargeback links
// https://docs.mollie.com/reference/v2/chargebacks-api/get-chargeback#parameters
type ChargebackLinks struct {
	Links map[string]struct {
		Self          Link  `json:"self"`
		Payment       Link  `json:"payment"`
		Settlement    *Link `json:"settlement,omitempty"`
		Documentation Link  `json:"documentation"`
	} `json:"_links"`
}

// ChargebackList is a list of payment chargeback objects and list metadata
// https://www.mollie.com/en/docs/reference/chargebacks/list#response
type ChargebackList struct {
	Data         []*Chargeback `json:"data"`
	ListMetadata `bson:",inline"`
}

// ChargebackService provides methods for creating and reading payment chargebacks
type ChargebackService struct {
	sling *sling.Sling
}

// NewChargebackService returns a new ChargebackService
func NewChargebackService(accessToken string) *ChargebackService {
	// Create mollie api client
	client := NewClient(accessToken)

	return &ChargebackService{
		sling: client,
	}
}

// FetchChargeback returns a payment chargeback
func (s *ChargebackService) Fetch(paymentId string, chargebackId string) (Chargeback, *http.Response, error) {
	chargeback := new(Chargeback)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Get(fmt.Sprintf("payments/%s/chargebacks/%s", paymentId, chargebackId)).Receive(chargeback, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *chargeback, resp, err
}

// ChargebackList returns all payment chargebacks created
func (s *ChargebackService) List(paymentId string, params *ListParams) (ChargebackList, *http.Response, error) {
	chargebacks := new(ChargebackList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path(fmt.Sprintf("payments/%s/chargebacks", paymentId)).QueryStruct(params).Receive(chargebacks, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}

	return *chargebacks, resp, err
}
