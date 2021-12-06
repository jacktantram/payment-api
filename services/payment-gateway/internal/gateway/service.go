//go:generate mockgen -source=service.go -destination=mocks/mocks.go -package=mocks

package gateway

import (
	"context"
	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/domain"
	"github.com/pkg/errors"
)

type Store interface {
	ExecInTransaction(ctx context.Context, fn func(ctx context.Context) error) error

	GetPayment(ctx context.Context, id string) (*paymentsV1.Payment, error)
	ListPaymentActions(ctx context.Context, filters *domain.ListPaymentActionFilters) ([]*paymentsV1.PaymentAction, error)

	CreatePayment(ctx context.Context, payment *paymentsV1.Payment) error
	CreatePaymentAction(ctx context.Context, action *paymentsV1.PaymentAction) error

	UpdatePayment(ctx context.Context, payment *paymentsV1.Payment, fields ...domain.UpdatePaymentField) error
	UpdatePaymentAction(ctx context.Context, action *paymentsV1.PaymentAction, fields ...domain.UpdatePaymentActionField) error
}

type IssuerGateway interface {
	CreateIssuerRequest(ctx context.Context, issuerRequest domain.IssuerRequest) (domain.IssuerResponse, error)
}

type Service struct {
	store         Store
	issuerGateway IssuerGateway
}

func NewService(store Store, gateway IssuerGateway) Service {
	return Service{store: store, issuerGateway: gateway}
}

