// +build integration

package store_test

import (
	"context"
	processorv1 "github.com/jacktantram/payments-api/build/go/rpc/paymentprocessor/v1"
	amountV1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	"github.com/jacktantram/payments-api/services/payment-processor/internal/domain"
	uuid "github.com/kevinburke/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_CreatePaymentAction(t *testing.T) {
	t.Parallel()
	var (
		paymentAction = &paymentsV1.PaymentAction{
			Amount:      1021,
			PaymentType: paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION,
		}
	)

	require.NoError(t, testStore.CreatePaymentAction(context.Background(), paymentAction))
	assert.NotEmpty(t, paymentAction.Id)
	assert.NotNil(t, paymentAction.CreatedAt)
}

func TestStore_CreatePayment(t *testing.T) {
	t.Parallel()
	t.Run("should fail creating a payment given that the action does not exist", func(t *testing.T) {
		payment := &paymentsV1.Payment{
			Amount: &amountV1.Money{
				MinorUnits: 1000,
				Currency:   "GBP",
			},
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_AUTHORIZED,
			ActionId:      uuid.NewV4().String(),
		}
		err := testStore.CreatePayment(context.Background(), payment)
		require.Error(t, err)
		assert.Equal(t, domain.ErrPaymentCreateActionDoesNotExist, err)
	})
	t.Run("should successfully create a payment given there is a payment action", func(t *testing.T) {
		var (
			paymentAction = &paymentsV1.PaymentAction{
				Amount:      1021,
				PaymentType: paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION,
			}
		)
		require.NoError(t, testStore.CreatePaymentAction(context.Background(), paymentAction))

		payment := &paymentsV1.Payment{
			Amount: &amountV1.Money{
				MinorUnits: 1000,
				Currency:   "GBP",
			},
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING,
			ActionId:      paymentAction.Id,
		}
		require.NoError(t, testStore.CreatePayment(context.Background(), payment))

		assert.NotEmpty(t, payment.Id)
		assert.NotNil(t, payment.CreatedAt)
	})
}

func TestStore_GetPayment(t *testing.T) {
	t.Parallel()
	t.Run("should successfully get a payment", func(t *testing.T) {
		var (
			paymentAction = &paymentsV1.PaymentAction{
				Amount:      1021,
				PaymentType: paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION,
			}
		)
		require.NoError(t, testStore.CreatePaymentAction(context.Background(), paymentAction))

		payment := &paymentsV1.Payment{
			Amount: &amountV1.Money{
				MinorUnits: 1000,
				Currency:   "GBP",
			},
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING,
			ActionId:      paymentAction.Id,
		}
		require.NoError(t, testStore.CreatePayment(context.Background(), payment))

		p, err := testStore.GetPayment(context.Background(), payment.Id)
		require.NoError(t, err)
		assert.Equal(t, payment, p)
	})
	t.Run("should return error for unknown payment", func(t *testing.T) {
		_, err := testStore.GetPayment(context.Background(), uuid.NewV4().String())
		require.Error(t, err)
		assert.Equal(t, domain.ErrNoPayment, err)
	})
}

func TestStore_ListPaymentActions(t *testing.T) {
	t.Parallel()
	t.Run("should successfully get a list of payment actions", func(t *testing.T) {
		var (
			paymentAction = &paymentsV1.PaymentAction{
				Amount:      1021,
				PaymentType: paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION,
			}
		)
		require.NoError(t, testStore.CreatePaymentAction(context.Background(), paymentAction))

		payment := &paymentsV1.Payment{
			Amount: &amountV1.Money{
				MinorUnits: 1000,
				Currency:   "GBP",
			},
			PaymentStatus: paymentsV1.PaymentStatus_PAYMENT_STATUS_PENDING,
			ActionId:      paymentAction.Id,
		}
		require.NoError(t, testStore.CreatePayment(context.Background(), payment))

		p, err := testStore.ListPaymentActions(context.Background(),
			&processorv1.ListPaymentActionFilters{ActionIds: []string{paymentAction.Id}})
		require.NoError(t, err)
		assert.Len(t, p, 1)

		assert.Equal(t, paymentAction, p[0])
	})
	t.Run("should return no actions if action doesnt exist", func(t *testing.T) {
		actions, err := testStore.ListPaymentActions(context.Background(), &processorv1.ListPaymentActionFilters{ActionIds: []string{uuid.NewV4().String()}})
		require.NoError(t, err)
		assert.NotNil(t, actions)
	})
}
