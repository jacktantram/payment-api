//go:generate mockgen --destination=mocks/mock_processor.go -package=mocks github.com/jacktantram/payments-api/build/go/rpc/paymentprocessor/v1 PaymentProcessorClient
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"strings"
	"time"

	processorv1 "github.com/jacktantram/payments-api/build/go/rpc/paymentprocessor/v1"
)

const (
	CardNumberLen = 16
	CVVLen        = 3
	CurrencyLen   = 3
	ExpiryMonLen  = 12
)

func HandleRoutes(h Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/authorize", h.AuthorizeHandler).Methods(http.MethodPost)
	r.HandleFunc("/capture", h.CaptureHandler).Methods(http.MethodPost)
	r.HandleFunc("/refund", h.RefundHandler).Methods(http.MethodPost)
	r.HandleFunc("/void", h.VoidHandler).Methods(http.MethodPost)

	return r
}

type Handler struct {
	processorClient processorv1.PaymentProcessorClient
}

func NewHandler(processorClient processorv1.PaymentProcessorClient) (Handler, error) {
	if processorClient == nil {
		return Handler{}, errors.New("processor client is nil")
	}
	return Handler{processorClient: processorClient}, nil
}

func (h Handler) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		http.Error(w, "no body supplied", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var authorizationRequest CreateAuthorizationRequest
	if err := json.NewDecoder(r.Body).Decode(&authorizationRequest); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	validateRequest := func() error {
		if authorizationRequest.Amount.MinorUnits == 0 {
			return errors.New("invalid amount.minor_units: cannot be zero")
		}
		if len(authorizationRequest.Amount.Currency) != CurrencyLen {
			return fmt.Errorf("invalid amount.currency: must be length of %d", CurrencyLen)
		}
		if authorizationRequest.Card == nil {
			return errors.New("missing payment method: cannot be empty")
		}
		if len(authorizationRequest.Card.CardNumber) != CardNumberLen && len(strings.ReplaceAll(authorizationRequest.Card.CardNumber, " ", "")) != CardNumberLen {
			return fmt.Errorf("invalid payment_method.card.card_number: length not equal to %d", CardNumberLen)
		}
		if len(authorizationRequest.Card.Cvv) != CVVLen {
			return fmt.Errorf("invalid payment_method.card.cvv: length not equal to %d", 3)
		}
		if authorizationRequest.Card.Expiry == nil {
			return errors.New("missing payment_method.card.expiry: cannot be empty")
		}
		if authorizationRequest.Card.Expiry.Month > 12 {
			return fmt.Errorf("invalid payment_method.card.expiry.month: expiry month cannot exceed %d", ExpiryMonLen)
		}
		if int(authorizationRequest.Card.Expiry.Year) < time.Now().Year() {
			return errors.New("missing payment_method.card.expiry.year: cannot be in the past")
		}
		return nil
	}
	if err := validateRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	logFields := log.Fields{
		"amount.minor_units": authorizationRequest.Amount.MinorUnits,
		"amount.currency":    authorizationRequest.Amount.Currency,
		// Ideally a card ID/token would be better
		"card.first_six": authorizationRequest.Card.CardNumber[:6],
		"card.last_four": authorizationRequest.Card.CardNumber[len(authorizationRequest.Card.CardNumber)-3:],
	}

	fn := func() error {
		paymentResponse, err := h.processorClient.CreatePayment(r.Context(), &processorv1.CreatePaymentRequest{
			Amount:        &authorizationRequest.Amount,
			PaymentMethod: &processorv1.CreatePaymentRequest_Card{Card: authorizationRequest.Card},
		})
		if err != nil {
			return err
		}
		paymentBytes, err := protojson.Marshal(paymentResponse.GetPayment())
		if err != nil {
			logFields["payment.id"] = paymentResponse.Payment.Id
			return err
		}
		_, err = w.Write(paymentBytes)
		if err != nil {
			return err
		}
		return nil
	}

	if err := fn(); err != nil {
		logFields["error"] = err
		log.WithFields(logFields).Error("failed to process authorization request")
		http.Error(w, "Oops something went wrong", http.StatusInternalServerError)
		return
	}
}

func (h Handler) CaptureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		http.Error(w, "no body supplied", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var captureRequest CreateCaptureRequest
	if err := json.NewDecoder(r.Body).Decode(&authorizationRequest); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

}

func (h Handler) RefundHandler(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) VoidHandler(w http.ResponseWriter, r *http.Request) {

}
