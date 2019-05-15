package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dghubble/sling"
	"github.com/rollick/decimal"
)

// OrderLine is an order line
// https://docs.mollie.com/reference/v2/orders-api/create-order#order-line-details
type OrderLine struct {
	Resource       *string `json:"resource,omitempty"`
	ID             *string `json:"id,omitempty"`
	OrderID        *string `json:"orderId,omitempty"`
	Type           string  `json:"type,omitempty"`
	Name           string  `json:"name"`
	Quantity       uint8   `json:"quantity"`
	UnitPrice      Amount  `json:"unitPrice"`
	DiscountAmount *Amount `json:"discountAmount,omitempty"`
	TotalAmount    Amount  `json:"totalAmount"`
	VatRate        string  `json:"vatRate"`
	VatAmount      Amount  `json:"vatAmount"`
	SKU            string  `json:"sku,omitempty"`
	ImageURL       string  `json:"imageUrl,omitempty"`
	ProductURL     string  `json:"productUrl,omitempty"`
}

// Order is a order object
// https://www.mollie.com/nl/docs/reference/orders/get#response
type Order struct {
	Resource            string       `json:"resource"`
	ID                  string       `json:"id"`
	ProfileID           string       `json:"profileId"`
	Method              string       `json:"method"`
	Mode                string       `json:"mode"`
	Amount              Amount       `json:"amount"`
	AmountCaptured      *Amount      `json:"amountCaptured,omitempty"`
	AmountRefunded      *Amount      `json:"amountRefunded,omitempty"`
	Status              string       `json:"status"`
	IsCancelable        bool         `json:"isCancelable"`
	BillingAddress      *Address     `json:"billingAddress"`
	ConsumerDateOfBirth string       `json:"consumerDateOfBirth,omitempty"`
	OrderNumber         string       `json:"orderNumber"`
	ShippingAddress     *Address     `json:"shippingAddress,omitempty"`
	Lines               []*OrderLine `json:"lines"`
	Locale              string       `json:"locale"`
	Metadata            interface{}  `json:"metadata"`
	RedirectURL         *string      `json:"redirectUrl"`
	WebhookURL          *string      `json:"webhookUrl,omitempty"`
	CreatedAt           time.Time    `json:"createdAt"`
	ExpiresAt           time.Time    `json:"expiresAt"`
	ExpiredAt           *time.Time   `json:"expiredAt,omitempty"`
	PaidAt              *time.Time   `json:"paidAt,omitempty"`
	AuthorizedAt        *time.Time   `json:"authorizedAt,omitempty"`
	CanceledAt          *time.Time   `json:"canceledAt,omitempty"`
	CompletedAt         *time.Time   `json:"completedAt,omitempty"`
	Links               OrderLinks   `json:"links"`
	Embedded            struct {
		Payments []*Payment `json:"payments"`
		Refunds  []*Refund  `json:"refunds"`
	} `json:"_embedded"`
}

// OrderLinks respresents the links object returned in a Order
// https://docs.mollie.com/reference/v2/orders-api/get-order#
type OrderLinks struct {
	Links map[string]struct {
		Self          Link  `json:"self"`
		Checkout      *Link `json:"checkout,omitempty"`
		Documentation Link  `json:"documentation"`
	}
}

// OrderList is a list of order objects and list metadata
// https://docs.mollie.com/reference/v2/orders-api/list-orders#parameters
type OrderList struct {
	Embedded struct {
		Orders []*Order `json:"orders"`
	} `json:"_embedded"`
	ListMetadata `bson:",inline"`
}

