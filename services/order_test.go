package services

import (
	"fmt"
	"testing"
)

func createOrder(r OrderRequest) (*Order, error) {
	s := testClient().OrderService
	o, res, err := s.Create(&r)
	if err != nil {
		return nil, fmt.Errorf("OrderService.Create: %s", err)
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf("OrderService.Create: create failed. +%s", res.Status)
	}

	return &o, err
}

func TestOrderCreate(t *testing.T) {
	s := testClient().OrderService
	r := testOrderRequest()
	o, err := createOrder(r)
	if err != nil {
		t.Fatalf("OrderService.Create: %s", err)
	}

	defer func() {
		_, res, err := s.Cancel(o.ID)
		if err != nil {
			t.Errorf("OrderService.Cancel: %s", err)
		}
		if res.StatusCode != 200 {
			t.Errorf("OrderService.Cancel: cancel failed. +%v", res)
		}
	}()

	if o.Amount != r.Amount {
		t.Errorf("OrderService.Create: amount not set correctly. +%v", o.ID)
	}
}

func TestOrderCreateWithPayment(t *testing.T) {
	s := testClient().OrderService
	r := testOrderRequest()
	p := OrderRequestPayment{
		CustomerReference: "abc123",
		WebhookURL:        "http://www.brown.com/hook",
	}

	r.Payment = p

	o, err := createOrder(r)
	if err != nil {
		t.Fatalf("OrderService.Create: %s", err)
	}

	defer func() {
		_, res, err := s.Cancel(o.ID)
		if err != nil {
			t.Errorf("OrderService.Cancel: %s", err)
		}
		if res.StatusCode != 200 {
			t.Errorf("OrderService.Cancel: cancel failed. +%v", res)
		}
	}()

	if o.Amount != r.Amount {
		t.Errorf("OrderService.Create: amount should be +%v", r.Amount)
	}

	if r.Payment.CustomerReference != p.CustomerReference {
		t.Errorf("OrderService.Create: payment.customerReference should be %s", p.CustomerReference)
	}
}

func TestOrderCreateWithDiscount(t *testing.T) {
	s := testClient().OrderService
	r := testOrderRequest()
	q := uint8(1)
	r.Lines = append(r.Lines, &OrderLineRequest{
		Name:     "Discount Code",
		Type:     "discount",
		Quantity: &q,
		UnitPrice: &Amount{
			Currency: "EUR",
			Value:    "-5.00",
		},
		TotalAmount: &Amount{
			Currency: "EUR",
			Value:    "-5.00",
		},
		VatRate: "21.00",
		VatAmount: &Amount{
			Currency: "EUR",
			Value:    "-0.87",
		},
	})

	// 1) create with incorrect total
	_, err := createOrder(r)
	if err == nil {
		t.Errorf("OrderService.Create: should not be possible to create discounted order with incorrect total. %s", err)
	}

	// 2) create with correct total
	r.Amount = Amount{
		Currency: "EUR",
		Value:    "5.50",
	}

	o, err := createOrder(r)
	if err != nil {
		t.Fatalf("OrderService.Create: %s", err)
	}

	defer func() {
		_, res, err := s.Cancel(o.ID)
		if err != nil {
			t.Errorf("OrderService.Cancel: %s", err)
		}
		if res.StatusCode != 200 {
			t.Errorf("OrderService.Cancel: cancel failed. +%v", res)
		}
	}()

	if o.Amount != r.Amount {
		t.Errorf("OrderService.Create: amount should be +%v", r.Amount)
	}

}

func TestOrderCreateWithGiftCard(t *testing.T) {
	s := testClient().OrderService
	r := testOrderRequest()
	q := uint8(1)
	r.Lines = append(r.Lines, &OrderLineRequest{
		Name:     "10 Euro gift card",
		Quantity: &q,
		UnitPrice: &Amount{
			Currency: "EUR",
			Value:    "-5.00",
		},
		TotalAmount: &Amount{
			Currency: "EUR",
			Value:    "-5.00",
		},
		VatRate: "21.00",
		VatAmount: &Amount{
			Currency: "EUR",
			Value:    "-0.87",
		},
	})

	// 1) create with incorrect total
	_, err := createOrder(r)
	if err == nil {
		t.Errorf("OrderService.Create: should not be possible to create discounted order with incorrect total. %s", err)
	}

	// 2) create with correct total
	r.Amount = Amount{
		Currency: "EUR",
		Value:    "5.50",
	}

	o, err := createOrder(r)
	if err != nil {
		t.Fatalf("OrderService.Create: %s", err)
	}

	defer func() {
		_, res, err := s.Cancel(o.ID)
		if err != nil {
			t.Errorf("OrderService.Cancel: %s", err)
		}
		if res.StatusCode != 200 {
			t.Errorf("OrderService.Cancel: cancel failed. +%v", res)
		}
	}()

	if o.Amount != r.Amount {
		t.Errorf("OrderService.Create: amount should be +%v", r.Amount)
	}

}

