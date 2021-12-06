// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
	v10 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	domain "github.com/jacktantram/payments-api/services/payment-gateway/internal/domain"
)

// MockGateway is a mock of Gateway interface.
type MockGateway struct {
	ctrl     *gomock.Controller
	recorder *MockGatewayMockRecorder
}

// MockGatewayMockRecorder is the mock recorder for MockGateway.
type MockGatewayMockRecorder struct {
	mock *MockGateway
}

// NewMockGateway creates a new mock instance.
func NewMockGateway(ctrl *gomock.Controller) *MockGateway {
	mock := &MockGateway{ctrl: ctrl}
	mock.recorder = &MockGatewayMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGateway) EXPECT() *MockGatewayMockRecorder {
	return m.recorder
}

// Capture mocks base method.
func (m *MockGateway) Capture(ctx context.Context, paymentID string, amount uint64) (*v10.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Capture", ctx, paymentID, amount)
	ret0, _ := ret[0].(*v10.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Capture indicates an expected call of Capture.
func (mr *MockGatewayMockRecorder) Capture(ctx, paymentID, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capture", reflect.TypeOf((*MockGateway)(nil).Capture), ctx, paymentID, amount)
}

// CreatePayment mocks base method.
func (m *MockGateway) CreatePayment(ctx context.Context, amount *v1.Money, method domain.PaymentMethod) (*v10.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePayment", ctx, amount, method)
	ret0, _ := ret[0].(*v10.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePayment indicates an expected call of CreatePayment.
func (mr *MockGatewayMockRecorder) CreatePayment(ctx, amount, method interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePayment", reflect.TypeOf((*MockGateway)(nil).CreatePayment), ctx, amount, method)
}

// Refund mocks base method.
func (m *MockGateway) Refund(ctx context.Context, paymentID string, amount uint64) (*v10.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refund", ctx, paymentID, amount)
	ret0, _ := ret[0].(*v10.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refund indicates an expected call of Refund.
func (mr *MockGatewayMockRecorder) Refund(ctx, paymentID, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refund", reflect.TypeOf((*MockGateway)(nil).Refund), ctx, paymentID, amount)
}

// Void mocks base method.
func (m *MockGateway) Void(ctx context.Context, paymentID string) (*v10.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Void", ctx, paymentID)
	ret0, _ := ret[0].(*v10.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Void indicates an expected call of Void.
func (mr *MockGatewayMockRecorder) Void(ctx, paymentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Void", reflect.TypeOf((*MockGateway)(nil).Void), ctx, paymentID)
}