package main

import (
	"github.com/jacktantram/payments-api/pkg/driver/v1/config"
	"github.com/jacktantram/payments-api/services/payment-processor/internal/transportgrpc"
	log "github.com/sirupsen/logrus"
	"net"
)

// Cfg represents the services config
type Cfg struct {
	MigrationPath string `envconfig:"MIGRATION_PATH" default:"/migrations"`
}

func main() {
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
