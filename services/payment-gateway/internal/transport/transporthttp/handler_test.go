package transporthttp_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/domain"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jacktantram/payments-api/services/payment-gateway/internal/transport/transporthttp"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/transport/transporthttp/mocks"

	"github.com/golang/mock/gomock"
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
		validRequest = transporthttp.CreateAuthorizationRequest{
			Card: &paymentsV1.PaymentMethodCard{
				CardNumber: "4000000000000119",
				Expiry: &paymentsV1.PaymentMethodCard_ExpiryDate{
					Month: transporthttp.ExpiryMonLen,
					Year:  uint32(time.Now().Year() + 1),
				},
				Cvv: "123",
			},
			Amount: &amountV1.Money{
				MinorUnits: 3020,
				Currency:   "GBP",
			},
		}
	)

	for _, tc := range []struct {
		description     string
		request         transporthttp.CreateAuthorizationRequest
		expStatusCode   int
		responseMessage string
		fn              func(mocks *mocks.MockGateway)
	}{
		{
			description: "should return error given that the amount minor units is zero",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry:     validRequest.Card.Expiry,
					Cvv:        validRequest.Card.Cvv,
				},
				Amount: &amountV1.Money{
					MinorUnits: 0,
					Currency:   validRequest.Amount.Currency,
				},
			},
			responseMessage: "invalid amount.minor_units: cannot be zero",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the amount currency is not a valid length",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry:     validRequest.Card.Expiry,
					Cvv:        validRequest.Card.Cvv,
				},
				Amount: &amountV1.Money{
					MinorUnits: validRequest.Amount.MinorUnits,
					Currency:   "GB",
				},
			},
			responseMessage: "invalid amount.currency: must be length of 3",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the card number is not a valid length",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: "12",
					Expiry:     validRequest.Card.Expiry,
					Cvv:        validRequest.Card.Cvv,
				},
				Amount: &amountV1.Money{
					MinorUnits: validRequest.Amount.MinorUnits,
					Currency:   validRequest.Amount.Currency,
				},
			},
			responseMessage: "invalid payment_method.card.card_number: length not equal to 16",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the card cvv is not a valid length",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry:     validRequest.Card.Expiry,
					Cvv:        "1",
				},
				Amount: &amountV1.Money{
					MinorUnits: validRequest.Amount.MinorUnits,
					Currency:   validRequest.Amount.Currency,
				},
			},
			responseMessage: "invalid payment_method.card.cvv: length not equal to 3",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the card expiry is not provided",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry:     nil,
					Cvv:        validRequest.Card.Cvv,
				},
				Amount: &amountV1.Money{
					MinorUnits: validRequest.Amount.MinorUnits,
					Currency:   validRequest.Amount.Currency,
				},
			},
			responseMessage: "missing payment_method.card.expiry: cannot be empty",
			expStatusCode:   http.StatusUnprocessableEntity,
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
				Amount: &amountV1.Money{
					MinorUnits: validRequest.Amount.MinorUnits,
					Currency:   validRequest.Amount.Currency,
				},
			},
			responseMessage: "invalid payment_method.card.expiry.month: expiry month cannot exceed 12",
			expStatusCode:   http.StatusUnprocessableEntity,
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
				Amount: &amountV1.Money{
					MinorUnits: validRequest.Amount.MinorUnits,
					Currency:   validRequest.Amount.Currency,
				},
			},
			responseMessage: "missing payment_method.card.expiry.year: cannot be in the past",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error if unable to create payment",
			request: transporthttp.CreateAuthorizationRequest{
				Card: &paymentsV1.PaymentMethodCard{
					CardNumber: validRequest.Card.CardNumber,
					Expiry: &paymentsV1.PaymentMethodCard_ExpiryDate{
						Month: validRequest.Card.Expiry.Month,
						Year:  validRequest.Card.Expiry.Year,
					},
					Cvv: validRequest.Card.Cvv,
				},
				Amount: &amountV1.Money{
					MinorUnits: validRequest.Amount.MinorUnits,
					Currency:   validRequest.Amount.Currency,
				},
			},
			responseMessage: "Oops something went wrong",
			fn: func(mocks *mocks.MockGateway) {
				mocks.
					EXPECT().
					CreatePayment(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("an error"))
			},
			expStatusCode: http.StatusInternalServerError,
		},
	} {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			var (
				ctrl        = gomock.NewController(t)
				mockGateway = mocks.NewMockGateway(ctrl)
			)
			if tc.fn != nil {
				tc.fn(mockGateway)
			}

			h, err := transporthttp.NewHandler(mockGateway)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			b, err := json.Marshal(&tc.request)
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
	//go:embed testdata/valid-authorization-request.json
	validAuthorizationRequest []byte
)

