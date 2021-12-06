package gateway_test

import (
	"context"
	"github.com/golang/mock/gomock"
	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/domain"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/gateway"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/gateway/mocks"
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
	//t.Parallel()
	//var (
	//	ctrl = gomock.NewController(t)
	//
	//	store         = mocks.NewMockStore(ctrl)
	//	issuerGateway = mocks.NewMockIssuerGateway(ctrl)
	//	amount        = &amountV1.Money{
	//		MinorUnits: 10000,
	//		Currency:   "GBP",
	//	}
	//	method = domain.PaymentMethod{Card: &paymentsV1.PaymentMethodCard{
	//		CardNumber: "10000000000000000",
	//	},
	//	}
	//	paymentType   = paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION
	//	paymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED
	//
	//	payment = &paymentsV1.Payment{
	//		Amount:        amount,
	//		PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING,
	//		PaymentMethod: &paymentsV1.Payment_Card{Card: method.Card},
	//	}
	//	paymentAction = paymentsV1.PaymentAction{
	//		Amount:      amount.MinorUnits,
	//		PaymentType: paymentType,
	//		PaymentId:   uuid.NewV4().String(),
	//	}
	//)
	//store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
	//	DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
	//		return fn(ctx)
	//	})
	//store.
	//	EXPECT().
	//	CreatePayment(gomock.Any(), payment).
	//	Return(nil)
	//
	//store.
	//	EXPECT().
	//	CreatePaymentAction(gomock.Any(), paymentAction).
	//	Return(nil)
	//
	//issuerGateway.EXPECT().CreateIssuerRequest(gomock.Any(), domain.IssuerRequest{
	//	Amount:        amount,
	//	OperationType: paymentType,
	//	PaymentMethod: method}).
	//	Return(domain.IssuerResponse{AuthCode: "00"}, nil)
	//
	//store.EXPECT().ExecInTransaction(gomock.Any(), gomock.Any()).
	//	DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
	//		return fn(ctx)
	//	})
	//
	//paymentAction.PaymentType = paymentType
	//store.EXPECT().
	//	UpdatePaymentAction(gomock.Any(), paymentAction, paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED).
	//	Return(nil)
	//
	//payment.PaymentStatus = paymentStatus
	//store.EXPECT().
	//	UpdatePayment(gomock.Any(), payment, paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED).
	//	Return(nil)
	//
	//service := gateway.NewService(store, issuerGateway)
	//_, err := service.CreatePayment(context.Background(), amount, method)
	//require.NoError(t, err)
}
