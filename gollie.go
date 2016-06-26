// Mollie API access (partial) using token authentication
package gollie

import "gopkg.in/rollick/gollie.v2/services"

//
// Client to wrap services
//

// Client is a tiny Mollie API client
type Client struct {
	MethodService       *services.MethodService
	PaymentService      *services.PaymentService
	CustomerService     *services.CustomerService
	SubscriptionService *services.SubscriptionService
	// TODO: Other service endpoints to be added
}

// NewClient returns a new Client
func NewClient(accessToken string) *Client {
	return &Client{
		MethodService:       services.NewMethodService(accessToken),
		PaymentService:      services.NewPaymentService(accessToken),
		CustomerService:     services.NewCustomerService(accessToken),
		SubscriptionService: services.NewSubscriptionService(accessToken),
	}
}
