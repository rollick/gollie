package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dghubble/sling"
	"github.com/shopspring/decimal"
)

// CustomerList is a list of customer objects and list metadata
// https://www.mollie.com/nl/docs/reference/customers/list#response
type CustomerList struct {
	Data         []*Customer `json:"data"`
	ListMetadata `bson:",inline"`
}

// Customer is a customer object
// https://www.mollie.com/nl/docs/reference/customers/get#response
type Customer struct {
	Resource  string    `json:"resource"`
	ID        string    `json:"id"`
	Mode      string    `json:"mode"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Locale    string    `json:"locale"`
	Metadata  string    `json:"metadata"`
	Methods   []string  `json:"recentlyUsedMethods"`
	CreatedAt time.Time `json:"createdDatetime"`
}

// CustomerPayment is a customer payment object
// https://www.mollie.com/nl/docs/reference/customers/get#response
type CustomerPayment struct {
	Resource    string          `json:"respurce"`
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	Mode        string          `json:"mode"`
	Method      string          `json:"method"`
	Status      string          `json:"status"`
	Locale      string          `json:"locale"`
	ProfileID   string          `json:"profileId"`
	CustomerID  string          `json:"customerId"`
	Metadata    interface{}     `json:"metadata"`
	Links       PaymentLinks    `json:"links"`
}

// CustomerRequest is a customer create request
// https://www.mollie.com/nl/docs/reference/customers/create#parameters
type CustomerRequest struct {
	Name     string      `json:"name,omitempty"`
	Email    string      `json:"email,omitempty"`
	Locale   string      `json:"locale,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}

// CustomerService provides methods for accessing customer records.
type CustomerService struct {
	sling *sling.Sling
}

// NewCustomerService returns a new CustomerService.
func NewCustomerService(accessToken string) *CustomerService {
	// Create mollie api client
	client := NewClient(accessToken)

	return &CustomerService{
		sling: client,
	}
}

// List returns all customers created.
func (s *CustomerService) List(params *ListParams) (CustomerList, *http.Response, error) {
	customers := new(CustomerList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path("customers").QueryStruct(params).Receive(customers, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}

	return *customers, resp, err
}

// Fetch returns a created customer
func (s *CustomerService) Fetch(customerId string) (Customer, *http.Response, error) {
	customer := new(Customer)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Get(fmt.Sprintf("customers/%s", customerId)).Receive(customer, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}
	return *customer, resp, err
}

// Create creates a new payment
func (s *CustomerService) Create(customerBody *CustomerRequest) (Customer, *http.Response, error) {
	customer := new(Customer)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post("customers").BodyJSON(customerBody).Receive(customer, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}
	return *customer, resp, err
}

// PaymentList returns all customer payments created
func (s *CustomerService) PaymentList(customerId string, params *ListParams) (PaymentList, *http.Response, error) {
	payments := new(PaymentList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path(fmt.Sprintf("customers/%s/payments", customerId)).QueryStruct(params).Receive(payments, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}

	return *payments, resp, err
}

// Payment creates a new customer payment
func (s *CustomerService) Payment(customerId string, paymentBody PaymentRequest) (CustomerPayment, *http.Response, error) {
	payment := new(CustomerPayment)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post(fmt.Sprintf("customers/%s/payments", customerId)).BodyJSON(paymentBody).Receive(payment, mollieError)
	if err == nil && mollieError.Err.Type != "" {
		err = mollieError
	}

	return *payment, resp, err
}
