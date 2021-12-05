package processor

import (
	"context"
	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
)

type Store interface {
	CreatePayment(ctx context.Context, payment *paymentsV1.Payment) error
	CreatePaymentAction(ctx context.Context, action paymentsV1.PaymentAction) error
}

type Service struct {
}

func (s Service) CreatePayment(ctx context.Context, amount *amountV1.Money, method *paymentv1)
