package domain

import (
	"database/sql"
	"errors"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	uuid "github.com/kevinburke/go.uuid"
	"time"
)

var (
	ErrPaymentCreateActionDoesNotExist = errors.New("action does not exist to create payment with")
)

type PaymentAction struct {
	ID          uuid.UUID    `db:"id"`
	Amount      int64        `db:"amount"`
	PaymentType PaymentType  `db:"payment_type"`
	CreatedAt   time.Time    `db:"created_at"`
	ProcessedAt sql.NullTime `db:"updated_at"`
}

type Payment struct {
	ID          uuid.UUID     `db:"id"`
	Amount      int64         `db:"amount"`
	Currency    string        `db:"currency"`
	Status      PaymentStatus `db:"status"`
	ActionID    uuid.UUID     `db:"action_id"`
	CreatedAt   time.Time     `db:"created_at"`
	ProcessedAt sql.NullTime  `db:"updated_at"`
}

type PaymentStatus string

const (
	PaymentStatusPending           PaymentStatus = "PENDING"
	PaymentStatusAuthorized        PaymentStatus = "AUTHORIZED"
	PaymentStatusCaptured          PaymentStatus = "CAPTURED"
	PaymentStatusPartiallyCaptured PaymentStatus = "PARTIALLY_CAPTURED"
	PaymentStatusRefunded          PaymentStatus = "REFUNDED"
	PaymentStatusPartiallyRefunded PaymentStatus = "PARTIALLY_REFUNDED"
	PaymentStatusVoided            PaymentStatus = "VOIDED"
)

func (p *PaymentStatus) FromProto(paymentStatus paymentsV1.PaymentStatus) error {
	switch paymentStatus {
	case paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING:
		*p = PaymentStatusPending
	case paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED:
		*p = PaymentStatusAuthorized
	case paymentsV1.PaymentStatus_PAYMENT_STATUS_CAPTURED:
		*p = PaymentStatusCaptured
	case paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED:
		*p = PaymentStatusPartiallyCaptured
	case paymentsV1.PaymentStatus_PAYMENT_STATUS_REFUNDED:
		*p = PaymentStatusRefunded
	case paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED:
		*p = PaymentStatusPartiallyRefunded
	case paymentsV1.PaymentStatus_PAYMENT_STATUS_VOIDED:
		*p = PaymentStatusVoided
	default:
		return errors.New("unknown")
	}
	return nil
}

type PaymentType string

const (
	PaymentTypeAuthorization PaymentType = "AUTHORIZATION"
	PaymentTypeCapture       PaymentType = "CAPTURE"
	PaymentTypeRefund        PaymentType = "REFUND"
	PaymentTypeVoid          PaymentType = "VOID"
)

func (p *PaymentType) FromProto(paymentType paymentsV1.PaymentType) error {
	switch paymentType {
	case paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION:
		*p = PaymentTypeAuthorization
	case paymentsV1.PaymentType_PAYMENT_TYPE_CAPTURE:
		*p = PaymentTypeCapture
	case paymentsV1.PaymentType_PAYMENT_TYPE_REFUND:
		*p = PaymentTypeRefund
	case paymentsV1.PaymentType_PAYMENT_TYPE_VOID:
		*p = PaymentTypeVoid
	default:
		return errors.New("unknown")
	}
	return nil
}
