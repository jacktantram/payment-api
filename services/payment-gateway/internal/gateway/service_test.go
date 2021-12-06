package gateway_test

import (
	"context"
	"github.com/golang/mock/gomock"
	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/domain"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/gateway"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/gateway/mocks"
	uuid "github.com/kevinburke/go.uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_CreatePayment_Error(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		description string
		fn          func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway)
		err         error
	}{
		{
			description: "should return an error if unable to create pending payment",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					CreatePayment(gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			description: "should return an error if unable to create payment action",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				store.
					EXPECT().
					CreatePayment(gomock.Any(), gomock.Any()).
					Return(nil)

				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			description: "should return an error if unable to call payment gateway",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					CreatePayment(gomock.Any(), gomock.Any()).
					Return(nil)

				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(nil)

				gateway.EXPECT().CreateIssuerRequest(gomock.Any(), gomock.Any()).
					Return(domain.IssuerResponse{}, errors.New("error"))

			},
			err: errors.New("error"),
		},
		{
			description: "should return an error if unable to update payment action",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					CreatePayment(gomock.Any(), gomock.Any()).
					Return(nil)

				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(nil)

				gateway.EXPECT().CreateIssuerRequest(gomock.Any(), gomock.Any()).
					Return(domain.IssuerResponse{AuthCode: "00"}, nil)

				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.EXPECT().
					UpdatePaymentAction(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("error"))

			},
			err: errors.Wrap(domain.ErrUpdatePaymentOutcome, "error"),
		},
		{
			description: "should return an error if unable to update payment",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					CreatePayment(gomock.Any(), gomock.Any()).
					Return(nil)

				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(nil)

				gateway.EXPECT().CreateIssuerRequest(gomock.Any(), gomock.Any()).
					Return(domain.IssuerResponse{AuthCode: "00"}, nil)

				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.EXPECT().
					UpdatePaymentAction(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				store.EXPECT().
					UpdatePayment(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("error"))

			},
			err: errors.Wrap(domain.ErrUpdatePaymentOutcome, "error"),
		}} {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			var (
				ctrl = gomock.NewController(t)

				mockStore         = mocks.NewMockStore(ctrl)
				mockIssuerGateway = mocks.NewMockIssuerGateway(ctrl)
			)
			if tc.fn != nil {
				tc.fn(mockStore, mockIssuerGateway)
			}
			service := gateway.NewService(mockStore, mockIssuerGateway)
			_, err := service.CreatePayment(context.Background(), &amountV1.Money{
				MinorUnits: 10000,
				Currency:   "GBP",
			}, domain.PaymentMethod{Card: &paymentsV1.PaymentMethodCard{
				CardNumber: "10000000000000000",
			},
			})
			require.Error(t, err)
			assert.Equal(t, tc.err.Error(), err.Error())
		})
	}
}

func TestService_CreatePayment_Success(t *testing.T) {
	t.Parallel()

	var (
		ctrl = gomock.NewController(t)

		store         = mocks.NewMockStore(ctrl)
		issuerGateway = mocks.NewMockIssuerGateway(ctrl)
		amount        = &amountV1.Money{
			MinorUnits: 10000,
			Currency:   "GBP",
		}
		method = domain.PaymentMethod{Card: &paymentsV1.PaymentMethodCard{
			CardNumber: "10000000000000000",
		},
		}

		paymentID = uuid.NewV4().String()
		payment   = &paymentsV1.Payment{
			Amount:        amount,
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING,
			PaymentMethod: &paymentsV1.Payment_Card{Card: method.Card},
		}
		paymentAction = &paymentsV1.PaymentAction{
			Amount:      amount.MinorUnits,
			PaymentType: paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION,
			PaymentId:   paymentID,
		}
	)

	store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})
	store.
		EXPECT().
		CreatePayment(gomock.Any(), payment).DoAndReturn(func(ctx context.Context, payment *paymentsV1.Payment) error {
		payment.Id = paymentID
		return nil
	}).
		Return(nil)

	store.
		EXPECT().
		CreatePaymentAction(gomock.Any(), paymentAction).
		Return(nil)

	issuerGateway.EXPECT().CreateIssuerRequest(gomock.Any(), domain.IssuerRequest{
		Amount:        amount,
		OperationType: paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION,
		PaymentMethod: method}).
		Return(domain.IssuerResponse{AuthCode: "00"}, nil)

	store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})

	store.EXPECT().
		UpdatePaymentAction(gomock.Any(), &paymentsV1.PaymentAction{
			Amount:       paymentAction.Amount,
			PaymentType:  paymentAction.PaymentType,
			ResponseCode: "00",
			PaymentId:    paymentAction.PaymentId,
		}, domain.UpdatePaymentActionFieldResponseCode).
		Return(nil)

	store.EXPECT().
		UpdatePayment(gomock.Any(), &paymentsV1.Payment{
			Id:            paymentID,
			Amount:        payment.Amount,
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
			PaymentMethod: payment.PaymentMethod,
		}, domain.UpdatePaymentFieldStatus).Return(nil)

	service := gateway.NewService(store, issuerGateway)
	_, err := service.CreatePayment(context.Background(), amount, method)
	require.NoError(t, err)
}