func TestHandler_AuthorizeHandler_Success(t *testing.T) {
	t.Parallel()
	var (
		ctrl        = gomock.NewController(t)
		mockGateway = mocks.NewMockGateway(ctrl)

		expPayment = &paymentsV1.Payment{
			Id: uuid.NewV4().String(),
			Amount: &amountV1.Money{
				MinorUnits: 1000,
				Currency:   "GBP",
			},
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
			CreatedAt:     timestamppb.Now(),
			UpdatedAt:     nil,
		}
	)

	mockGateway.EXPECT().CreatePayment(gomock.Any(), &amountV1.Money{
		MinorUnits: 2212,
		Currency:   "GBP",
	},
		domain.PaymentMethod{
			Card: &paymentsV1.PaymentMethodCard{
				CardNumber: "4000 0000 0000 0119",
				Expiry: &paymentsV1.PaymentMethodCard_ExpiryDate{
					Month: 12,
					Year:  2029,
				},
				Cvv: "123",
			}}).
		Return(expPayment, nil)

	h, err := transporthttp.NewHandler(mockGateway)
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

func TestHandler_CaptureHandler_Error(t *testing.T) {
	t.Parallel()
	var (
		validRequest = transporthttp.CreateCaptureRequest{
			PaymentID: uuid.NewV4().String(),
			Amount:    3020,
		}
	)

	for _, tc := range []struct {
		description     string
		request         transporthttp.CreateCaptureRequest
		expStatusCode   int
		responseMessage string
		fn              func(mocks *mocks.MockGateway)
	}{
		{
			description: "should return error given that the payment id is empty",
			request: transporthttp.CreateCaptureRequest{
				PaymentID: "",
				Amount:    validRequest.Amount,
			},

			responseMessage: "invalid payment_id: cannot be empty",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the amount is zero",
			request: transporthttp.CreateCaptureRequest{
				PaymentID: validRequest.PaymentID,
				Amount:    0,
			},
			responseMessage: "invalid amount: cannot be zero",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description:     "should return error if unable to create capture",
			request:         validRequest,
			responseMessage: "Oops something went wrong",
			fn: func(mocks *mocks.MockGateway) {
				mocks.
					EXPECT().
					Capture(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("an error"))
			},
			expStatusCode: http.StatusInternalServerError,
		},
		{
			description:     "should return error if unable to payment not found",
			request:         validRequest,
			responseMessage: "payment not found",
			fn: func(mocks *mocks.MockGateway) {
				mocks.
					EXPECT().
					Capture(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, domain.ErrNoPayment)
			},
			expStatusCode: http.StatusNotFound,
		},
		{
			description:     "should return error if capture not allowed",
			request:         validRequest,
			responseMessage: "capture not allowed",
			fn: func(mocks *mocks.MockGateway) {
				mocks.
					EXPECT().
					Capture(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, domain.ErrNotPermitted)
			},
			expStatusCode: http.StatusForbidden,
		},
	} {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			var (
				ctrl        = gomock.NewController(t)
				mockGateway = mocks.NewMockGateway(ctrl)
			)
			if tc.fn != nil {
				tc.fn(mockGateway)
			}

			h, err := transporthttp.NewHandler(mockGateway)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			b, err := json.Marshal(&tc.request)
			require.NoError(t, err)

			h.CaptureHandler(recorder, httptest.NewRequest(http.MethodPost, "/capture", bytes.NewReader(b)))
			assert.Equal(t, tc.expStatusCode, recorder.Code)
			respBody, err := ioutil.ReadAll(recorder.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tc.responseMessage)
		})

	}
}

var (
	//go:embed testdata/valid-capture-request.json
	validCaptureRequest []byte
)

