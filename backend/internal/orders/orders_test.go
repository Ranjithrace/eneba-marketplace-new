package orders

import "testing"

func TestNormalizeAndFingerprint(t *testing.T) {
	req := CheckoutRequest{Items: []Line{{OfferID: 9, Quantity: 1}, {OfferID: 3, Quantity: 2}, {OfferID: 9, Quantity: 1}}, PaymentMethod: "approve", IdempotencyKey: "request-abcdef123"}
	out, e := Normalize(req)
	if e != nil {
		t.Fatal(e)
	}
	if len(out.Items) != 2 || out.Items[0].OfferID != 3 || out.Items[1].Quantity != 2 {
		t.Fatalf("not normalized: %+v", out)
	}
	a, _ := Fingerprint(req)
	b, _ := Fingerprint(CheckoutRequest{Items: []Line{{OfferID: 3, Quantity: 2}, {OfferID: 9, Quantity: 2}}, PaymentMethod: "approve", IdempotencyKey: "different-key-123"})
	if a != b {
		t.Fatal("equivalent carts differ")
	}
	other := req
	other.PaymentMethod = "decline"
	c, _ := Fingerprint(other)
	if a == c {
		t.Fatal("different payment paths share fingerprint")
	}
}
func TestRejectMalformedCheckout(t *testing.T) {
	cases := []CheckoutRequest{
		{Items: []Line{{OfferID: 1, Quantity: 21}}, PaymentMethod: "approve", IdempotencyKey: "request-abcdef123"},
		{Items: []Line{{OfferID: 0, Quantity: 1}}, PaymentMethod: "approve", IdempotencyKey: "request-abcdef123"},
		{Items: []Line{{OfferID: 1, Quantity: 1}}, PaymentMethod: "real_card", IdempotencyKey: "request-abcdef123"},
		{Items: []Line{{OfferID: 1, Quantity: 1}}, PaymentMethod: "approve", IdempotencyKey: "short"},
	}
	for _, c := range cases {
		if _, e := Normalize(c); e == nil {
			t.Errorf("accepted: %+v", c)
		}
	}
}
