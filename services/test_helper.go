package services

import "os"

// TestClient is a tiny Mollie API test client
type TestClient struct {
	MethodService       *MethodService
	PaymentService      *PaymentService
	CustomerService     *CustomerService
	MandateService      *MandateService
	SubscriptionService *SubscriptionService
	OrderService        *OrderService
}

// NewTestClient returns a new TestClient
func NewTestClient(accessToken string) *TestClient {
	return &TestClient{
		MethodService:       NewMethodService(accessToken),
		PaymentService:      NewPaymentService(accessToken),
		CustomerService:     NewCustomerService(accessToken),
		MandateService:      NewMandateService(accessToken),
		SubscriptionService: NewSubscriptionService(accessToken),
		OrderService:        NewOrderService(accessToken),
	}
}

func testClient() *TestClient {
	return NewTestClient(os.Getenv("MOLLIE_ACCESS_TOKEN"))
}

func testOrderRequest() OrderRequest {
	amt := Amount{
		Currency: "EUR",
		Value:    "10.50",
	}
	l := []*OrderLineRequest{}
	q := uint8(1)
	l = append(l, &OrderLineRequest{
		Name:     "Grape",
		Quantity: &q,
		UnitPrice: &Amount{
			Currency: "EUR",
			Value:    "10.50",
		},
		TotalAmount: &Amount{
			Currency: "EUR",
			Value:    "10.50",
		},
		VatRate: "21.00",
		VatAmount: &Amount{
			Currency: "EUR",
			Value:    "1.82",
		},
	})
	a := Address{
		StreetAndNumber: "1 Riverview",
		PostalCode:      "1234CD",
		City:            "Riverville",
		Country:         "AU",
		GivenName:       "Brian",
		FamilyName:      "Brown",
		Email:           "brian@brown.com",
	}
	ru := "http://www.brown.com/payment"
	hu := "http://www.brown.com/hook"

	return OrderRequest{
		Amount:          amt,
		RedirectURL:     ru,
		WebhookURL:      hu,
		Locale:          "en-GB",
		OrderNumber:     "12345abcde",
		Lines:           l,
		BillingAddress:  a,
		ShippingAddress: &a,
	}
}
