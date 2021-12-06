package domain

import (
	v1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
)

type IssuerRequest struct {
	Amount        *v1.Money
	OperationType paymentsV1.PaymentType
	PaymentMethod PaymentMethod
}

type IssuerResponse struct {
	AuthCode string
}
