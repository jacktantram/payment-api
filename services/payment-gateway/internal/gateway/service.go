package gateway

import (
	"context"
	processorv1 "github.com/jacktantram/payments-api/build/go/rpc/paymentprocessor/v1"
	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/domain"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Store interface {
	ExecInTransaction(ctx context.Context, fn func(ctx context.Context) error) error

	GetPayment(ctx context.Context, id string) (*paymentsV1.Payment, error)
	ListPaymentActions(ctx context.Context, filters *processorv1.ListPaymentActionFilters) ([]*paymentsV1.PaymentAction, error)

	CreatePayment(ctx context.Context, payment *paymentsV1.Payment) error
	CreatePaymentAction(ctx context.Context, action *paymentsV1.PaymentAction) error

	UpdatePayment(ctx context.Context, payment *paymentsV1.Payment, fields ...processorv1.UpdatePaymentField) error
	UpdatePaymentAction(ctx context.Context, action *paymentsV1.PaymentAction, fields ...processorv1.UpdatePaymentActionField) error
}

type IssuerGateway interface {
	CreateIssuerRequest(ctx context.Context, issuerRequest domain.IssuerRequest) (domain.IssuerResponse, error)
}

type Service struct {
	store         Store
	issuerGateway IssuerGateway
	processorv1.UnimplementedPaymentProcessorServer
	S *grpc.Server
}

func NewService(store Store, gateway IssuerGateway) Service {
	s := Service{S: grpc.NewServer(), store: store, issuerGateway: gateway}
	processorv1.RegisterPaymentProcessorServer(s.S, s)
	return s
}

func (s Service) CreatePayment(ctx context.Context, amount *amountV1.Money, method domain.PaymentMethod) (*processorv1.CreatePaymentResponse, error) {
	paymentType := paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION
	var (
		payment       *paymentsV1.Payment
		paymentAction *paymentsV1.PaymentAction
	)
	if err := s.store.ExecInTransaction(ctx, func(ctx context.Context) error {
		paymentAction = &paymentsV1.PaymentAction{
			Amount:      amount.MinorUnits,
			PaymentType: paymentType,
		}
		if err := s.store.CreatePaymentAction(ctx, paymentAction); err != nil {
			return err
		}

		payment = &paymentsV1.Payment{
			Amount:        amount,
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING,
		}
		if err := s.store.CreatePayment(ctx, payment); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	issuerRequest :=

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
		if err = s.store.UpdatePaymentAction(ctx, paymentAction, processorv1.UpdatePaymentActionField_UPDATE_PAYMENT_ACTION_RESPONSE_CODE); err != nil {
			return err
		}

		// This will need more work on mappings
		if issuerSuccess(issuerResponse.AuthCode) {
			payment.PaymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED
		} else {
			payment.PaymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_DECLINED
		}
		if err = s.store.UpdatePayment(ctx, payment, processorv1.UpdatePaymentField_UPDATE_PAYMENT_FIELD_STATUS); err != nil {
			return err
		}
		return nil
	}); err != nil {
		// will need to alert on this as payment was successful
		return nil, errors.Wrap(domain.ErrUpdatePaymentOutcome, err.Error())
	}

	return &processorv1.CreatePaymentResponse{Payment: payment}, nil
}

// Capture is responsible for capturing funds in a payment.
// The amount cannot exceed the existing auth amount and cannot capture unless the payment is in an authorized or partially captured state.
// It can also not exceed the existing successful payment action amounts
func (s Service) Capture(ctx context.Context, request *processorv1.CreateCaptureRequest) (*processorv1.CreateCaptureResponse, error) {
	paymentType := paymentsV1.PaymentType_PAYMENT_TYPE_CAPTURE
	var (
		payment       *paymentsV1.Payment
		paymentAction *paymentsV1.PaymentAction
		sumAction     uint64
		amount        = request.GetAmount()
	)
	if err := s.store.ExecInTransaction(ctx, func(ctx context.Context) error {
		payment, err := s.store.GetPayment(ctx, request.GetPaymentId())
		if err != nil {
			return err
		}

		if amount > payment.Amount.GetMinorUnits() {
			return domain.ErrInvalidAmount
		}
		switch payment.PaymentStatus {
		case paymentsV1.PaymentStatus_PAYMENT_STATUS_DECLINED, paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED, paymentsV1.PaymentStatus_PAYMENT_STATUS_REFUNDED, paymentsV1.PaymentStatus_PAYMENT_STATUS_VOIDED:
			return domain.ErrUnprocessable
		}

		actions, err := s.store.ListPaymentActions(ctx, &processorv1.ListPaymentActionFilters{ActionIds: []string{paymentAction.Id}})
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
		if amount > sumAction {
			return domain.ErrInvalidAmount
		}
		if err := s.store.CreatePaymentAction(ctx, &paymentsV1.PaymentAction{
			Amount:      amount,
			PaymentType: paymentType,
		}); err != nil {
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
	})
	if err != nil {
		return nil, err
	}

	if err = s.store.ExecInTransaction(ctx, func(ctx context.Context) error {
		paymentAction.ResponseCode = issuerResponse.AuthCode
		if err = s.store.UpdatePaymentAction(ctx, paymentAction, processorv1.UpdatePaymentActionField_UPDATE_PAYMENT_ACTION_RESPONSE_CODE); err != nil {
			return err
		}

		// This will need more work on mappings
		if issuerSuccess(issuerResponse.AuthCode) {
			if sumAction+amount == payment.Amount.MinorUnits {
				payment.PaymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_CAPTURED
			} else {
				payment.PaymentStatus = paymentsV1.PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED
			}
			if err := s.store.UpdatePayment(ctx, payment, processorv1.UpdatePaymentField_UPDATE_PAYMENT_FIELD_STATUS); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		// will need to alert on this as payment was successful
		return nil, errors.Wrap(domain.ErrUpdatePaymentOutcome, err.Error())
	}
	return &processorv1.CreateCaptureResponse{Payment: payment}, nil
}

func (s Service) Refund(ctx context.Context, request *processorv1.CreateRefundRequest) (*processorv1.CreateRefundResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Void(ctx context.Context, request *processorv1.CreateVoidRequest) (*processorv1.CreateVoidResponse, error) {
	//TODO implement me
	panic("implement me")
}

func issuerSuccess(code string) bool {
	return code == "00"
}
