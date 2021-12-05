package main

import (
	"context"
	"github.com/jacktantram/payments-api/pkg/driver/v1/config"
	"github.com/jacktantram/payments-api/services/payment-processor/internal/transportgrpc"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
)

const (
	serviceName = "payment-processor"
)

// Cfg represents the services config
type Cfg struct {
	config.HTTPConfig
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

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	s := transportgrpc.NewServer()
	if err := s.Serve(lis); err != nil {
		log.WithError(err).Fatal("service shutting down")
	}
}
