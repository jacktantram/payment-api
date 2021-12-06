package store

import (
	"context"
	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/domain"
	"github.com/jmoiron/sqlx"
	uuid "github.com/kevinburke/go.uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

func (r Store) GetPayment(ctx context.Context, id string) (*paymentsV1.Payment, error) {
	var p domain.Payment
	if err := r.connFromContext(ctx).QueryRowxContext(ctx, "SELECT * FROM payment WHERE id=$1", uuid.FromStringOrNil(id)).StructScan(&p); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, domain.ErrNoPayment
		}
		return nil, err
	}

	pbPayment := &paymentsV1.Payment{
		Id: p.ID.String(),
		Amount: &amountV1.Money{
			MinorUnits: uint64(p.Amount),
			Currency:   p.Currency,
		},
		PaymentMethod: &paymentsV1.Payment_Card{
			Card: &paymentsV1.PaymentMethodCard{
				CardNumber: p.CardNumber,
			}},
		PaymentStatus: p.Status.ToProto(),
		CreatedAt:     timestamppb.New(p.CreatedAt),
	}
	if p.UpdatedAt.Valid {
		pbPayment.UpdatedAt = timestamppb.New(p.UpdatedAt.Time)
	}
	return pbPayment, nil
}

func (r Store) ListPaymentActions(ctx context.Context, filters *domain.ListPaymentActionFilters) ([]*paymentsV1.PaymentAction, error) {
	arg := map[string]interface{}{}
	if len(filters.PaymentIDs) != 0 {
		arg["payment_id"] = filters.PaymentIDs
	}
	query, args, err := sqlx.Named("SELECT * FROM payment_action WHERE payment_id=:payment_id", arg)
	if err != nil {
		return nil, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}

	query = r.db.DB.Rebind(query)
	rows, err := r.connFromContext(ctx).Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	paymentActions := make([]*paymentsV1.PaymentAction, 0)
	for rows.Next() {
		var action domain.PaymentAction
		if err := rows.StructScan(&action); err != nil {
			return nil, err
		}
		paymentAction := &paymentsV1.PaymentAction{
			Id:           action.ID.String(),
			Amount:       uint64(action.Amount),
			PaymentType:  action.PaymentType.ToProto(),
			ResponseCode: action.ResponseCode.String,
			PaymentId:    action.PaymentID.String(),
			CreatedAt:    timestamppb.New(action.CreatedAt),
		}
		if action.ProcessedAt.Valid {
			paymentAction.ProcessedAt = timestamppb.New(action.ProcessedAt.Time)
		}
		paymentActions = append(paymentActions, paymentAction)
	}
	return paymentActions, nil
}

func (r Store) CreatePayment(ctx context.Context, payment *paymentsV1.Payment) error {
	var paymentStatus domain.PaymentStatus
	if err := paymentStatus.FromProto(payment.PaymentStatus); err != nil {
		return err
	}

	rows, err := r.connFromContext(ctx).NamedQueryContext(ctx, `
		INSERT INTO payment (amount, currency, status, card_number)
		VALUES(:amount,:currency,:status,:card_number)
		RETURNING id, created_at;
		`, &domain.Payment{
		Amount:     int64(payment.Amount.MinorUnits),
		Status:     paymentStatus,
		Currency:   payment.Amount.Currency,
		CardNumber: strings.ReplaceAll(payment.GetCard().GetCardNumber(), " ", ""),
	})
	if err != nil {
		return err
	}
	if !rows.Next() {
		return errors.New("row unaffected")
	}
	var (
		id        string
		createdAt time.Time
	)
	if err = rows.Scan(&id, &createdAt); err != nil {
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
		INSERT INTO payment_action (amount, payment_type,payment_id)
		VALUES(:amount,:payment_type,:payment_id)
		RETURNING id,created_at
		`, &domain.PaymentAction{
		Amount:      int64(action.Amount),
		PaymentType: paymentType,
		PaymentID:   uuid.FromStringOrNil(action.PaymentId),
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Constraint == "payment_action_payment_id_fkey" {
				return domain.ErrNoPaymentForAction
			}
		}
		return err
	}

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

// TODO(Jack): Update these to use dynamic update statements
func (r Store) UpdatePayment(ctx context.Context, payment *paymentsV1.Payment, fields ...domain.UpdatePaymentField) error {
	var paymentStatus domain.PaymentStatus
	if err := paymentStatus.FromProto(payment.PaymentStatus); err != nil {
		return err
	}

	execContext, err := r.connFromContext(ctx).ExecContext(ctx, `UPDATE payment SET status=$1,updated_at=now() where id=$2`, paymentStatus, payment.Id)
	if err != nil {
		return err
	}
	_, err = execContext.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

// TODO(Jack): Update these to use dynamic update statements
func (r Store) UpdatePaymentAction(ctx context.Context, action *paymentsV1.PaymentAction, fields ...domain.UpdatePaymentActionField) error {
	execContext, err := r.connFromContext(ctx).ExecContext(ctx, `UPDATE payment_action SET response_code=$1, processed_at=now() where id=$2`, action.ResponseCode, action.Id)
	if err != nil {
		return err
	}
	_, err = execContext.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