func TestOrderUpdate(t *testing.T) {
	s := testClient().OrderService
	r := testOrderRequest()
	o, err := createOrder(r)
	if err != nil {
		t.Fatalf("OrderService.Create: %s", err)
	}

	ba := *o.BillingAddress
	ba.City = "Lakeside"

	ur := OrderUpdateRequest{
		BillingAddress: &ba,
		OrderNumber:    "abcde12345",
	}

	no, _, err := s.Update(o.ID, &ur)
	if err != nil {
		t.Fatalf("OrderService.Update: %s", err)
	}

	defer func() {
		_, res, err := s.Cancel(o.ID)
		if err != nil {
			t.Errorf("OrderService.Cancel: %s", err)
		}
		if res.StatusCode != 200 {
			t.Errorf("OrderService.Cancel: cancel failed. +%v", res)
		}
	}()

	if no.OrderNumber != ur.OrderNumber {
		t.Errorf("OrderService.Update: orderNumber should be +%v", ur.OrderNumber)
	}

	if no.BillingAddress.City != ur.BillingAddress.City {
		t.Errorf("OrderService.Update: billingAddress.city should be +%v", ur.BillingAddress.City)
	}

	if no.ShippingAddress.City != o.ShippingAddress.City {
		t.Errorf("OrderService.Update: shippingAddress.city should be +%v", o.ShippingAddress.City)
	}
}

func TestOrderLineUpdate(t *testing.T) {
	s := testClient().OrderService
	r := testOrderRequest()

	o, err := createOrder(r)
	if err != nil {
		t.Fatalf("OrderService.Create: %s", err)
	}

	l := o.Lines[0]
	ol := OrderLineRequest{
		Name: "New line name",
	}

	no, res, err := s.UpdateLine(o.ID, *l.ID, &ol)
	if err != nil {
		t.Fatalf("OrderService.Update: %s", err)
	}
	if res.StatusCode != 200 {
		t.Errorf("OrderService.Cancel: cancel failed. +%v", res)
	}

	defer func() {
		_, res, err := s.Cancel(o.ID)
		if err != nil {
			t.Errorf("OrderService.Cancel: %s", err)
		}
		if res.StatusCode != 200 {
			t.Errorf("OrderService.Cancel: cancel failed. +%v", res)
		}
	}()

	if no.Lines[0].Name != ol.Name {
		t.Errorf("OrderService.UpdateOrderLine: name should be +%v", ol.Name)
	}
}

func TestOrderPaymentCreate(t *testing.T) {
	tc := testClient()
	s := tc.OrderService
	r := testOrderRequest()
	o, err := createOrder(r)
	if err != nil {
		t.Fatalf("OrderService.Create: %s", err)
	}

	// Cancel the existing order payment before trying to
	// create a new payment.
	e := "payments"
	fo, _, err := s.Fetch(o.ID, &e)
	if err != nil {
		t.Fatalf("OrderService.Fetch: %s", err)
	}

	if fo.Embedded.Payments != nil {
		for _, p := range fo.Embedded.Payments {
			pa, _, err := tc.PaymentService.Cancel(p.ID)
			if err != nil {
				t.Fatalf("PaymentService.Cancel: %s", err)
			}
			if pa.Status != "canceled" {
				t.Errorf("PaymentService.Cancel: payment not canceled %s. Status is %s", p.ID, p.Status)
			}
		}
	}

	//pr := OrderPaymentRequest{}
	//_, _, err = s.CreatePayment(o.ID, &pr)
	//if err != nil {
	//	t.Fatalf("OrderService.CreatePayment: %s", err)
	//}
}

func TestOrderList(t *testing.T) {
	s := testClient().OrderService
	l := uint16(1)
	p := ListParams{
		Limit: &l,
	}

	_, _, err := s.List(&p)
	if err != nil {
		t.Fatalf("OrderService.List: %s", err)
	}
}

func TestOrderListWithParams(t *testing.T) {
	s := testClient().OrderService
	r := testOrderRequest()
	o, err := createOrder(r)
	if err != nil {
		t.Fatalf("OrderService.Create: %s", err)
	}

	l := uint16(1)
	p := ListParams{
		From:  &o.ID,
		Limit: &l,
	}

	ol, _, err := s.List(&p)
	if err != nil {
		t.Errorf("OrderService.List: %s", err)
	}

	if len(ol.Embedded.Orders) == 0 {
		t.Error("OrderService.List: at least one order expected")
	}

	f := ol.Embedded.Orders[0]
	if f.ID != o.ID {
		t.Errorf("OrderService.List: first order should be %s. Got: %s", o.ID, f.ID)
	}
}
