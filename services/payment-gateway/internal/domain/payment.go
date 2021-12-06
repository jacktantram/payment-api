package domain

import (
	"database/sql"
	"errors"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	uuid "github.com/kevinburke/go.uuid"
	"time"
)

var (
	ErrNoPaymentForAction = errors.New("payment does not exist to create payment with")

	ErrUnprocessable = errors.New("unprocessable entity")

	ErrUpdatePaymentOutcome = errors.New("unable to update payment outcome")

	ErrNoPayment    = errors.New("no payment found")
	ErrNotPermitted = errors.New("not permitted")
)

type PaymentAction struct {
	ID           uuid.UUID      `db:"id"`
	Amount       int64          `db:"amount"`
	PaymentType  PaymentType    `db:"payment_type"`
	ResponseCode sql.NullString `db:"response_code"`
	PaymentID    uuid.UUID      `db:"payment_id"`
	CreatedAt    time.Time      `db:"created_at"`
	ProcessedAt  sql.NullTime   `db:"processed_at"`
}

type ListPaymentActionFilters struct {
	PaymentIDs []string
}

type UpdatePaymentActionField int

const (
	UpdatePaymentActionFieldResponseCode UpdatePaymentActionField = 0
)

type Payment struct {
	ID         uuid.UUID     `db:"id"`
	Amount     int64         `db:"amount"`
	Currency   string        `db:"currency"`
	Status     PaymentStatus `db:"status"`
	CardNumber string        `db:"card_number"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  sql.NullTime  `db:"updated_at"`
}

type UpdatePaymentField int

const (
	UpdatePaymentFieldStatus UpdatePaymentField = 0
)

type PaymentMethod struct {
	Card *paymentsV1.PaymentMethodCard
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
	PaymentStatusDeclined          PaymentStatus = "DECLINED"
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
	case paymentsV1.PaymentStatus_PAYMENT_STATUS_DECLINED:
		*p = PaymentStatusDeclined
	default:
		return errors.New("unknown")
	}
	return nil
}

func (p PaymentStatus) ToProto() paymentsV1.PaymentStatus {
	switch p {
	case PaymentStatusPending:
		return paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING
	case PaymentStatusAuthorized:
		return paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED
	case PaymentStatusPartiallyCaptured:
		return paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED
	case PaymentStatusCaptured:
		return paymentsV1.PaymentStatus_PAYMENT_STATUS_CAPTURED
	case PaymentStatusPartiallyRefunded:
		return paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED
	case PaymentStatusRefunded:
		return paymentsV1.PaymentStatus_PAYMENT_STATUS_REFUNDED
	case PaymentStatusVoided:
		return paymentsV1.PaymentStatus_PAYMENT_STATUS_VOIDED
	case PaymentStatusDeclined:
		return paymentsV1.PaymentStatus_PAYMENT_STATUS_DECLINED
	}
	return paymentsV1.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
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

func (p PaymentType) ToProto() paymentsV1.PaymentType {
	switch p {
	case PaymentTypeAuthorization:
		return paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION
	case PaymentTypeCapture:
		return paymentsV1.PaymentType_PAYMENT_TYPE_CAPTURE
	case PaymentTypeRefund:
		return paymentsV1.PaymentType_PAYMENT_TYPE_REFUND
	case PaymentTypeVoid:
		return paymentsV1.PaymentType_PAYMENT_TYPE_VOID
	default:
		return paymentsV1.PaymentType_PAYMENT_TYPE_UNSPECIFIED
	}
}
