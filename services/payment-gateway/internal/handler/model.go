package handler

import (
	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
)

// CreateAuthorizationRequest is the request used to create an authorization
type CreateAuthorizationRequest struct {
	Card   *paymentsV1.PaymentMethodCard
	Amount *amountV1.Money
}

// CreateCaptureRequest  is the request used to perform a capture towards a payment.
type CreateCaptureRequest struct {
	PaymentID string `json:"payment_id"`
	Amount    uint64 `json:"amount"`
}

// CreateRefundRequest is the request used to create a refund towards a payment
type CreateRefundRequest struct {
	PaymentID string `json:"payment_id"`
	Amount    uint64 `json:"amount"`
}

// CreateVoidRequest is the request used to void a payment.
type CreateVoidRequest struct {
	PaymentID string `json:"payment_id"`
}