// OrderRequestPayment represents optional order payment details
// https://docs.mollie.com/reference/v2/orders-api/create-order#payment-parameters
type OrderRequestPayment struct {
	ConsumerAccount   string `json:"consumerAccount,omitempty"`
	CustomerID        string `json:"customerId,omitempty"`
	CustomerReference string `json:"customerReference,omitempty"`
	Issuer            string `json:"issuer,omitempty"`
	MandateID         string `json:"mandateId,omitempty"`
	SequenceType      string `json:"sequenceType,omitempty"`
	VoucherNumber     string `json:"voucherNumber,omitempty"`
	VoucherPin        string `json:"voucherPin,omitempty"`
	WebhookURL        string `json:"webhookUrl,omitempty"`
}

// OrderRequest is a order request
// https://docs.mollie.com/reference/v2/orders-api/create-order#parameters
type OrderRequest struct {
	Amount              Amount              `json:"amount"`
	OrderNumber         string              `json:"orderNumber"`
	Lines               []*OrderLineRequest `json:"lines"`
	BillingAddress      Address             `json:"billingAddress"`
	ShippingAddress     *Address            `json:"shippingAddress,omitempty"`
	ConsumerDateOfBirth string              `json:"consumerDateOfBirth,omitempty"`
	RedirectURL         string              `json:"redirectUrl,omitempty"`
	WebhookURL          string              `json:"webhookUrl,omitempty"`
	Locale              string              `json:"locale"`
	Method              interface{}         `json:"method,omitempty"`
	Payment             OrderRequestPayment `json:"payment,omitempty"`
	Metadata            interface{}         `json:"metadata,omitempty"`
}

// OrderUpdateRequest is a order update request
// https://docs.mollie.com/reference/v2/orders-api/update-order#parameters
type OrderUpdateRequest struct {
	OrderNumber     string   `json:"orderNumber,omitempty"`
	BillingAddress  *Address `json:"billingAddress,omitempty"`
	ShippingAddress *Address `json:"shippingAddress,omitempty"`
}

// OrderRefund is an order refund response
// https://www.mollie.com/en/docs/reference/refunds/get#response
type OrderRefund struct {
	ID             string          `json:"id"`
	Order          Order           `json:"order"`
	Amount         decimal.Decimal `json:"amount"`
	Status         string          `json:"status"`
	RefundDatetime *time.Time      `json:"refundDatetime"`
}

// OrderRefundLine represents a line of an order to be refunded
// https://docs.mollie.com/reference/v2/orders-api/create-order-refund#parameters
type OrderRefundLine struct {
	ID       string  `json:"id"`
	Quantity *uint8  `json:"quantity,omitempty"`
	Amount   *Amount `json:"amount,omitempty"`
}

// OrderRefundRequest is an order refund request
// https://www.mollie.com/en/docs/reference/refunds/create
type OrderRefundRequest struct {
	Lines       []OrderRefundLine `json:"lines"`
	Description string            `json:"description,omitempty"`
}

// OrderLineRequest is an order line request
// https://docs.mollie.com/reference/v2/orders-api/update-orderline#parameters
type OrderLineRequest struct {
	Type           string  `json:"type,omitempty"`
	Name           string  `json:"name,omitempty"`
	Quantity       *uint8  `json:"quantity,omitempty"`
	UnitPrice      *Amount `json:"unitPrice,omitempty"`
	DiscountAmount *Amount `json:"discountAmount,omitempty"`
	TotalAmount    *Amount `json:"totalAmount,omitempty"`
	VatRate        string  `json:"vatRate,omitempty"`
	VatAmount      *Amount `json:"vatAmount,omitempty"`
	SKU            string  `json:"sku,omitempty,omitempty"`
	ImageURL       string  `json:"imageUrl,omitempty"`
	ProductURL     string  `json:"productUrl,omitempty"`
}

// OrderPaymentRequest is an order payment request
// https://docs.mollie.com/reference/v2/orders-api/create-order-payment#parameters
type OrderPaymentRequest struct {
}

// OrderRefundList is a list of order refund objects and list metadata
// https://www.mollie.com/en/docs/reference/refunds/list#response
type OrderRefundList struct {
	Embedded struct {
		OrderRefunds []*OrderRefund `json:"refunds"`
	} `json:"_embedded"`
	ListMetadata `bson:",inline"`
}

