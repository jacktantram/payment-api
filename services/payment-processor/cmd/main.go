package main

import (
	"context"
	"errors"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	"github.com/jacktantram/payments-api/pkg/driver/v1/config"
	"github.com/jacktantram/payments-api/pkg/driver/v1/postgres"
	"github.com/jacktantram/payments-api/services/payment-processor/internal/domain"
	"github.com/jacktantram/payments-api/services/payment-processor/internal/processor"
	"github.com/jacktantram/payments-api/services/payment-processor/internal/store"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
)

// Cfg represents the services config
type Cfg struct {
	DatabaseURI   string `envconfig:"DATABASE_URI"`
	MigrationPath string `envconfig:"MIGRATION_PATH" default:"/migrations"`
}

func main() {
	cfg := &Cfg{}
	if err := config.LoadConfig(cfg); err != nil {
		log.WithError(err).Fatalf("unable to load config")
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	client, err := postgres.NewClient(cfg.DatabaseURI, "postgres")
	if err != nil {
		log.WithError(err).Fatal("failed to setup postgres client")
	}

	s := processor.NewService(store.NewStore(client), FakeGateway{})
	if err := s.S.Serve(lis); err != nil {
		log.WithError(err).Fatal("service shutting down")
	}
}

type FakeGateway struct {
}

var bannedCards = map[string]paymentsV1.PaymentType{
	"4000000000000119": paymentsV1.PaymentType_PAYMENT_TYPE_AUTHORIZATION,
	"4000000000000259": paymentsV1.PaymentType_PAYMENT_TYPE_CAPTURE,
	"4000000000003238": paymentsV1.PaymentType_PAYMENT_TYPE_REFUND,
}

func (f FakeGateway) CreateIssuerRequest(ctx context.Context, issuerRequest domain.IssuerRequest) (domain.IssuerResponse, error) {
	if issuerRequest.PaymentMethod.Card == nil {
		return domain.IssuerResponse{}, errors.New("unsupported method")
	}
	cardNumber := strings.ReplaceAll(issuerRequest.PaymentMethod.Card.CardNumber, " ", "")
	if methodType, ok := bannedCards[cardNumber]; ok {
		if methodType == issuerRequest.OperationType {
			return domain.IssuerResponse{AuthCode: "12"}, nil
		}
	}
	return domain.IssuerResponse{AuthCode: "00"}, nil
}
