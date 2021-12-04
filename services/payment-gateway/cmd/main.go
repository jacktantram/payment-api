package main

import (
	"context"
	"github.com/jacktantram/payments-api/pkg/driver/v1/config"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/transporthttp"
	log "github.com/sirupsen/logrus"
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
}

func main(){
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

	handler, err:= transporthttp.NewHandler()
	if err!=nil{
		log.WithError(err).Fatalf("unable to setup handler")
	}

	srv := &http.Server{
		Handler:      transporthttp.HandleRoutes(handler),
		Addr:         ":8000",
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		ReadTimeout: time.Duration(cfg.ReadTimeout) * time.Second,
	}
	if err=srv.ListenAndServe();err!=nil{
		log.WithError(err).Fatal("unable to listen and serve")
	}
}
