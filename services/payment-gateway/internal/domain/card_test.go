package domain_test

import (
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateCardNumber(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		cardNumber string
		bool
	}{
		{cardNumber: "4603 1110 9388 0019", bool: true},
		{cardNumber: "4603111093880019", bool: true},
		{cardNumber: "4603111093880019", bool: true},
		{cardNumber: "5555555555554444", bool: true},
		{cardNumber: "4000 0000 0000 0119", bool: true},
		{cardNumber: "1111111111111111", bool: false},
		{cardNumber: "1841835786578528", bool: false},
		{cardNumber: "1841 8357 8657 8528", bool: false},
	} {
		tc := tc
		t.Run(tc.cardNumber, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.bool, domain.ValidCardNumber(tc.cardNumber))
		})
	}
}
