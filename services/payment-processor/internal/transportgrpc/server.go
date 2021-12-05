package transportgrpc

import (
	"context"
	"google.golang.org/grpc"

	paymentprocessorv1 "github.com/jacktantram/payments-api/build/go/rpc/paymentprocessor/v1"
)

type Server struct {
	paymentprocessorv1.UnimplementedPaymentProcessorServer
}

func NewServer() *grpc.Server {
	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	paymentprocessorv1.RegisterPaymentProcessorServer(s, &Server{})
	return s
}

func (s Server) CreatePayment(ctx context.Context, request *paymentprocessorv1.CreatePaymentRequest) (*paymentprocessorv1.CreatePaymentResponse, error) {
	//TODO implement me
	request.
		panic("implement me")
}

func (s Server) Capture(ctx context.Context, request *paymentprocessorv1.CreateCaptureRequest) (*paymentprocessorv1.CreateCaptureResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) Refund(ctx context.Context, request *paymentprocessorv1.CreateRefundRequest) (*paymentprocessorv1.CreateRefundResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) Void(ctx context.Context, request *paymentprocessorv1.CreateVoidRequest) (*paymentprocessorv1.CreateVoidResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetPayment(ctx context.Context, request *paymentprocessorv1.GetPaymentRequest) (*paymentprocessorv1.GetPaymentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) ListPaymentActions(ctx context.Context, request *paymentprocessorv1.ListPaymentActionsRequest) (*paymentprocessorv1.ListPaymentActionsResponse, error) {
	//TODO implement me
	panic("implement me")
}