func TestService_Capture_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		description string
		amount      uint64
		fn          func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway)
		err         error
	}{
		{
			description: "given an error fetching payment",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			description: "given capture amount exceeds payment amount",
			amount:      1200,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is already captured",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_CAPTURED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is already refunded",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_REFUNDED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is partially refunded",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is voided",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_VOIDED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is pending",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is declined",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_DECLINED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is partially captured and new capture exceeds exceeds total payment amount",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)

				store.
					EXPECT().
					ListPaymentActions(gomock.Any(), gomock.Any()).
					Return([]*paymentsV1.PaymentAction{{
						Amount:       500,
						PaymentType:  paymentsV1.PaymentType_PAYMENT_TYPE_CAPTURE,
						ResponseCode: "00",
					}}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given that payment action creation fails",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			description: "given that unable to call issuer",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(nil)
				gateway.
					EXPECT().
					CreateIssuerRequest(gomock.Any(), gomock.Any()).
					Return(domain.IssuerResponse{}, errors.New("error"))

			},
			err: errors.New("error"),
		},
		{
			description: "given that unable to update payment action",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(nil)
				gateway.
					EXPECT().
					CreateIssuerRequest(gomock.Any(), gomock.Any()).
					Return(domain.IssuerResponse{}, nil)

				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.EXPECT().
					UpdatePaymentAction(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
			},
			err: errors.Wrap(domain.ErrUpdatePaymentOutcome, "error"),
		}} {

		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			var (
				ctrl = gomock.NewController(t)

				mockStore         = mocks.NewMockStore(ctrl)
				mockIssuerGateway = mocks.NewMockIssuerGateway(ctrl)
			)
			if tc.fn != nil {
				tc.fn(mockStore, mockIssuerGateway)
			}
			service := gateway.NewService(mockStore, mockIssuerGateway)
			_, err := service.Capture(context.Background(), "id", tc.amount)
			require.Error(t, err)
			assert.Equal(t, tc.err.Error(), err.Error())
		})
	}
}

func TestService_Capture_Success(t *testing.T) {
	t.Parallel()

	t.Run("should be able to capture", func(t *testing.T) {
		var (
			ctrl = gomock.NewController(t)

			store       = mocks.NewMockStore(ctrl)
			mockGateway = mocks.NewMockIssuerGateway(ctrl)
		)
		store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})
		store.
			EXPECT().
			GetPayment(gomock.Any(), gomock.Any()).
			Return(&paymentsV1.Payment{
				PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
				Amount: &amountV1.Money{
					MinorUnits: 1000,
				},
			}, nil)
		store.
			EXPECT().
			CreatePaymentAction(gomock.Any(), gomock.Any()).
			Return(nil)
		mockGateway.
			EXPECT().
			CreateIssuerRequest(gomock.Any(), gomock.Any()).
			Return(domain.IssuerResponse{AuthCode: "00"}, nil)

		store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})
		store.EXPECT().
			UpdatePaymentAction(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		store.EXPECT().
			UpdatePayment(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		service := gateway.NewService(store, mockGateway)
		payment, err := service.Capture(context.Background(), "id", 1000)
		require.NoError(t, err)
		assert.Equal(t, paymentsV1.PaymentStatus_PAYMENT_STATUS_CAPTURED, payment.PaymentStatus, payment)
	})
	t.Run("should be able to partially capture", func(t *testing.T) {
		var (
			ctrl = gomock.NewController(t)

			store       = mocks.NewMockStore(ctrl)
			mockGateway = mocks.NewMockIssuerGateway(ctrl)
		)
		store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})
		store.
			EXPECT().
			GetPayment(gomock.Any(), gomock.Any()).
			Return(&paymentsV1.Payment{
				PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
				Amount: &amountV1.Money{
					MinorUnits: 1000,
				},
			}, nil)
		store.
			EXPECT().
			CreatePaymentAction(gomock.Any(), gomock.Any()).
			Return(nil)
		mockGateway.
			EXPECT().
			CreateIssuerRequest(gomock.Any(), gomock.Any()).
			Return(domain.IssuerResponse{AuthCode: "00"}, nil)

		store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})
		store.EXPECT().
			UpdatePaymentAction(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		store.EXPECT().
			UpdatePayment(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		service := gateway.NewService(store, mockGateway)
		payment, err := service.Capture(context.Background(), "id", 500)
		require.NoError(t, err)
		assert.Equal(t, paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED, payment.PaymentStatus, payment)
	})

}

