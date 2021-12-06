package main

import (
	"context"
	paymentsV1 "github.com/jacktantram/payments-api/build/go/shared/payment/v1"
	"github.com/jacktantram/payments-api/pkg/driver/v1/config"
	"github.com/jacktantram/payments-api/pkg/driver/v1/postgres"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/domain"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/gateway"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/store"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/transport/transporthttp"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

// Cfg represents the services config
type Cfg struct {
	config.HTTPConfig
	Hostnames struct {
		PaymentProcessor string `envconfig:"PAYMENT_PROCESSOR_HOSTNAME"`
	}
	DatabaseURI   string `envconfig:"DATABASE_URI"`
	MigrationPath string `envconfig:"MIGRATION_PATH" default:"/migrations"`
}

func main() {
	cfg := &Cfg{}

	if err := config.LoadConfig(cfg); err != nil {
		log.WithError(err).Fatalf("unable to load config")
	}

	client, err := postgres.NewClient(cfg.DatabaseURI, "postgres")
	if err != nil {
		log.WithError(err).Fatal("failed to setup postgres client")
	}
	defer client.DB.Close()

	if err = client.Migrate(cfg.MigrationPath); err != nil {
		log.WithError(err).Fatalf("unable to migrate")
	}

	h, err := transporthttp.NewHandler(gateway.NewService(store.NewStore(client), FakeGateway{}))
	if err != nil {
		log.WithError(err).Fatalf("unable to setup transporthttp")
	}

	srv := &http.Server{
		Handler:      transporthttp.HandleRoutes(h),
		Addr:         ":8080",
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
	}

	if err = srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("unable to listen and serve")
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
