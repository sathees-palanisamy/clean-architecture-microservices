package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct_AvailableQty(t *testing.T) {
	p := &Product{
		TotalQty:    10,
		ReservedQty: 3,
	}

	assert.Equal(t, 7, p.AvailableQty())
}