func TestService_Refund_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		description string
		amount      uint64
		fn          func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway)
		err         error
	}{
		{
			description: "given an error fetching payment",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			description: "given refund amount exceeds payment amount",
			amount:      1200,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is already refunded",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_REFUNDED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is voided",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_VOIDED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is pending",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is declined",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_DECLINED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is authorized",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is partially captured and new capture exceeds exceeds total payment amount",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)

				store.
					EXPECT().
					ListPaymentActions(gomock.Any(), gomock.Any()).
					Return([]*paymentsV1.PaymentAction{{
						Amount:       500,
						PaymentType:  paymentsV1.PaymentType_PAYMENT_TYPE_CAPTURE,
						ResponseCode: "00",
					}}, nil)
			},
			err: domain.ErrNotPermitted,
		},

		{
			description: "given payment is partially funded and refunded amount-amount already captured exceeds exceeds total payment amount",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)

				store.
					EXPECT().
					ListPaymentActions(gomock.Any(), gomock.Any()).
					Return([]*paymentsV1.PaymentAction{{
						Amount:       500,
						PaymentType:  paymentsV1.PaymentType_PAYMENT_TYPE_CAPTURE,
						ResponseCode: "00",
					},
						{
							Amount:       100,
							PaymentType:  paymentsV1.PaymentType_PAYMENT_TYPE_REFUND,
							ResponseCode: "00",
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given that payment action creation fails",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)

				store.
					EXPECT().
					ListPaymentActions(gomock.Any(), gomock.Any()).
					Return([]*paymentsV1.PaymentAction{{}}, nil)
				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			description: "given that unable to call issuer",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
				store.
					EXPECT().
					ListPaymentActions(gomock.Any(), gomock.Any()).
					Return([]*paymentsV1.PaymentAction{{}}, nil)
				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(nil)
				gateway.
					EXPECT().
					CreateIssuerRequest(gomock.Any(), gomock.Any()).
					Return(domain.IssuerResponse{}, errors.New("error"))

			},
			err: errors.New("error"),
		},
		{
			description: "given that unable to update payment action",
			amount:      1000,
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
				store.
					EXPECT().
					ListPaymentActions(gomock.Any(), gomock.Any()).
					Return([]*paymentsV1.PaymentAction{{}}, nil)
				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(nil)
				gateway.
					EXPECT().
					CreateIssuerRequest(gomock.Any(), gomock.Any()).
					Return(domain.IssuerResponse{}, nil)

				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.EXPECT().
					UpdatePaymentAction(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
			},
			err: errors.Wrap(domain.ErrUpdatePaymentOutcome, "error"),
		}} {

		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			var (
				ctrl = gomock.NewController(t)

				mockStore         = mocks.NewMockStore(ctrl)
				mockIssuerGateway = mocks.NewMockIssuerGateway(ctrl)
			)
			if tc.fn != nil {
				tc.fn(mockStore, mockIssuerGateway)
			}
			service := gateway.NewService(mockStore, mockIssuerGateway)
			_, err := service.Refund(context.Background(), "id", tc.amount)
			require.Error(t, err)
			assert.Equal(t, tc.err.Error(), err.Error())
		})
	}
}

func TestService_Void_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		description string
		fn          func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway)
		err         error
	}{
		{
			description: "given an error fetching payment",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			description: "given payment is refunded",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_REFUNDED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is partially refunded",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is already voided",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_VOIDED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is pending",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is declined",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_DECLINED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is captured",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_CAPTURED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given payment is partially captured",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_CAPTURED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
			},
			err: domain.ErrNotPermitted,
		},
		{
			description: "given that payment action creation fails",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)

				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			description: "given that unable to call issuer",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(nil)
				gateway.
					EXPECT().
					CreateIssuerRequest(gomock.Any(), gomock.Any()).
					Return(domain.IssuerResponse{}, errors.New("error"))

			},
			err: errors.New("error"),
		},
		{
			description: "given that unable to update payment action",
			fn: func(store *mocks.MockStore, gateway *mocks.MockIssuerGateway) {
				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.
					EXPECT().
					GetPayment(gomock.Any(), gomock.Any()).
					Return(&paymentsV1.Payment{
						PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
						Amount: &amountV1.Money{
							MinorUnits: 1000,
						},
					}, nil)
				store.
					EXPECT().
					CreatePaymentAction(gomock.Any(), gomock.Any()).
					Return(nil)
				gateway.
					EXPECT().
					CreateIssuerRequest(gomock.Any(), gomock.Any()).
					Return(domain.IssuerResponse{}, nil)

				store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})
				store.EXPECT().
					UpdatePaymentAction(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("error"))
			},
			err: errors.Wrap(domain.ErrUpdatePaymentOutcome, "error"),
		}} {

		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			var (
				ctrl = gomock.NewController(t)

				mockStore         = mocks.NewMockStore(ctrl)
				mockIssuerGateway = mocks.NewMockIssuerGateway(ctrl)
			)
			if tc.fn != nil {
				tc.fn(mockStore, mockIssuerGateway)
			}
			service := gateway.NewService(mockStore, mockIssuerGateway)
			_, err := service.Void(context.Background(), "id")
			require.Error(t, err)
			assert.Equal(t, tc.err.Error(), err.Error())
		})
	}
}