// OrderService provides methods for creating and reading orders
type OrderService struct {
	sling *sling.Sling
}

// NewOrderService returns a new OrderService
func NewOrderService(accessToken string) *OrderService {
	// Create mollie api client
	client := NewClient(accessToken)

	return &OrderService{
		sling: client,
	}
}

// List returns the accessible orders
func (s *OrderService) List(params *ListParams) (OrderList, *http.Response, error) {
	orders := new(OrderList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path("orders").QueryStruct(params).Receive(orders, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}

	return *orders, resp, err
}

// Fetch returns an existing order
func (s *OrderService) Fetch(orderId string, embed *string) (Order, *http.Response, error) {
	order := new(Order)
	mollieError := new(MollieError)
	path := fmt.Sprintf("orders/%s", orderId)
	if embed != nil {
		path = fmt.Sprintf("%s?embed=%s", path, *embed)
	}
	log.Printf(" ==============> %s ", path)
	resp, err := s.sling.New().Get(path).Receive(order, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *order, resp, err
}

// Create creates a new order
func (s *OrderService) Create(orderBody *OrderRequest) (Order, *http.Response, error) {
	order := new(Order)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post("orders").BodyJSON(orderBody).Receive(order, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *order, resp, err
}

// Cancel will cancel an order if possible
func (s *OrderService) Cancel(orderId string) (Order, *http.Response, error) {
	order := new(Order)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Delete(fmt.Sprintf("orders/%s", orderId)).Receive(order, mollieError)
	log.Printf("+%v", mollieError.Links)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *order, resp, err
}

// Update updates an existing order
func (s *OrderService) Update(id string, orderBody *OrderUpdateRequest) (Order, *http.Response, error) {
	order := new(Order)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Patch(fmt.Sprintf("orders/%s", id)).BodyJSON(orderBody).Receive(order, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *order, resp, err
}

// UpdateLine updates an existing order line
func (s *OrderService) UpdateLine(orderId string, lineId string, lineBody *OrderLineRequest) (Order, *http.Response, error) {
	order := new(Order)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Patch(fmt.Sprintf("orders/%s/lines/%s", orderId, lineId)).BodyJSON(lineBody).Receive(order, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *order, resp, err
}

// CreatePayment creates a new payment for the order
func (s *OrderService) CreatePayment(orderId string, paymentBody *OrderPaymentRequest) (Payment, *http.Response, error) {
	payment := new(Payment)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post(fmt.Sprintf("orders/%s/payments", orderId)).BodyJSON(paymentBody).Receive(payment, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *payment, resp, err
}

// CreateRefund creates a new order refund
func (s *OrderService) CreateRefund(orderId string, refundBody *OrderRefundRequest) (OrderRefund, *http.Response, error) {
	refund := new(OrderRefund)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Post(fmt.Sprintf("orders/%s/refunds", orderId)).BodyJSON(refundBody).Receive(refund, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *refund, resp, err
}

// FetchRefund returns a order refund
func (s *OrderService) FetchRefund(orderId string, refundId string) (OrderRefund, *http.Response, error) {
	refund := new(OrderRefund)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Get(fmt.Sprintf("orders/%s/refunds/%s", orderId, refundId)).Receive(refund, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}
	return *refund, resp, err
}

// RefundList returns all order refunds created
func (s *OrderService) RefundList(orderId string, params *ListParams) (OrderRefundList, *http.Response, error) {
	refunds := new(OrderRefundList)
	mollieError := new(MollieError)
	resp, err := s.sling.New().Path(fmt.Sprintf("orders/%s/refunds", orderId)).QueryStruct(params).Receive(refunds, mollieError)
	if err == nil && mollieError.Status >= 300 {
		err = mollieError
	}

	return *refunds, resp, err
}
