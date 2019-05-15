package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

// Mandate is a customer mandate object
// https://www.mollie.com/en/docs/reference/mandates/create#response
type Mandate struct {
	Resource         string `json:"resource"`
	Id               string `json:"id"`
	Status           string `json:"status"`
	Method           string `json:"method"`
	CustomerId       string `json:"customerId"`
	Details          MandateDetails
	MandateReference string     `json:"mandateReference"`
	SignatureDate    string     `json:"signatureDate"`
	CreatedDateTime  *time.Time `json:"createdDateTime"`
}

// MandateDetails is the payment method details for a customer mandate
// https://www.mollie.com/en/docs/reference/mandates/get#response
type MandateDetails struct {
	ConsumerName    string `json:"consumerName"`
	ConsumerAccount string `json:"consumerAccount"`
	ConsumerBic     string `json:"consumerBic"`
	CardHolder      string `json:"cardHolder"`
	CardNumber      string `json:"cardNumber"`
	CardLabel       string `json:"cardLabel"`
	CardFingerprint string `json:"cardFingerprint"`
	CardExpiryDate  string `json:"cardExpiryDate"`
}

// MandateService provides methods for accessing customer mandate records.
type MandateService struct {
	sling *sling.Sling
}

// NewMandateService returns a new MandateService.
func NewMandateService(accessToken string) *MandateService {
	// Create mollie api client
	client := NewClient(accessToken)

	return &MandateService{
		sling: client,
	}
}

// MandateList is a list of customer mandate objects and list metadata
// https://www.mollie.com/en/docs/reference/mandates/list#response
type MandateList struct {
	Data         []*Mandate `json:"data"`
	ListMetadata `bson:",inline"`
}

// MandateList returns a list of mandates for a customer
func (s *MandateService) List(customerId string, params *ListParams) (MandateList, *http.Response, error) {
	mandates := new(MandateList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path(fmt.Sprintf("customers/%s/mandates", customerId)).QueryStruct(params).Receive(mandates, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}

	return *mandates, resp, err
}

// Mandate creates a new customer mandate
func (s *MandateService) Create(customerId string, mandateBody PaymentRequest) (Mandate, *http.Response, error) {
	mandate := new(Mandate)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post(fmt.Sprintf("customers/%s/mandates", customerId)).BodyJSON(mandateBody).Receive(mandate, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}

	return *mandate, resp, err
}

// MandateFetch returns a customer mandate
func (s *MandateService) Fetch(customerId string, mandateId string) (Mandate, *http.Response, error) {
	mandate := new(Mandate)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path(fmt.Sprintf("customers/%s/mandates/%s", customerId, mandateId)).Receive(mandate, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}

	return *mandate, resp, err
}