func TestService_Refund_Success(t *testing.T) {
	t.Parallel()
	t.Run("should be able to fully refund", func(t *testing.T) {
		var (
			ctrl = gomock.NewController(t)

			store             = mocks.NewMockStore(ctrl)
			mockIssuerGateway = mocks.NewMockIssuerGateway(ctrl)
		)

		store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})
		store.
			EXPECT().
			GetPayment(gomock.Any(), gomock.Any()).
			Return(&paymentsV1.Payment{
				PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED,
				Amount: &amountV1.Money{
					MinorUnits: 1000,
				},
			}, nil)
		store.
			EXPECT().
			ListPaymentActions(gomock.Any(), gomock.Any()).
			Return([]*paymentsV1.PaymentAction{{}}, nil)
		store.
			EXPECT().
			CreatePaymentAction(gomock.Any(), gomock.Any()).
			Return(nil)
		mockIssuerGateway.
			EXPECT().
			CreateIssuerRequest(gomock.Any(), gomock.Any()).
			Return(domain.IssuerResponse{AuthCode: "00"}, nil)

		store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})
		store.EXPECT().
			UpdatePaymentAction(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		store.EXPECT().
			UpdatePayment(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		service := gateway.NewService(store, mockIssuerGateway)
		payment, err := service.Refund(context.Background(), "id", 1000)
		require.NoError(t, err)
		assert.Equal(t, paymentsV1.PaymentStatus_PAYMENT_STATUS_REFUNDED, payment.PaymentStatus)
	})
	t.Run("should be able to partially refund", func(t *testing.T) {
		var (
			ctrl = gomock.NewController(t)

			store             = mocks.NewMockStore(ctrl)
			mockIssuerGateway = mocks.NewMockIssuerGateway(ctrl)
		)

		store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})
		store.
			EXPECT().
			GetPayment(gomock.Any(), gomock.Any()).
			Return(&paymentsV1.Payment{
				PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED,
				Amount: &amountV1.Money{
					MinorUnits: 1000,
				},
			}, nil)
		store.
			EXPECT().
			ListPaymentActions(gomock.Any(), gomock.Any()).
			Return([]*paymentsV1.PaymentAction{{}}, nil)
		store.
			EXPECT().
			CreatePaymentAction(gomock.Any(), gomock.Any()).
			Return(nil)
		mockIssuerGateway.
			EXPECT().
			CreateIssuerRequest(gomock.Any(), gomock.Any()).
			Return(domain.IssuerResponse{AuthCode: "00"}, nil)

		store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
				return fn(ctx)
			})
		store.EXPECT().
			UpdatePaymentAction(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		store.EXPECT().
			UpdatePayment(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		service := gateway.NewService(store, mockIssuerGateway)
		payment, err := service.Refund(context.Background(), "id", 500)
		require.NoError(t, err)
		assert.Equal(t, paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED, payment.PaymentStatus)
	})

}

func TestService_Void_Success(t *testing.T) {
	t.Parallel()

	var (
		ctrl = gomock.NewController(t)

		store             = mocks.NewMockStore(ctrl)
		mockIssuerGateway = mocks.NewMockIssuerGateway(ctrl)
	)

	store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})
	store.
		EXPECT().
		GetPayment(gomock.Any(), gomock.Any()).
		Return(&paymentsV1.Payment{
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
			Amount: &amountV1.Money{
				MinorUnits: 1000,
			},
		}, nil)

	store.
		EXPECT().
		CreatePaymentAction(gomock.Any(), gomock.Any()).
		Return(nil)
	mockIssuerGateway.
		EXPECT().
		CreateIssuerRequest(gomock.Any(), gomock.Any()).
		Return(domain.IssuerResponse{AuthCode: "00"}, nil)

	store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})
	store.EXPECT().
		UpdatePaymentAction(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	store.EXPECT().
		UpdatePayment(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	service := gateway.NewService(store, mockIssuerGateway)
	payment, err := service.Void(context.Background(), "id")
	require.NoError(t, err)
	assert.Equal(t, paymentsV1.PaymentStatus_PAYMENT_STATUS_VOIDED, payment.PaymentStatus)
}
