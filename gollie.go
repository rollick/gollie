package gollie

import "github.com/rollick/gollie/services"

//
// Client to wrap services
//

// Client is a tiny Mollie API client
type Client struct {
	MethodService       *services.MethodService
	PaymentService      *services.PaymentService
	CustomerService     *services.CustomerService
	MandateService      *services.MandateService
	SubscriptionService *services.SubscriptionService
	OrderService        *services.OrderService
}

// NewClient returns a new Client
func NewClient(accessToken string) *Client {
	return &Client{
		MethodService:       services.NewMethodService(accessToken),
		PaymentService:      services.NewPaymentService(accessToken),
		CustomerService:     services.NewCustomerService(accessToken),
		MandateService:      services.NewMandateService(accessToken),
		SubscriptionService: services.NewSubscriptionService(accessToken),
		OrderService:        services.NewOrderService(accessToken),
	}
}
