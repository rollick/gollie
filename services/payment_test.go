package services

import (
	"os"
	"testing"
)

func TestPaymentCreate(t *testing.T) {
	s := NewPaymentService(os.Getenv("MOLLIE_ACCESS_TOKEN"))
	amt := Amount{
		Currency: "EUR",
		Value:    "10.50",
	}
	p, res, err := s.Create(&PaymentRequest{
		Amount:      amt,
		Description: "A red bucket",
		RedirectUrl: "http://localhost/payment",
	})
	if err != nil {
		t.Fatalf("PaymentService.Create: %s", err)
	}

	if res.StatusCode != 201 {
		t.Errorf("PaymentService.Create: create failed. +%s", res.Status)
	}

	//defer func() {
	//	_, res, err = s.Cancel(p.ID)
	//	if err != nil {
	//		t.Errorf("PaymentService.Cancel: %s", err)
	//	}
	//	if res.StatusCode != 200 {
	//		t.Errorf("PaymentService.Cancel: cancel failed. +%v", res)
	//	}
	//}()

	if p.Amount != amt {
		t.Errorf("PaymentService.Create: amount not set correctly. +%v", p.ID)
	}
}