func TestHandler_CaptureHandler_Success(t *testing.T) {
	t.Parallel()
	var (
		ctrl        = gomock.NewController(t)
		mockGateway = mocks.NewMockGateway(ctrl)

		expPayment = &paymentsV1.Payment{
			Id: uuid.NewV4().String(),
			Amount: &amountV1.Money{
				MinorUnits: 1000,
				Currency:   "GBP",
			},
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED,
			CreatedAt:     timestamppb.Now(),
			UpdatedAt:     nil,
		}
	)

	mockGateway.EXPECT().Capture(gomock.Any(), "a6921fc3-a7e3-4661-909b-b3c6c77837ce", uint64(2212)).
		Return(expPayment, nil)

	h, err := transporthttp.NewHandler(mockGateway)
	require.NoError(t, err)
	recorder := httptest.NewRecorder()

	h.CaptureHandler(recorder, httptest.NewRequest(http.MethodPost, "/capture", bytes.NewReader(validCaptureRequest)))
	assert.Equal(t, http.StatusOK, recorder.Code)
	respBody, err := ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	var paymentResponse paymentsV1.Payment
	require.NoError(t, protojson.Unmarshal(respBody, &paymentResponse))

	// easier to compare when test fails. proto.Equal also works but not as readable
	assert.Equal(t, expPayment.String(), paymentResponse.String())
}

func TestHandler_RefundHandler_Error(t *testing.T) {
	t.Parallel()
	var (
		validRequest = transporthttp.CreateRefundRequest{
			PaymentID: uuid.NewV4().String(),
			Amount:    3020,
		}
	)

	for _, tc := range []struct {
		description     string
		request         transporthttp.CreateRefundRequest
		expStatusCode   int
		responseMessage string
		fn              func(mocks *mocks.MockGateway)
	}{
		{
			description: "should return error given that the payment id is empty",
			request: transporthttp.CreateRefundRequest{
				PaymentID: "",
				Amount:    validRequest.Amount,
			},

			responseMessage: "invalid payment_id: cannot be empty",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description: "should return error given that the amount is zero",
			request: transporthttp.CreateRefundRequest{
				PaymentID: validRequest.PaymentID,
				Amount:    0,
			},
			responseMessage: "invalid amount: cannot be zero",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description:     "should return error if unable to create refund",
			request:         validRequest,
			responseMessage: "Oops something went wrong",
			fn: func(mocks *mocks.MockGateway) {
				mocks.
					EXPECT().
					Refund(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("an error"))
			},
			expStatusCode: http.StatusInternalServerError,
		},
		{
			description:     "should return error if unable to payment not found",
			request:         validRequest,
			responseMessage: "payment not found",
			fn: func(mocks *mocks.MockGateway) {
				mocks.
					EXPECT().
					Refund(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, domain.ErrNoPayment)
			},
			expStatusCode: http.StatusNotFound,
		},
		{
			description:     "should return error if refund not allowed",
			request:         validRequest,
			responseMessage: "refund not allowed",
			fn: func(mocks *mocks.MockGateway) {
				mocks.
					EXPECT().
					Refund(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, domain.ErrNotPermitted)
			},
			expStatusCode: http.StatusForbidden,
		},
	} {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			var (
				ctrl        = gomock.NewController(t)
				mockGateway = mocks.NewMockGateway(ctrl)
			)
			if tc.fn != nil {
				tc.fn(mockGateway)
			}

			h, err := transporthttp.NewHandler(mockGateway)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			b, err := json.Marshal(&tc.request)
			require.NoError(t, err)

			h.RefundHandler(recorder, httptest.NewRequest(http.MethodPost, "/refund", bytes.NewReader(b)))
			assert.Equal(t, tc.expStatusCode, recorder.Code)
			respBody, err := ioutil.ReadAll(recorder.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tc.responseMessage)
		})

	}
}

var (
	//go:embed testdata/valid-refund-request.json
	validRefundRequest []byte
)

func TestHandler_RefundHandler_Success(t *testing.T) {
	t.Parallel()
	var (
		ctrl        = gomock.NewController(t)
		mockGateway = mocks.NewMockGateway(ctrl)

		expPayment = &paymentsV1.Payment{
			Id: uuid.NewV4().String(),
			Amount: &amountV1.Money{
				MinorUnits: 1000,
				Currency:   "GBP",
			},
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED,
			CreatedAt:     timestamppb.Now(),
			UpdatedAt:     nil,
		}
	)

	mockGateway.EXPECT().Refund(gomock.Any(), "a6921fc3-a7e3-4661-909b-b3c6c77837ce", uint64(2212)).
		Return(expPayment, nil)

	h, err := transporthttp.NewHandler(mockGateway)
	require.NoError(t, err)
	recorder := httptest.NewRecorder()

	h.RefundHandler(recorder, httptest.NewRequest(http.MethodPost, "/refund", bytes.NewReader(validRefundRequest)))
	assert.Equal(t, http.StatusOK, recorder.Code)
	respBody, err := ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	var paymentResponse paymentsV1.Payment
	require.NoError(t, protojson.Unmarshal(respBody, &paymentResponse))

	// easier to compare when test fails. proto.Equal also works but not as readable
	assert.Equal(t, expPayment.String(), paymentResponse.String())
}

