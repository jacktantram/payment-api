//+build integration

package store_test

import (
	"github.com/jacktantram/payments-api/pkg/driver/v1/postgres"
	"github.com/jacktantram/payments-api/services/payment-gateway/internal/store"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

var (
	testStore store.Store
)

func TestMain(m *testing.M) {
	postgresClient, err := postgres.NewClient("postgres://postgres:postgres@localhost:5432?sslmode=disable", "payment")
	if err != nil {
		log.Fatal(err)
	}
	if err := postgresClient.Migrate("../migrations"); err != nil {
		log.Fatal(err)
	}
	testStore = store.NewStore(postgresClient)
	exitVal := m.Run()
	os.Exit(exitVal)
}
