package catalog

import "testing"

func TestFilterValidation(t *testing.T) {
	for _, f := range []Filters{{MinPrice: -1}, {MinPrice: 3000, MaxPrice: 1000}, {Search: string(make([]byte, 101))}, {Page: 10001}} {
		if Validate(f) == nil {
			t.Errorf("accepted invalid filter: %+v", f)
		}
	}
	if err := Validate(Filters{Search: "minecrft", MinPrice: 0, MaxPrice: 3000, Page: 2}); err != nil {
		t.Fatal(err)
	}
}