func TestHandler_VoidHandler_Error(t *testing.T) {
	t.Parallel()
	var (
		validRequest = transporthttp.CreateVoidRequest{
			PaymentID: uuid.NewV4().String(),
		}
	)

	for _, tc := range []struct {
		description     string
		request         transporthttp.CreateVoidRequest
		expStatusCode   int
		responseMessage string
		fn              func(mocks *mocks.MockGateway)
	}{
		{
			description: "should return error given that the payment id is empty",
			request: transporthttp.CreateVoidRequest{
				PaymentID: "",
			},

			responseMessage: "invalid payment_id: cannot be empty",
			expStatusCode:   http.StatusUnprocessableEntity,
		},
		{
			description:     "should return error if unable to perform void",
			request:         validRequest,
			responseMessage: "Oops something went wrong",
			fn: func(mocks *mocks.MockGateway) {
				mocks.
					EXPECT().
					Void(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("an error"))
			},
			expStatusCode: http.StatusInternalServerError,
		},
		{
			description:     "should return error if unable to payment not found",
			request:         validRequest,
			responseMessage: "payment not found",
			fn: func(mocks *mocks.MockGateway) {
				mocks.
					EXPECT().
					Void(gomock.Any(), gomock.Any()).
					Return(nil, domain.ErrNoPayment)
			},
			expStatusCode: http.StatusNotFound,
		},
		{
			description:     "should return error if void not allowed",
			request:         validRequest,
			responseMessage: "void not allowed",
			fn: func(mocks *mocks.MockGateway) {
				mocks.
					EXPECT().
					Void(gomock.Any(), gomock.Any()).
					Return(nil, domain.ErrNotPermitted)
			},
			expStatusCode: http.StatusForbidden,
		},
	} {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			var (
				ctrl        = gomock.NewController(t)
				mockGateway = mocks.NewMockGateway(ctrl)
			)
			if tc.fn != nil {
				tc.fn(mockGateway)
			}

			h, err := transporthttp.NewHandler(mockGateway)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()

			b, err := json.Marshal(&tc.request)
			require.NoError(t, err)

			h.VoidHandler(recorder, httptest.NewRequest(http.MethodPost, "/void", bytes.NewReader(b)))
			assert.Equal(t, tc.expStatusCode, recorder.Code)
			respBody, err := ioutil.ReadAll(recorder.Body)
			require.NoError(t, err)
			assert.Contains(t, string(respBody), tc.responseMessage)
		})
	}
}

var (
	//go:embed testdata/valid-void-request.json
	validVoidRequest []byte
)

func TestHandler_VoidHandler_Success(t *testing.T) {
	t.Parallel()
	var (
		ctrl        = gomock.NewController(t)
		mockGateway = mocks.NewMockGateway(ctrl)

		expPayment = &paymentsV1.Payment{
			Id: uuid.NewV4().String(),
			Amount: &amountV1.Money{
				MinorUnits: 1000,
				Currency:   "GBP",
			},
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_VOIDED,
			CreatedAt:     timestamppb.Now(),
			UpdatedAt:     nil,
		}
	)

	mockGateway.EXPECT().Void(gomock.Any(), "a6921fc3-a7e3-4661-909b-b3c6c77837ce").
		Return(expPayment, nil)

	h, err := transporthttp.NewHandler(mockGateway)
	require.NoError(t, err)
	recorder := httptest.NewRecorder()

	h.VoidHandler(recorder, httptest.NewRequest(http.MethodPost, "/void", bytes.NewReader(validVoidRequest)))
	assert.Equal(t, http.StatusOK, recorder.Code)
	respBody, err := ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	var paymentResponse paymentsV1.Payment
	require.NoError(t, protojson.Unmarshal(respBody, &paymentResponse))

	// easier to compare when test fails. proto.Equal also works but not as readable
	assert.Equal(t, expPayment.String(), paymentResponse.String())
}
