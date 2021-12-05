package transporthttp_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jacktantram/payments-api/services/payment-gateway/internal/transporthttp"

	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_AuthorizeHandler_Error(t *testing.T) {
	t.Parallel()
	var (
		validRequest = transporthttp.CreateAuthorizationRequest{
			Card: &paymentsV1.PaymentMethodCard{
				CardNumber: "4000000000000119",
				Expiry: &paymentsV1.PaymentMethodCard_ExpiryDate{
					Month: transporthttp.ExpiryMonLen,
					Year:  uint32(time.Now().Year() + 1),
				},
				Cvv: "123",
			},
			Amount: amountV1.Money{
				MinorUnits: 3020,
				Currency:   "GBP",
			},
		}
	)

	for _, tc := range []struct {
		description     string
		request         transporthttp.CreateAuthorizationRequest
		responseMessage string
	}{
		{
			description: "should return error given that the amount minor units is zero",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry:     validRequest.Card.Expiry,
					Cvv:        validRequest.Card.Cvv,
				},
				Amount: amountV1.Money{
					MinorUnits: 0,
					Currency:   validRequest.Amount.Currency,
				},
			},
			responseMessage: "invalid amount.minor_units: cannot be zero",
		},
		{
			description: "should return error given that the amount currency is not a valid length",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry:     validRequest.Card.Expiry,
					Cvv:        validRequest.Card.Cvv,
				},
				Amount: amountV1.Money{
					MinorUnits: validRequest.Amount.MinorUnits,
					Currency:   "GB",
				},
			},
			responseMessage: "invalid amount.currency: must be length of 3",
		},
		{
			description: "should return error given that the card number is not a valid length",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: "12",
					Expiry:     validRequest.Card.Expiry,
					Cvv:        validRequest.Card.Cvv,
				},
				Amount: validRequest.Amount,
			},
			responseMessage: "invalid payment_method.card.card_number: length not equal to 16",
		},
		{
			description: "should return error given that the card cvv is not a valid length",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry:     validRequest.Card.Expiry,
					Cvv:        "1",
				},
				Amount: validRequest.Amount,
			},
			responseMessage: "invalid payment_method.card.cvv: length not equal to 3",
		},
		{
			description: "should return error given that the card expiry is not provided",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry:     nil,
					Cvv:        validRequest.Card.Cvv,
				},
				Amount: validRequest.Amount,
			},
			responseMessage: "missing payment_method.card.expiry: cannot be empty",
		},
		{
			description: "should return error given that the card expiry months exceeds 12",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry: &paymentsV1.PaymentMethodCard_ExpiryDate{
						Month: 13,
						Year:  validRequest.Card.Expiry.Year,
					},
					Cvv: validRequest.Card.Cvv,
				},
				Amount: validRequest.Amount,
			},
			responseMessage: "invalid payment_method.card.expiry.month: expiry month cannot exceed 12",
		},
		{
			description: "should return error given that the card expiry is before current year",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry: &paymentsV1.PaymentMethodCard_ExpiryDate{
						Month: validRequest.Card.Expiry.Month,
						Year:  uint32(time.Now().Year() - 1),
					},
					Cvv: validRequest.Card.Cvv,
				},
				Amount: validRequest.Amount,
			},
			responseMessage: "missing payment_method.card.expiry.year: cannot be in the past",
		},
	} {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			h, err := transporthttp.NewHandler()
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			b, err := json.Marshal(tc.request)
			require.NoError(t, err)

			h.AuthorizeHandler(recorder, httptest.NewRequest(http.MethodPost, "/authorize", bytes.NewReader(b)))
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			respBody, err := ioutil.ReadAll(recorder.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tc.responseMessage)
		})

	}

}

func TestHandler_AuthorizeHandler_Success(t *testing.T) {

}
