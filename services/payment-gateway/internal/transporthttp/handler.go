package transporthttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type PaymentGateway interface{}

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
}

func NewHandler() (Handler, error) {
	return Handler{}, nil
}

func (h Handler) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		http.Error(w, "no body supplied", http.StatusBadRequest)
		return
	}
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
		// could strip spaces to allow a customer to type like 4000 0000
		if len(authorizationRequest.Card.CardNumber) != CardNumberLen {
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}

func (h Handler) CaptureHandler(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) RefundHandler(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) VoidHandler(w http.ResponseWriter, r *http.Request) {

}
