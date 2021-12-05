package handler_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jacktantram/payments-api/services/payment-gateway/internal/handler"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/handler/mocks"

	"github.com/golang/mock/gomock"
	processorv1 "github.com/jacktantram/payments-api/build/go/rpc/paymentprocessor/v1"
	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	uuid "github.com/kevinburke/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHandler_AuthorizeHandler_Error(t *testing.T) {
	t.Parallel()
	var (
		validRequest = handler.CreateAuthorizationRequest{
			Card: &paymentsV1.PaymentMethodCard{
				CardNumber: "4000000000000119",
				Expiry: &paymentsV1.PaymentMethodCard_ExpiryDate{
					Month: handler.ExpiryMonLen,
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
		request         handler.CreateAuthorizationRequest
		expStatusCode   int
		responseMessage string
		fn              func(mocks *mocks.MockPaymentProcessorClient)
	}{
		{
			description: "should return error given that the amount minor units is zero",
			request: handler.CreateAuthorizationRequest{
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
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the amount currency is not a valid length",
			request: handler.CreateAuthorizationRequest{
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
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the card number is not a valid length",
			request: handler.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: "12",
					Expiry:     validRequest.Card.Expiry,
					Cvv:        validRequest.Card.Cvv,
				},
				Amount: validRequest.Amount,
			},
			responseMessage: "invalid payment_method.card.card_number: length not equal to 16",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the card cvv is not a valid length",
			request: handler.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry:     validRequest.Card.Expiry,
					Cvv:        "1",
				},
				Amount: validRequest.Amount,
			},
			responseMessage: "invalid payment_method.card.cvv: length not equal to 3",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the card expiry is not provided",
			request: handler.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry:     nil,
					Cvv:        validRequest.Card.Cvv,
				},
				Amount: validRequest.Amount,
			},
			responseMessage: "missing payment_method.card.expiry: cannot be empty",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the card expiry months exceeds 12",
			request: handler.CreateAuthorizationRequest{
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
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the card expiry is before current year",
			request: handler.CreateAuthorizationRequest{
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
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description:     "should return error if unable to create payment",
			request:         validRequest,
			responseMessage: "Oops something went wrong",
			fn: func(mocks *mocks.MockPaymentProcessorClient) {
				mocks.
					EXPECT().
					CreatePayment(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("an error"))
			},
			expStatusCode: http.StatusInternalServerError,
		},
	} {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			var (
				ctrl                = gomock.NewController(t)
				mockProcessorClient = mocks.NewMockPaymentProcessorClient(ctrl)
			)
			if tc.fn != nil {
				tc.fn(mockProcessorClient)
			}

			h, err := handler.NewHandler(mockProcessorClient)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			b, err := json.Marshal(tc.request)
			require.NoError(t, err)

			h.AuthorizeHandler(recorder, httptest.NewRequest(http.MethodPost, "/authorize", bytes.NewReader(b)))
			assert.Equal(t, tc.expStatusCode, recorder.Code)
			respBody, err := ioutil.ReadAll(recorder.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tc.responseMessage)
		})

	}

}

var (
	//go:embed testdata/authorization/valid-authorization-request.json
	validAuthorizationRequest []byte
)

func TestHandler_AuthorizeHandler_Success(t *testing.T) {
	t.Parallel()
	var (
		ctrl                = gomock.NewController(t)
		mockProcessorClient = mocks.NewMockPaymentProcessorClient(ctrl)

		expPayment = &paymentsV1.Payment{
			Id: uuid.NewV4().String(),
			Amount: &amountV1.Money{
				MinorUnits: 1000,
				Currency:   "GBP",
			},
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
			ActionId:      uuid.NewV4().String(),
			CreatedAt:     timestamppb.Now(),
			UpdatedAt:     nil,
		}
	)

	mockProcessorClient.EXPECT().CreatePayment(gomock.Any(), &processorv1.CreatePaymentRequest{
		Amount: &amountV1.Money{
			MinorUnits: 2212,
			Currency:   "GBP",
		},
		PaymentMethod: &processorv1.CreatePaymentRequest_Card{
			Card: &paymentsV1.PaymentMethodCard{
				CardNumber: "4000 0000 0000 0119",
				Expiry: &paymentsV1.PaymentMethodCard_ExpiryDate{
					Month: 12,
					Year:  2029,
				},
				Cvv: "123",
			},
		}}).
		Return(&processorv1.CreatePaymentResponse{Payment: expPayment}, nil)

	h, err := handler.NewHandler(mockProcessorClient)
	require.NoError(t, err)
	recorder := httptest.NewRecorder()

	h.AuthorizeHandler(recorder, httptest.NewRequest(http.MethodPost, "/authorize", bytes.NewReader(validAuthorizationRequest)))
	assert.Equal(t, http.StatusOK, recorder.Code)
	respBody, err := ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	var paymentResponse paymentsV1.Payment
	require.NoError(t, protojson.Unmarshal(respBody, &paymentResponse))

	// easier to compare when test fails. proto.Equal also works but not as readable
	assert.Equal(t, expPayment.String(), paymentResponse.String())
}
