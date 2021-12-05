package store

import (
	"context"
	"fmt"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	"github.com/jacktantram/payments-api/services/payment-processor/internal/domain"
	uuid "github.com/kevinburke/go.uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (r Store) CreatePayment(ctx context.Context, payment *paymentsV1.Payment) error {
	var paymentStatus domain.PaymentStatus
	if err := paymentStatus.FromProto(payment.PaymentStatus); err != nil {
		return err
	}
	rows, err := r.connFromContext(ctx).NamedQueryContext(ctx, `
		INSERT INTO payment (amount, currency, status, action_id)
		VALUES(:amount,:currency,:status,:action_id)
		RETURNING (id, created_at);
		`, &domain.Payment{
		Amount:   int64(payment.Amount.MinorUnits),
		Status:   paymentStatus,
		Currency: payment.Amount.Currency,
		ActionID: uuid.FromStringOrNil(payment.ActionId),
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Constraint == "payment_action_id_fkey" {
				return domain.ErrPaymentCreateActionDoesNotExist
			}
		}

		return err
	}
	if !rows.Next() {
		return errors.New("row unaffected")
	}
	var (
		id        string
		createdAt time.Time
	)
	if err = rows.Scan(&id, createdAt); err != nil {
		return errors.Wrap(err, "unable to scan row")
	}
	payment.Id = id
	payment.CreatedAt = timestamppb.New(createdAt)
	return nil
}

func (r Store) CreatePaymentAction(ctx context.Context, action *paymentsV1.PaymentAction) error {
	var paymentType domain.PaymentType
	if err := paymentType.FromProto(action.PaymentType); err != nil {
		return err
	}
	rows, err := r.connFromContext(ctx).NamedQueryContext(ctx, `
		INSERT INTO payment_action (amount, payment_type)
		VALUES(:amount,:payment_type)
		RETURNING id,created_at
		`, &domain.PaymentAction{
		Amount:      int64(action.Amount.MinorUnits),
		PaymentType: paymentType,
	})
	if err != nil {
		return err
	}
	fmt.Println(rows)

	if !rows.Next() {
		return errors.New("row unaffected")
	}
	var (
		id        uuid.UUID
		createdAt time.Time
	)
	if err = rows.Scan(&id, &createdAt); err != nil {
		return errors.Wrap(err, "unable to scan row")
	}
	action.Id = id.String()
	action.CreatedAt = timestamppb.New(createdAt)
	return nil
}
