package transporthttp

import (
	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
)

type CreateAuthorizationRequest struct{
	Card *paymentsV1.PaymentMethodCard
	Amount amountV1.Money
}



type CreatePaymentResponse struct{
	Payment *paymentsV1.Payment
}