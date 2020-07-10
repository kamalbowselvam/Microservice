package data

import (
		"testing"
		"github.com/stretchr/testify/assert"
)

func TestProductMissingNameReturnsErr(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.22,
		SKU: "abc-efg-hij",
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 0)
}