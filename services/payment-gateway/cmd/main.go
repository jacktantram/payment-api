package main

import (
	"context"
	processorv1 "github.com/jacktantram/payments-api/build/go/rpc/paymentprocessor/v1"
	"github.com/jacktantram/payments-api/pkg/driver/v1/config"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/handler"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	serviceName = "payment-gateway"
)

// Cfg represents the services config
type Cfg struct {
	config.HTTPConfig
	Hostnames struct {
		PaymentProcessor string `envconfig:"PAYMENT_PROCESSOR_HOSTNAME"`
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		<-c
		cancel()
	}()

	cfg := &Cfg{}
	if err := config.LoadConfig(cfg); err != nil {
		log.WithError(err).Fatalf("unable to load config")
	}

	grpcClient, err := grpc.Dial(cfg.Hostnames.PaymentProcessor, grpc.WithInsecure())
	if err != nil {
		log.WithError(err).Fatalf("unable to load config")
	}

	h, err := handler.NewHandler(processorv1.NewPaymentProcessorClient(grpcClient))
	if err != nil {
		log.WithError(err).Fatalf("unable to setup handler")
	}

	srv := &http.Server{
		Handler:      handler.HandleRoutes(h),
		Addr:         ":8080",
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
	}
	if err = srv.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("unable to listen and serve")
	}
}