func (s Service) CreatePayment(ctx context.Context, amount *amountV1.Money, method domain.PaymentMethod) (*paymentsV1.Payment, error) {
	paymentType := paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION
	var (
		payment       *paymentsV1.Payment
		paymentAction *paymentsV1.PaymentAction
	)
	if err := s.store.ExecInTransaction(ctx, func(ctx context.Context) error {
		payment = &paymentsV1.Payment{
			Amount:        amount,
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING,
			PaymentMethod: &paymentsV1.Payment_Card{Card: method.Card},
		}
		if err := s.store.CreatePayment(ctx, payment); err != nil {
			return err
		}
		paymentAction = &paymentsV1.PaymentAction{
			Amount:      amount.MinorUnits,
			PaymentType: paymentType,
			PaymentId:   payment.Id,
		}
		if err := s.store.CreatePaymentAction(ctx, paymentAction); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// If more time would have done this part asynchronously from a PaymentCreatedEvent.
	issuerResponse, err := s.issuerGateway.CreateIssuerRequest(ctx, domain.IssuerRequest{
		Amount:        amount,
		OperationType: paymentType,
		PaymentMethod: method})
	if err != nil {
		return nil, err
	}

	if err = s.store.ExecInTransaction(ctx, func(ctx context.Context) error {
		paymentAction.ResponseCode = issuerResponse.AuthCode
		if err = s.store.UpdatePaymentAction(ctx, paymentAction, domain.UpdatePaymentActionFieldResponseCode); err != nil {
			return err
		}

		// This will need more work on mappings
		if issuerSuccess(issuerResponse.AuthCode) {
			payment.PaymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED
		} else {
			payment.PaymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_DECLINED
		}
		if err := s.store.UpdatePayment(ctx, payment, domain.UpdatePaymentFieldStatus); err != nil {
			return err
		}
		return nil
	}); err != nil {
		// will need to alert on this as payment was successful
		return nil, errors.Wrap(domain.ErrUpdatePaymentOutcome, err.Error())
	}

	return payment, nil
}

// Capture is responsible for capturing funds in a payment.
// The amount cannot exceed the existing auth amount and cannot capture unless the payment is in an authorized or partially captured state.
// It can also not exceed the existing successful payment action amounts
func (s Service) Capture(ctx context.Context, paymentID string, amount uint64) (*paymentsV1.Payment, error) {
	paymentType := paymentsV1.PaymentType_PAYMENT_TYPE_CAPTURE
	var (
		payment       *paymentsV1.Payment
		paymentAction *paymentsV1.PaymentAction
		sumAction     uint64
	)
	if err := s.store.ExecInTransaction(ctx, func(ctx context.Context) error {
		var err error
		payment, err = s.store.GetPayment(ctx, paymentID)
		if err != nil {
			return err
		}

		if amount > payment.Amount.GetMinorUnits() {
			return domain.ErrNotPermitted
		}

		if payment.PaymentStatus != paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED && payment.PaymentStatus != paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED {
			return domain.ErrNotPermitted

		}

		if payment.PaymentStatus == paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED {
			actions, err := s.store.ListPaymentActions(ctx, &domain.ListPaymentActionFilters{PaymentIDs: []string{paymentID}})
			if err != nil {
				return err
			}

			for _, action := range actions {
				if issuerSuccess(action.ResponseCode) {
					if action.PaymentType == paymentsV1.PaymentType_PAYMENT_TYPE_CAPTURE {
						sumAction += action.Amount
					}
				}
			}
			if amount > payment.Amount.MinorUnits-sumAction {
				return domain.ErrNotPermitted
			}
		}

		paymentAction = &paymentsV1.PaymentAction{
			Amount:      amount,
			PaymentType: paymentType,
			PaymentId:   paymentID,
		}
		if err = s.store.CreatePaymentAction(ctx, paymentAction); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	issuerResponse, err := s.issuerGateway.CreateIssuerRequest(ctx, domain.IssuerRequest{
		Amount: &amountV1.Money{
			MinorUnits: amount,
			Currency:   payment.Amount.Currency,
		},
		OperationType: paymentType,
		PaymentMethod: domain.PaymentMethod{Card: payment.GetCard()}})
	if err != nil {
		return nil, err
	}

	if err = s.store.ExecInTransaction(ctx, func(ctx context.Context) error {
		paymentAction.ResponseCode = issuerResponse.AuthCode
		if err = s.store.UpdatePaymentAction(ctx, paymentAction, domain.UpdatePaymentActionFieldResponseCode); err != nil {
			return err
		}

		// This will need more work on mappings
		if issuerSuccess(issuerResponse.AuthCode) {
			if sumAction+amount == payment.Amount.MinorUnits {
				payment.PaymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_CAPTURED
			} else {
				payment.PaymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED
			}
			if err := s.store.UpdatePayment(ctx, payment, domain.UpdatePaymentFieldStatus); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {

		// will need to alert on this as payment was successful
		return nil, errors.Wrap(domain.ErrUpdatePaymentOutcome, err.Error())
	}
	return payment, nil
}

func (s Service) Refund(ctx context.Context, paymentID string, amount uint64) (*paymentsV1.Payment, error) {
	paymentType := paymentsV1.PaymentType_PAYMENT_TYPE_REFUND
	var (
		payment       *paymentsV1.Payment
		paymentAction *paymentsV1.PaymentAction
		sumAction     uint64
	)
	if err := s.store.ExecInTransaction(ctx, func(ctx context.Context) error {
		var err error
		payment, err = s.store.GetPayment(ctx, paymentID)
		if err != nil {
			return err
		}

		if amount > payment.Amount.GetMinorUnits() {
			return domain.ErrNotPermitted
		}

		if payment.PaymentStatus != paymentsV1.PaymentStatus_PAYMENT_STATUS_CAPTURED && payment.PaymentStatus != paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED && payment.PaymentStatus != paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED {
			return domain.ErrNotPermitted
		}

		actions, err := s.store.ListPaymentActions(ctx, &domain.ListPaymentActionFilters{PaymentIDs: []string{paymentID}})
		if err != nil {
			return err
		}

		for _, action := range actions {
			if issuerSuccess(action.ResponseCode) {
				if action.PaymentType == paymentsV1.PaymentType_PAYMENT_TYPE_CAPTURE {
					sumAction += action.Amount
				}
				if action.PaymentType == paymentsV1.PaymentType_PAYMENT_TYPE_REFUND {
					sumAction -= action.Amount
				}
			}
		}
		if amount > sumAction && sumAction > 0 {
			return domain.ErrNotPermitted
		}
		paymentAction = &paymentsV1.PaymentAction{
			Amount:      amount,
			PaymentType: paymentType,
			PaymentId:   paymentID,
		}
		if err = s.store.CreatePaymentAction(ctx, paymentAction); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	issuerResponse, err := s.issuerGateway.CreateIssuerRequest(ctx, domain.IssuerRequest{
		Amount: &amountV1.Money{
			MinorUnits: amount,
			Currency:   payment.Amount.Currency,
		},
		OperationType: paymentType,
		PaymentMethod: domain.PaymentMethod{Card: payment.GetCard()}})
	if err != nil {
		return nil, err
	}

	if err = s.store.ExecInTransaction(ctx, func(ctx context.Context) error {
		paymentAction.ResponseCode = issuerResponse.AuthCode
		if err = s.store.UpdatePaymentAction(ctx, paymentAction, domain.UpdatePaymentActionFieldResponseCode); err != nil {
			return err
		}

		// This will need more work on mappings
		if issuerSuccess(issuerResponse.AuthCode) {
			if sumAction+amount == payment.Amount.MinorUnits || amount == sumAction {
				payment.PaymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_REFUNDED
			} else {
				payment.PaymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED
			}
			if err := s.store.UpdatePayment(ctx, payment, domain.UpdatePaymentFieldStatus); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		// will need to alert on this as payment was successful
		return nil, errors.Wrap(domain.ErrUpdatePaymentOutcome, err.Error())
	}
	return payment, nil
}

func (s Service) Void(ctx context.Context, paymentID string) (*paymentsV1.Payment, error) {
	paymentType := paymentsV1.PaymentType_PAYMENT_TYPE_VOID
	var (
		payment       *paymentsV1.Payment
		paymentAction *paymentsV1.PaymentAction
	)
	if err := s.store.ExecInTransaction(ctx, func(ctx context.Context) error {
		var err error
		payment, err = s.store.GetPayment(ctx, paymentID)
		if err != nil {
			return err
		}

		if payment.PaymentStatus != paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED {
			return domain.ErrNotPermitted
		}

		paymentAction = &paymentsV1.PaymentAction{
			Amount:      payment.Amount.GetMinorUnits(),
			PaymentType: paymentType,
			PaymentId:   paymentID,
		}
		if err = s.store.CreatePaymentAction(ctx, paymentAction); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	issuerResponse, err := s.issuerGateway.CreateIssuerRequest(ctx, domain.IssuerRequest{
		Amount: &amountV1.Money{
			MinorUnits: payment.Amount.GetMinorUnits(),
			Currency:   payment.Amount.Currency,
		},
		OperationType: paymentType,
		PaymentMethod: domain.PaymentMethod{Card: payment.GetCard()}})
	if err != nil {
		return nil, err
	}

	if err = s.store.ExecInTransaction(ctx, func(ctx context.Context) error {
		paymentAction.ResponseCode = issuerResponse.AuthCode
		if err = s.store.UpdatePaymentAction(ctx, paymentAction, domain.UpdatePaymentActionFieldResponseCode); err != nil {
			return err
		}
		// This will need more work on mappings
		if issuerSuccess(issuerResponse.AuthCode) {
			payment.PaymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_VOIDED
			if err = s.store.UpdatePayment(ctx, payment, domain.UpdatePaymentFieldStatus); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		// will need to alert on this as payment was successful
		return nil, errors.Wrap(domain.ErrUpdatePaymentOutcome, err.Error())
	}
	return payment, nil
}

func issuerSuccess(code string) bool {
	return code == "00"
}
