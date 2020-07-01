package data

import "testing"

func TestChecksValidation(t *testing.T) {

	p := &Product{
		Name: "Kamal",
		Price: 1.00,
		SKU: "abs-abc-efg",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}

}
